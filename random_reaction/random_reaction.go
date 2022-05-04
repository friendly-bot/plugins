package main

import (
	"math/rand"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type RandomReaction struct {
	chance    int
	reactions map[string]int
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewOnMessage(cfg *viper.Viper) (api.OnMessage, error) {
	rr := &RandomReaction{
		chance:    cfg.GetInt("chance"),
		reactions: make(map[string]int),
	}

	return rr, cfg.UnmarshalKey("reactions", &rr.reactions)
}

func (p RandomReaction) OnMessage(msg *slack.MessageEvent, ctx api.Context) error {
	if msg.Hidden {
		return nil
	}

	ir := slack.ItemRef{Channel: msg.Channel, Timestamp: msg.Timestamp}

	for reaction, probability := range p.reactions {
		r := rand.Intn(p.chance)

		ctx.Logger.WithFields(logrus.Fields{
			"reaction": reaction,
			"random":   r,
		}).Debug("roll")

		if r >= probability {
			continue
		}

		if err := ctx.RTM.AddReaction(reaction, ir); err != nil {
			ctx.Logger.WithField("reaction", reaction).Error(err)
		}
	}

	return nil
}
