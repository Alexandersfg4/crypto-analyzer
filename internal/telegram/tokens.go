package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func (c *Client) handleTokens(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	text := bot.EscapeMarkdown(update.Message.Text)
	tokens, err := getItemsFromText(text)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("failed to get tokens from text")
		c.sendMessage(ctx, chatID, "failed to get tokens from text: "+err.Error())
		return
	}

	config, err := c.configStorage.SaveTokens(tokens)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("failed to save tokens")
		c.sendMessage(ctx, chatID, "failed to save tokens: "+err.Error())
		return
	}

	c.sendMessage(ctx, chatID, fmt.Sprintf("Tokens updated successfully: %s", strings.Join(config.Tokens, ", ")))
}

func getItemsFromText(text string) ([]string, error) {
	msgs := strings.Split(text, " ")
	if len(msgs) < 2 {
		return nil, errors.New("invalid input")
	}

	return strings.Split(msgs[1], ","), nil
}
