package telegram

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleProtocols(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error
	text := bot.EscapeMarkdown(update.Message.Text)
	c.protocols, err = getItemsFromText(text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      err.Error(),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      fmt.Sprintf("Protocols updated successfully: %s", strings.Join(c.protocols, ", ")),
		ParseMode: models.ParseModeMarkdown,
	})
}
