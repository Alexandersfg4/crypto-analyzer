package config

import (
	"errors"
	"flag"
	"os"
	"strings"
)

const (
	envCoinstatsAPIKey     = "COINSTATS_API_KEY"
	envCoinmarketcapAPIKey = "API_KEY_COINMARKETCAP"
	envTelegramToken       = "TELEGRAM_API_TOKEN"
	envOpenRouterAPIKey    = "OPENROUTER_API_KEY"
)

type Config struct {
	CoinstatsAPIKey     string
	CoinmarketcapAPIKey string
	TelegramToken       string
	TelegramChatID      int64
	OpenRouterAPIKey    string
	OpenRouterModel     string
	Tokens              []string
	Protocols           []string
}

func Load() (*Config, error) {
	cfg := &Config{}

	tokens := flag.String("tokens", "", "comma-separated list of tokens to include in the report")
	protocols := flag.String("protocols", "", "comma-separated list of DeFi protocols to include in the report")
	telegramChatID := flag.Int64("telegram-chat-id", 0, "telegram chat id to send the report to")
	openRouterModel := flag.String("openrouter-model", "", "OpenRouter model used for AI summary")
	flag.Parse()

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

	openRouterKey, ok := os.LookupEnv(envOpenRouterAPIKey)
	if !ok {
		return nil, errors.New("env OPENROUTER_API_KEY not found")
	}
	cfg.OpenRouterAPIKey = openRouterKey

	cfg.TelegramChatID = *telegramChatID
	if cfg.TelegramChatID == 0 {
		return nil, errors.New("flag -telegram-chat-id must be provided")
	}

	cfg.OpenRouterModel = strings.TrimSpace(*openRouterModel)
	if cfg.OpenRouterModel == "" {
		return nil, errors.New("flag -openrouter-model must be provided")
	}

	cfg.Tokens = splitCSV(*tokens)
	if len(cfg.Tokens) == 0 {
		return nil, errors.New("flag -tokens must be provided")
	}

	cfg.Protocols = splitCSV(*protocols)
	if len(cfg.Protocols) == 0 {
		return nil, errors.New("flag -protocols must be provided")
	}

	return cfg, nil
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		result = append(result, value)
	}

	return result
}
