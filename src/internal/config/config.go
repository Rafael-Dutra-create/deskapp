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
}

func NewConfig() *Config {
	cfg := &Config{
		DBDSN: os.Getenv("DATABASE_URL"),
		Port: os.Getenv("PORT"),
		mode: os.Getenv("MODE"),
	}
	cfg.Setup()

	return cfg
}


func (c *Config) Setup() {
	if c.Port == "" {
		c.Port = "8006"
	}

	if c.DBDSN == "" {
        c.DBDSN = "postgresql://root:gjdpakfnFpL0z1Kps@db01.controllab.com:5432/dados_prod_rc39?sslmode=disable"
    }
}

func (c *Config) GetMode() utils.MODE {
	switch strings.ToUpper(c.mode) {
		case "DEBUG": return utils.DEBUG
		default: return utils.RELEASE
	}
}