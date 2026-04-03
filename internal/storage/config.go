package storage

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

type Config struct {
	fileName string
}

func NewConfig(fileName string) *Config {
	_, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.WriteFile(fileName, []byte("{}"), 0644)
		} else {
			panic(err)
		}
	}

	return &Config{fileName: fileName}
}

func (c *Config) SaveTokens(tokens []string) (models.Config, error) {
	config, err := c.Read()
	if err != nil {
		return models.Config{}, err
	}
	config.Tokens = tokens

	data, err := json.Marshal(config)
	if err != nil {
		return models.Config{}, err
	}
	err = os.WriteFile(c.fileName, data, 0644)
	if err != nil {
		return models.Config{}, err
	}

	return *config, nil
}

func (c *Config) SaveProtocols(protocols []string) (models.Config, error) {
	config, err := c.Read()
	if err != nil {
		return models.Config{}, err
	}
	config.Protocols = protocols

	data, err := json.Marshal(config)
	if err != nil {
		return models.Config{}, err
	}
	err = os.WriteFile(c.fileName, data, 0644)
	if err != nil {
		return models.Config{}, err
	}

	return *config, nil
}

func (c *Config) SaveCronNextExecutionTime(t time.Time) (models.Config, error) {
	config, err := c.Read()
	if err != nil {
		return models.Config{}, err
	}
	config.CronNextExecutionTime = t

	data, err := json.Marshal(config)
	if err != nil {
		return models.Config{}, err
	}
	err = os.WriteFile(c.fileName, data, 0644)
	if err != nil {
		return models.Config{}, err
	}

	return *config, nil
}

func (c *Config) Read() (*models.Config, error) {
	data, err := os.ReadFile(c.fileName)
	if err != nil {
		return nil, err
	}
	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
