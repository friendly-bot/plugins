package main

import (
	"regexp"
	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

var Config Configuration

type (
	Configuration struct {
		Message string `json:"message"`
	}

	UseInviteNotAt struct {
		message string
	}
)

func NewFeature(c *Configuration) *UseInviteNotAt {
	return &UseInviteNotAt{
		message: c.Message,
	}
}

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

func (f *UseInviteNotAt) Run(ctx *bot.Context) error {
	_, _, err := ctx.RTM.PostMessage(
		ctx.MsgEvent.Channel,
		f.message,
		slack.PostMessageParameters{
			ThreadTimestamp: ctx.MsgEvent.Timestamp,
			UnfurlLinks:     true,
			UnfurlMedia:     true,
		},
	)

	return err
}
