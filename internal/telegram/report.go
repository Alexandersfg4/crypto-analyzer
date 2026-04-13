package telegram

import (
	"context"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	log "github.com/sirupsen/logrus"
)

const (
	timeoutReportExecution = time.Minute * 3
)

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

	ctx, cancel := context.WithTimeout(ctx, timeoutReportExecution)
	defer cancel()

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
