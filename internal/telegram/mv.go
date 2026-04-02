package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

func auth(userID int64) func(next bot.HandlerFunc) bot.HandlerFunc {
	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			if update.Message.From.ID != userID {
				log.WithFields(log.Fields{
					"user_id":  update.Message.From.ID,
					"username": update.Message.From.Username,
				}).Info("unauthorized user")
				return
			}

			next(ctx, bot, update)
		}
	}
}
