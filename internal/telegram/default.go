package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const defaultMessage = `/report - generate crypto report
/tokens ETH,USDT - set observed tokens
/protocols AAVE,UNISWAP - set observed DEX protocols`

func handleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update != nil && update.Message != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      defaultMessage,
			ParseMode: models.ParseModeMarkdown,
		})
	}
}
