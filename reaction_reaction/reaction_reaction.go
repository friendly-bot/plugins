package main

import (
	"regexp"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type ReactionReaction struct {
	reactions map[string]*regexp.Regexp
}

func NewReactionAdded(cfg *viper.Viper) (api.OnReactionAdded, error) {
	reactions := cfg.GetStringMapString("reactions")

	rr := &ReactionReaction{
		reactions: make(map[string]*regexp.Regexp, len(reactions)),
	}

	for emoji, regex := range reactions {
		rgx, err := regexp.Compile(regex)
		if err != nil {
			return nil, err
		}

		rr.reactions[emoji] = rgx
	}

	return rr, nil
}

func (p ReactionReaction) OnReactionAdded(ev *slack.ReactionAddedEvent, ctx api.Context) error {
	ir := slack.ItemRef{Channel: ev.Item.Channel, Timestamp: ev.Item.Timestamp}

	for reaction, regex := range p.reactions {
		l := ctx.Logger.WithField("reaction", reaction)
		l.WithField("regex", regex.String()).Debug("try regex")

		if regex.MatchString(ev.Reaction) {
			if err := ctx.RTM.AddReaction(reaction, ir); err != nil {
				l.Error(err)
			}
		}
	}

	return nil
}
