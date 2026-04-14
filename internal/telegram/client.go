package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/cron"

	internal_models "github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

type Generator interface {
	Generate(ctx context.Context, cfg internal_models.Config) (internal_models.Report, error)
}

type ConfigStorer interface {
	Read() (*internal_models.Config, error)
	Save(cfg *internal_models.Config) error
}

type Client struct {
	b             *bot.Bot
	r             Generator
	configStorage ConfigStorer
	reportCron    *cron.Cron
}

func New(apiToken string, userID int64, r Generator, store ConfigStorer) (*Client, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handleHelp),
		bot.WithMiddlewares(auth(userID), recoverFromPanic()),
	}

	b, err := bot.New(apiToken, opts...)
	if nil != err {
		return nil, fmt.Errorf("failed init new bot: %w", err)
	}

	c := &Client{
		b:             b,
		r:             r,
		configStorage: store,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommandStartOnly, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "report", bot.MatchTypeCommandStartOnly, c.handleReport)
	b.RegisterHandler(bot.HandlerTypeMessageText, "tokens", bot.MatchTypeCommandStartOnly, c.handleTokens)
	b.RegisterHandler(bot.HandlerTypeMessageText, "protocols", bot.MatchTypeCommandStartOnly, c.handleProtocols)
	b.RegisterHandler(bot.HandlerTypeMessageText, "config", bot.MatchTypeCommandStartOnly, c.handleConfig)
	b.RegisterHandler(bot.HandlerTypeMessageText, "cron", bot.MatchTypeCommandStartOnly, c.handleCron)
	b.RegisterHandler(bot.HandlerTypeMessageText, "model", bot.MatchTypeCommandStartOnly, c.handleModel)

	return c, nil
}

func (c *Client) Start(ctx context.Context) {
	cfg, err := c.configStorage.Read()
	if err != nil {
		panic(err)
	}

	if !cfg.CronNextExecutionTime.IsZero() && cfg.ChatID != 0 {

		now := time.Now()
		cr := cron.New(cfg.CronNextExecutionTime)

		if cfg.CronNextExecutionTime.Before(now) {
			cr.Reset(cfg.CronNextExecutionTime)
		}

		go cr.Run(ctx, func() {
			c.sendReport(ctx, cfg.ChatID)
			c.updateConfig(ctx, cfg.ChatID, func(cfg *internal_models.Config) {
				cfg.CronNextExecutionTime = c.reportCron.ExecutionTime()
			})
		})

		c.reportCron = cr
	}

	c.b.Start(ctx)
}

func (c *Client) sendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      processText(text),
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("send messge with err ", text)
		c.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: models.ParseModeMarkdownV1,
		})
	}

	return nil
}

func (c *Client) updateConfig(ctx context.Context, chatID int64, upater func(*internal_models.Config)) error {
	cfg, err := c.configStorage.Read()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("read config with error")
		c.sendMessage(ctx, chatID, "read config with error: "+err.Error())
		return err
	}

	upater(cfg)

	err = c.configStorage.Save(cfg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("save config with error")
		c.sendMessage(ctx, chatID, "save config with error: "+err.Error())
		return err
	}

	return nil
}
