package main

import (
	"math/rand"
	"time"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

var Config Configuration

type (
	Configuration struct {
		Chance    int            `json:"chance"`
		Reactions map[string]int `json:"reactions"`
	}

	RandomReaction struct {
		chance    int
		reactions map[string]int
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewFeature(c *Configuration) bot.Feature {
	return &RandomReaction{
		chance:    c.Chance,
		reactions: c.Reactions,
	}
}

func (f *RandomReaction) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

func (f *RandomReaction) Run(ctx *bot.Context) error {
	ir := slack.ItemRef{Channel: ctx.MsgEvent.Channel, Timestamp: ctx.MsgEvent.Timestamp}

	for reaction, probability := range f.reactions {
		r := rand.Intn(f.chance)

		ctx.Log.WithFields(logrus.Fields{
			"reaction": reaction,
			"random":   r,
		}).Debug("roll")

		if r >= probability {
			continue
		}

		if err := ctx.RTM.AddReaction(reaction, ir); err != nil {
			ctx.Log.WithField("reaction", reaction).Error(err)
		}
	}

	return nil
}
