package telegram

import (
	"context"
	"fmt"
	"strings"

	internal_models "github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func (c *Client) handleProtocols(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	text := bot.EscapeMarkdown(update.Message.Text)
	protocols, err := parseCommandArgs(text)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle protocols with error")
		c.sendMessage(ctx, chatID, "handle protocols with error: "+err.Error())
		return
	}

	err = c.updateConfig(ctx, chatID, func(cfg *internal_models.Config) {
		cfg.Protocols = protocols
	})
	if err != nil {
		return
	}

	c.sendMessage(ctx, chatID, fmt.Sprintf("Protocols updated successfully: %s", strings.Join(protocols, ", ")))
}
