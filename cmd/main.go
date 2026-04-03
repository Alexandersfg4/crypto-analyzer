package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinmarketcap"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/defillama"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/Alexandersfg4/crypto-analyzer/internal/telegram"
)

const (
	envCoinstatsAPIKey     = "COINSTATS_API_KEY"
	envCoinmarketcapAPIKey = "API_KEY_COINMARKETCAP"
	envTelegramToken       = "TELEGRAM_API_TOKEN"
	envTelegramUserID      = "TELEGRAM_USER_ID"
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()

	apiKeyCoinstats, ok := os.LookupEnv(envCoinstatsAPIKey)
	if !ok {
		log.Fatal("env COINSTATS_API_KEY not found")
	}

	apiKeyCoinmarketcap, ok := os.LookupEnv(envCoinmarketcapAPIKey)
	if !ok {
		log.Fatal("env API_KEY_COINMARKETCAP not found")
	}

	apiKeyTelegram, ok := os.LookupEnv(envTelegramToken)
	if !ok {
		log.Fatal("env TELEGRAM_API_TOKEN not found")
	}

	telegramUserID, ok := os.LookupEnv(envTelegramUserID)
	if !ok {
		log.Fatal("env TELEGRAM_USER_ID not found")
	}

	telegramUserIDInt, err := strconv.ParseInt(telegramUserID, 10, 64)
	if err != nil {
		log.Fatal("failed to parse TELEGRAM_USER_ID: ", err.Error())
	}

	ctx := context.Background()

	coinmarketcapSrv := coinmarketcap.NewService(apiKeyCoinmarketcap)
	coinstatsSrv := coinstats.NewService(apiKeyCoinstats)
	defillamaSrv := defillama.NewService()

	r := report.New(coinmarketcapSrv, coinstatsSrv, defillamaSrv)

	tg, err := telegram.New(apiKeyTelegram, telegramUserIDInt, r)
	if err != nil {
		fmt.Println("failed to create telegram client: ", err.Error())
		os.Exit(1)
	}

	tg.Start(ctx)
}
