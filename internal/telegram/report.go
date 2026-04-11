package telegram

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const maxMessageLength = 4096

func (c *Client) handleReport(ctx context.Context, b *bot.Bot, update *models.Update) {
	c.sendReport(ctx, update.Message.Chat.ID)
}

func (c *Client) sendReport(ctx context.Context, chatID int64) {
	config, err := c.configStorage.Read()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while reading config: "+err.Error())
		return
	}
	err = config.Validate()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while config validation: "+err.Error())
		return
	}

	data, err := c.r.Generate(ctx, *config)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Info("handle report with error")
		c.sendMessage(ctx, chatID, "got err while generating report: "+err.Error())
		return
	}

	c.sendMessage(ctx, chatID, data.MarketCap)
	c.sendMessage(ctx, chatID, data.Tokens.InPortfolio)
	c.sendMessage(ctx, chatID, data.Tokens.GainersAndLoosers)
	c.sendMessage(ctx, chatID, data.AISummary)
}

func (c *Client) sendMessage(ctx context.Context, chatID int64, text string) error {
	_, err := c.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      processText(text),
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("send messge with err ", text)
		c.b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatID,
			Text:      text,
			ParseMode: models.ParseModeMarkdownV1,
		})
	}

	return nil
}
