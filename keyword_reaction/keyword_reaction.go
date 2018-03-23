package main

import (
	"fmt"
	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

// Config structure set by the bot api
var Config Configuration

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Reactions is a map[reaction][]keyword to use when message is posted
		Reactions map[string][]string `json:"reactions"`
	}

	// KeywordReaction implement bot.Feature
	KeywordReaction struct {
		reactions map[string][]string
	}
)

// NewFeature return interface bot.Feature used by the bot for run it
func NewFeature() bot.Feature {
	return &KeywordReaction{
		reactions: Config.Reactions,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *KeywordReaction) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

// Run the feature, triggered by event new message
func (f *KeywordReaction) Run(ctx *bot.Context) error {
	ir := slack.ItemRef{Channel: ctx.MsgEvent.Channel, Timestamp: ctx.MsgEvent.Timestamp}
	// add extra space for matching with single word
	sentence := fmt.Sprintf(" %s ", ctx.MsgEvent.Text)

	for reaction, keywords := range f.reactions {
		l := ctx.Log.WithField("reaction", reaction)
		l.Debug("search keywords")

		if contains(sentence, keywords) {
			if err := ctx.RTM.AddReaction(reaction, ir); err != nil {
				l.Error(err)
			}
		}
	}

	return nil
}

func contains(sentence string, keywords []string) bool {
	for _, keyword := range keywords {
		// add extra space for react only on full word
		if strings.Contains(sentence, fmt.Sprintf(" %s ", keyword)) {
			return true
		}
	}

	return false
}
