package telegram

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleReport(ctx context.Context, b *bot.Bot, update *models.Update) {
	var (
		err  error
		buff *bytes.Buffer
	)
	if len(c.protocols) == 0 || len(c.tokens) == 0 {
		err = fmt.Errorf("no tokens or protocols provided")
	} else {
		buff, err = c.r.Generate(ctx, c.tokens, c.protocols)
	}
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "got err while generating report: " + err.Error(),
			ParseMode: models.ParseModeMarkdown,
		})
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      buff.String(),
		ParseMode: models.ParseModeMarkdown,
	})
}
