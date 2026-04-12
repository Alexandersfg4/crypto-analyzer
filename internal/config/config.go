package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	envCoinstatsAPIKey     = "COINSTATS_API_KEY"
	envCoinmarketcapAPIKey = "API_KEY_COINMARKETCAP"
	envTelegramToken       = "TELEGRAM_API_TOKEN"
	envTelegramUserID      = "TELEGRAM_USER_ID"
	envOpenRouterAPIKey    = "OPENROUTER_API_KEY"
)

type Config struct {
	CoinstatsAPIKey     string
	CoinmarketcapAPIKey string
	TelegramToken       string
	TelegramUserID      int64
	OpenRouterAPIKey    string
}

func Load() (*Config, error) {
	cfg := &Config{}

	coinstatsKey, ok := os.LookupEnv(envCoinstatsAPIKey)
	if !ok {
		return nil, errors.New("env COINSTATS_API_KEY not found")
	}
	cfg.CoinstatsAPIKey = coinstatsKey

	coinmarketcapKey, ok := os.LookupEnv(envCoinmarketcapAPIKey)
	if !ok {
		return nil, errors.New("env API_KEY_COINMARKETCAP not found")
	}
	cfg.CoinmarketcapAPIKey = coinmarketcapKey

	telegramToken, ok := os.LookupEnv(envTelegramToken)
	if !ok {
		return nil, errors.New("env TELEGRAM_API_TOKEN not found")
	}
	cfg.TelegramToken = telegramToken

	telegramUserID, ok := os.LookupEnv(envTelegramUserID)
	if !ok {
		return nil, errors.New("env TELEGRAM_USER_ID not found")
	}
	userID, err := strconv.ParseInt(telegramUserID, 10, 64)
	if err != nil {
		return nil, errors.New("failed to parse TELEGRAM_USER_ID: " + err.Error())
	}
	cfg.TelegramUserID = userID

	openRouterKey, ok := os.LookupEnv(envOpenRouterAPIKey)
	if !ok {
		return nil, errors.New("env OPENROUTER_API_KEY not found")
	}
	cfg.OpenRouterAPIKey = openRouterKey

	return cfg, nil
}
