package config

import (
	"github.com/asim/go-micro/v3/config"
	"github.com/asim/go-micro/v3/config/source"
	"github.com/asim/go-micro/v3/logger"
)

func NewConfigProvider(source source.Source) config.Config {
	cfg, _ := config.NewConfig()
	if err := cfg.Load(source); err != nil {
		logger.Log(logger.ErrorLevel, err)
	}

	return cfg
}
