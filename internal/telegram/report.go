package telegram

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Alexandersfg4/crypto-analyzer/internal/formatter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const maxMessageLength = 4096

func (c *Client) handleReport(ctx context.Context, b *bot.Bot, update *models.Update) {
	c.sendReport(ctx, update.Message.Chat.ID)
}

func (c *Client) sendReport(ctx context.Context, chatID int64) {
	var (
		err  error
		data report.Data
	)
	if len(c.protocols) == 0 || len(c.tokens) == 0 {
		err = fmt.Errorf("no tokens or protocols provided")
	} else {
		data, err = c.r.Generate(ctx, c.tokens, c.protocols)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while generating report: "+err.Error())
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
}

func (c *Client) sendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("send messge with err")
		c.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   text,
		})
	}

	return nil
}
