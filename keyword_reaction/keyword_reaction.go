package main

import (
	"fmt"
	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

var Config Configuration

type (
	Configuration struct {
		Reactions map[string][]string `json:"reactions"`
	}

	KeywordReaction struct {
		reactions map[string][]string
	}
)

func NewFeature(c *Configuration) bot.Feature {
	return &KeywordReaction{
		reactions: c.Reactions,
	}
}

func (f *KeywordReaction) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

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
