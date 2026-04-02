package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleConfig(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := fmt.Sprintf("*Config*\nTokens: %s\nProtocols: %s", strings.Join(c.tokens, ", "), strings.Join(c.protocols, ", "))

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
}
