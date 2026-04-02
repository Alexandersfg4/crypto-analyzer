package telegram

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/cron"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const reportInterval = time.Hour * 24

func (c *Client) handleCron(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error
	text := bot.EscapeMarkdown(update.Message.Text)

	duration, err := getDurationFromText(text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      err.Error(),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	c.reportCron = cron.New(duration)

	go c.reportCron.Run(func() {
		c.sendReport(ctx, update.Message.Chat.ID)
	})

	log.WithFields(log.Fields{
		"chatID": update.Message.Chat.ID,
	}).Info("cron set successfully")
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "cron set successfully",
		ParseMode: models.ParseModeMarkdown,
	})
}

func getDurationFromText(text string) (time.Duration, error) {
	msgs := strings.Split(text, " ")
	if len(msgs) < 2 {
		return 0, errors.New("invalid input")
	}

	return time.ParseDuration(msgs[1])
}
