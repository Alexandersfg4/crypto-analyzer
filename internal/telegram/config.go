package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleConfig(ctx context.Context, b *bot.Bot, update *models.Update) {
	config, err := c.configStorage.Read()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "failed to read config" + err.Error(),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}
	err = config.Validate()
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "failed to validate config" + err.Error(),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	text := fmt.Sprintf("*Config*\nTokens: %s\nProtocols: %s",
		strings.Join(config.Tokens, ", "), strings.Join(config.Protocols, ", "))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
}
