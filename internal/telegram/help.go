package telegram

import (
	"context"

	_ "embed"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

//go:embed data/help.txt
var helpMessage string

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update != nil && update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      helpMessage,
			ParseMode: models.ParseModeMarkdown,
		})
	}
}
