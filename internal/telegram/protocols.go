package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleProtocols(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := bot.EscapeMarkdown(update.Message.Text)
	protocols := strings.Split(text, ",")
	c.protocols = protocols

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Protocols updated successfully: %s", text),
		ParseMode: models.ParseModeMarkdown,
	})
}
