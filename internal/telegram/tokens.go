package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleTokens(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := bot.EscapeMarkdown(update.Message.Text)
	tokens := strings.Split(text, ",")
	c.tokens = tokens

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Tokens updated successfully: %s", text),
		ParseMode: models.ParseModeMarkdown,
	})
}
