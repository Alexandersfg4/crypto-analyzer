package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleTokens(ctx context.Context, b *bot.Bot, update *models.Update) {
	var err error
	text := bot.EscapeMarkdown(update.Message.Text)
	c.tokens, err = getItemsFromText(text)
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
		Text:      fmt.Sprintf("Tokens updated successfully: %s", strings.Join(c.tokens, ", ")),
		ParseMode: models.ParseModeMarkdown,
	})
}

func getItemsFromText(text string) ([]string, error) {
	msgs := strings.Split(text, " ")
	if len(msgs) < 2 {
		return nil, errors.New("invalid input")
	}

	return strings.Split(msgs[1], ","), nil
}
