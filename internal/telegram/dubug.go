package telegram

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func (c *Client) handleDebug(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	data, err := os.ReadFile("text.txt")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("handle debug with error")
		c.sendMessage(ctx, chatID, fmt.Sprintf("got err while reading file: %v", err))
		return
	}

	c.sendMessage(ctx, chatID, string(data))
}
