package telegram

import (
	"context"
	"fmt"

	internal_models "github.com/Alexandersfg4/crypto-analyzer/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	b      *bot.Bot
	chatID int64
}

func New(apiToken string, chatID int64) (*Client, error) {
	b, err := bot.New(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed init new bot: %w", err)
	}

	return &Client{
		b:      b,
		chatID: chatID,
	}, nil
}

func (c *Client) SendReport(ctx context.Context, report internal_models.Report) error {
	messages := []string{
		report.MarketCap,
		report.Tokens.InPortfolio,
		report.Tokens.GainersAndLoosers,
		report.AISummary,
	}

	for _, message := range messages {
		if message == "" {
			continue
		}

		if err := c.sendMessage(ctx, message); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) sendMessage(ctx context.Context, text string) error {
	_, err := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    c.chatID,
		Text:      processText(text),
		ParseMode: models.ParseModeMarkdown,
	})
	if err == nil {
		return nil
	}

	log.WithFields(log.Fields{
		"err": err,
	}).Warn("failed to send message with markdownv2, retrying with markdown")

	_, fallbackErr := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    c.chatID,
		Text:      text,
		ParseMode: models.ParseModeMarkdownV1,
	})
	if fallbackErr != nil {
		return fmt.Errorf("failed to send telegram message: markdownv2: %w, markdown: %v", err, fallbackErr)
	}

	return nil
}
