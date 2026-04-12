package main

import (
	"context"
	"flag"
	"os"

	"github.com/Alexandersfg4/crypto-analyzer/internal/config"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinmarketcap"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/coinstats"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/defillama"
	"github.com/Alexandersfg4/crypto-analyzer/internal/providers/openrouter"
	"github.com/Alexandersfg4/crypto-analyzer/internal/report"
	"github.com/Alexandersfg4/crypto-analyzer/internal/storage"
	"github.com/Alexandersfg4/crypto-analyzer/internal/telegram"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()

	coinmarketcapSrv := coinmarketcap.New(cfg.CoinmarketcapAPIKey)
	coinstatsSrv := coinstats.New(cfg.CoinstatsAPIKey)
	defillamaSrv := defillama.New()
	openRouterSrv := openrouter.New(cfg.OpenRouterAPIKey)

	r := report.New(coinmarketcapSrv, coinstatsSrv, defillamaSrv, openRouterSrv)

	configStorage := storage.New("crypto-analyzer.json")
	tg, err := telegram.New(cfg.TelegramToken, cfg.TelegramUserID, r, configStorage)
	if err != nil {
		logrus.Fatal("failed to create telegram client: ", err)
	}

	tg.Start(ctx)
}
