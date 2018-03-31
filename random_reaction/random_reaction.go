package main

import (
	"math/rand"
	"time"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Chance is used to calculate the chance to add reaction (Probability/Chance)
		Chance int `json:"chance"`

		// Reactions is map[reaction]probability, if random[0;Chance[ < Probability -> add reaction
		Reactions map[string]int `json:"reactions"`
	}

	// RandomReaction implement bot.feature
	RandomReaction struct {
		chance    int
		reactions map[string]int
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Feature used by the bot for run it
func NewFeature(conf *Configuration) bot.Feature {
	return &RandomReaction{
		chance:    conf.Chance,
		reactions: conf.Reactions,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *RandomReaction) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

// Run the feature, triggered by event new message
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
