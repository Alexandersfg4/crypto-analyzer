package telegram

import (
	"bytes"
	"context"

	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const maxMessageLength = 4096

func (c *Client) handleReport(ctx context.Context, b *bot.Bot, update *models.Update) {
	c.sendReport(ctx, update.Message.Chat.ID)
}

func (c *Client) sendReport(ctx context.Context, chatID int64) {
	var data report.Data
	config, err := c.configStorage.Read()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while reading config: "+err.Error())
		return
	}
	err = config.Validate()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while config validation: "+err.Error())
		return
	}
	if err == nil {
		data, err = c.r.Generate(ctx)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while generating report: "+err.Error())
		return
	}

	cap := &bytes.Buffer{}
	formatter.MarketCap(cap, data.MarketCap)
	formatter.FearAndGreed(cap, data.FeatAndGreed)
	c.sendMessage(ctx, chatID, cap.String())

	tokens := &bytes.Buffer{}
	formatter.Coins(tokens, data.ListingsLatest, config.Tokens)
	c.sendMessage(ctx, chatID, tokens.String())

	news := &bytes.Buffer{}
	formatter.News(news, data.News)
	c.sendMessage(ctx, chatID, news.String())

	protocols := &bytes.Buffer{}
	formatter.Protocols(protocols, data.Protocols, config.Protocols)
	c.sendMessage(ctx, chatID, protocols.String())
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
