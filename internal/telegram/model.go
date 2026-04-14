package telegram

import (
	"context"
	"fmt"

	internal_models "github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func (c *Client) handleModel(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text
	model, err := parseCommandArg(text)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle model with error")
		c.sendMessage(ctx, chatID, "handle model with error: "+err.Error())
		return
	}

	err = c.updateConfig(ctx, chatID, func(cfg *internal_models.Config) {
		cfg.OpenrouterModel = model
	})
	if err != nil {
		return
	}

	c.sendMessage(ctx, chatID, fmt.Sprintf("Model updated successfully: %s", model))
}
