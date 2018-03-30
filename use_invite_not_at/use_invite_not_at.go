package main

import (
	"regexp"
	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Message to send to the user
		Message       string   `json:"message"`
		ReactionsGood []string `json:"reactions_good"`
		ReactionsBad  []string `json:"reactions_bad"`
	}

	// UseInviteNotAt implement bot.Cron
	UseInviteNotAt struct {
		message       string
		reactionsGood []string
		reactionsBad  []string
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Cron used by the bot
func NewFeature(conf *Configuration) *UseInviteNotAt {
	return &UseInviteNotAt{
		message:       conf.Message,
		reactionsGood: conf.ReactionsGood,
		reactionsBad:  conf.ReactionsBad,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *UseInviteNotAt) Skip(ctx *bot.Context) (bool, string, error) {
	matched, err := regexp.MatchString("^<@\\w*>$", strings.Trim(ctx.MsgEvent.Text, " "))

	if err != nil {
		return true, "an error occurred", err
	}

	if !matched {
		return true, "is not an invitation", nil
	}

	members, err := ctx.Bot.GetListMembersByChannelID(ctx.MsgEvent.Channel, true)

	if err != nil {
		return true, "an error occurred", err
	}

	u := strings.Trim(ctx.MsgEvent.Text, "<>@")

	for _, member := range members {
		if member == u {
			return true, "user already in channel", nil
		}
	}

	return false, "", nil
}

// Run the cron
func (f *UseInviteNotAt) Run(ctx *bot.Context) error {
	_, ts, err := ctx.RTM.PostMessage(
		ctx.MsgEvent.Channel,
		f.message,
		slack.PostMessageParameters{
			ThreadTimestamp: ctx.MsgEvent.Timestamp,
			UnfurlLinks:     true,
			UnfurlMedia:     true,
		},
	)

	if err != nil {
		return err
	}

	badIR := slack.ItemRef{Channel: ctx.MsgEvent.Channel, Timestamp: ctx.MsgEvent.Timestamp}
	goodIR := slack.ItemRef{Channel: ctx.MsgEvent.Channel, Timestamp: ts}

	for _, reaction := range f.reactionsBad {
		if err := ctx.RTM.AddReaction(reaction, badIR); err != nil {
			return err
		}
	}

	for _, reaction := range f.reactionsGood {
		if err := ctx.RTM.AddReaction(reaction, goodIR); err != nil {
			return err
		}
	}

	return nil
}
