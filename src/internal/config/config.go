package config

import (
	"deskapp/src/internal/utils"
	"os"
	"strings"
)

type Config struct {
	DBDSN string
	Port string
	mode string
	server string
}

func NewConfig() *Config {
	cfg := &Config{
		DBDSN: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
		mode: os.Getenv("MODE"),
		server: os.Getenv("SERVER"),
	}
	cfg.Setup()

	return cfg
}


func (c *Config) Setup() {
	logger := utils.NewLogger()
	if c.Port == "" {
		c.Port = "8006"
	}

	if c.DBDSN == "" {
        logger.Warning("Nenhum Banco conectado!")
    }
}

func (c *Config) GetMode() utils.MODE {
	switch strings.ToUpper(c.mode) {
		case "DEBUG": return utils.DEBUG
		default: return utils.RELEASE
	}
}

func (c *Config) IsServer() bool {
	return strings.ToUpper(c.server) == "TRUE"
}