package main

import (
	"regexp"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type KeywordReaction struct {
	reactions map[string]*regexp.Regexp
}

func NewOnMessage(cfg *viper.Viper) (api.OnMessage, error) {
	reactions := cfg.GetStringMapString("reactions")

	kr := &KeywordReaction{
		reactions: make(map[string]*regexp.Regexp, len(reactions)),
	}

	for emoji, regex := range reactions {
		rgx, err := regexp.Compile(regex)
		if err != nil {
			return nil, err
		}

		kr.reactions[emoji] = rgx
	}

	return kr, nil
}

func (p KeywordReaction) OnMessage(msg *slack.MessageEvent, ctx api.Context) error {
	ir := slack.ItemRef{Channel: msg.Channel, Timestamp: msg.Timestamp}

	for reaction, regex := range p.reactions {
		l := ctx.Logger.WithField("reaction", reaction)
		l.WithField("regex", regex.String()).Debug("try regex")

		if regex.MatchString(msg.Text) {
			if err := ctx.RTM.AddReaction(reaction, ir); err != nil {
				l.Error(err)
			}
		}
	}

	return nil
}
