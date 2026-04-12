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

const cronTimeLayout = "15:04"

func (c *Client) handleCron(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error
	chatID := update.Message.Chat.ID
	text := bot.EscapeMarkdown(update.Message.Text)

	arg, err := parseCommandArg(text)
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("invalid input: %s, expected format: %s", err.Error(), cronTimeLayout))
		return
	}

	t, err := time.Parse(cronTimeLayout, arg)
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("invalid input: %s, expected format: %s", err.Error(), cronTimeLayout))
		return
	}

	if c.reportCron == nil {
		c.reportCron = cron.New(t)

		go c.reportCron.Run(ctx, func() {
			c.sendReport(ctx, update.Message.Chat.ID)
			c.updateConfig(ctx, chatID, func(cfg *internal_models.Config) {
				cfg.CronNextExecutionTime = c.reportCron.ExecutionTime()
			})
		})
	} else {
		c.reportCron.Reset(t)
	}

	err = c.updateConfig(ctx, chatID, func(cfg *internal_models.Config) {
		cfg.CronNextExecutionTime = c.reportCron.ExecutionTime()
		cfg.ChatID = chatID
	})
	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"chatID": update.Message.Chat.ID,
	}).Info("cron set successfully")
	c.sendMessage(ctx, chatID, "cron set successfully")
}
