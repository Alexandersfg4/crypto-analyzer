package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	_ "embed"

	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (c *Client) handleDebug(ctx context.Context, b *bot.Bot, update *models.Update) {
	var data report.Data
	chatID := update.Message.Chat.ID

	jsonData, err := os.ReadFile("text.json")
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("got err while reading file: %v", err))
		return
	}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		c.sendMessage(ctx, chatID, fmt.Sprintf("got err while unmarshalling data: %v", err))
		return
	}

	cap := &bytes.Buffer{}
	formatter.MarketCap(cap, data.MarketCap)
	formatter.FearAndGreed(cap, data.FeatAndGreed)
	c.sendMessage(ctx, chatID, cap.String())

	tokens := &bytes.Buffer{}
	formatter.Coins(tokens, data.ListingsLatest, c.tokens)
	c.sendMessage(ctx, chatID, tokens.String())

	news := &bytes.Buffer{}
	formatter.News(news, data.News)
	c.sendMessage(ctx, chatID, news.String())

	protocols := &bytes.Buffer{}
	formatter.Protocols(protocols, data.Protocols, c.protocols)
	c.sendMessage(ctx, chatID, protocols.String())
	c.sendReport(ctx, update.Message.Chat.ID)
}
