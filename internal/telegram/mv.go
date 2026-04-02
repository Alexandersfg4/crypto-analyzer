package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Auth(userID int64) func(next bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.Message.From.ID != userID {
				return
			}

			next(ctx, bot, update)
		}
	}
}
