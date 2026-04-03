package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleConfig(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	config, err := c.configStorage.Read()
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("failed to read config: %s", err.Error()))
		return
	}
	err = config.Validate()
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("failed to validate config: %s", err.Error()))
		return
	}

	text := fmt.Sprintf("*Config*\nTokens: %s\nProtocols: %s\nCron Next Execution Time: %s",
		strings.Join(config.Tokens, ", "),
		strings.Join(config.Protocols, ", "),
		config.CronNextExecutionTime.Format(time.RFC3339))

	c.sendMessage(ctx, chatID, text)
}
