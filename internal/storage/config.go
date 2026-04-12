package storage

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Alexandersfg4/crypto-analyzer/internal/models"
)

func New(fileName string) *Config {
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

type Config struct {
	fileName string
}

func (c *Config) Save(newConfig *models.Config) error {
	data, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(c.fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
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
