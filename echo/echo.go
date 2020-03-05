package main

import (
	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type Echo struct {}

func NewOnMessage(_ *viper.Viper) (api.OnMessage, error) {
	return &Echo{}, nil
}

func (a Echo) OnMessage(msg *slack.MessageEvent, ctx api.Context) error {
	ctx.Logger.Info(msg.Text)

	ctx.Logger.WithField("value", ctx.Cache.Get("toto")).Info("before")
	ctx.Cache.Set("toto", "titi", 0)
	ctx.Logger.WithField("value", ctx.Cache.Get("toto")).Info("after")

	return nil
}
