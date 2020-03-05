package main

import (
	"github.com/friendly-bot/friendly-bot/api"
	"github.com/spf13/viper"
)

type Ping struct {
	message string
}

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	return &Ping{message: cfg.GetString("message")}, nil
}

func (p Ping) Run(ctx api.Context) error {
	ctx.Logger.Info(p.message)

	return nil
}
