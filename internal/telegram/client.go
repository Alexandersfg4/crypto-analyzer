package telegram

import (
	"context"
	"fmt"

	"github.com/Alexandersfg4/crypto-analyzer/internal/cron"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/go-telegram/bot"
)

// Send any text message to the bot after the bot has been started

type Client struct {
	b          *bot.Bot
	r          *report.Report
	tokens     []string
	protocols  []string
	reportCron *cron.Cron
}

func New(apiToken string, r *report.Report) (*Client, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handleDefault),
	}

	b, err := bot.New(apiToken, opts...)
	if nil != err {
		return nil, fmt.Errorf("failed init new bot: %w", err)
	}

	c := &Client{
		b: b,
		r: r,
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "help", bot.MatchTypeCommand, handleDefault)
	b.RegisterHandler(bot.HandlerTypeMessageText, "report", bot.MatchTypeCommand, c.handleReport)
	b.RegisterHandler(bot.HandlerTypeMessageText, "tokens", bot.MatchTypeCommandStartOnly, c.handleTokens)
	b.RegisterHandler(bot.HandlerTypeMessageText, "protocols", bot.MatchTypeCommandStartOnly, c.handleProtocols)
	b.RegisterHandler(bot.HandlerTypeMessageText, "config", bot.MatchTypeCommandStartOnly, c.handleConfig)
	b.RegisterHandler(bot.HandlerTypeMessageText, "cron", bot.MatchTypeCommandStartOnly, c.handleCron)
	b.RegisterHandler(bot.HandlerTypeMessageText, "debug", bot.MatchTypeCommandStartOnly, c.handleDebug)

	return c, nil
}

func (c *Client) Start(ctx context.Context) {
	c.b.Start(ctx)
}
