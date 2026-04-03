package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func (c *Client) handleProtocols(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	text := bot.EscapeMarkdown(update.Message.Text)
	protocols, err := getItemsFromText(text)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle protocols with error")
		c.sendMessage(ctx, chatID, "handle protocols with error: "+err.Error())
		return
	}

	config, err := c.configStorage.SaveProtocols(protocols)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("failed to save protocols")
		c.sendMessage(ctx, chatID, "failed to save protocols: "+err.Error())
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Protocols updated successfully: %s", strings.Join(config.Protocols, ", ")),
		ParseMode: models.ParseModeMarkdown,
	})
}
