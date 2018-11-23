package main

import (
	"time"
	"fmt"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		Day     int    `json:"day"`
		Month   int    `json:"month"`
		User    string `json:"user"`
		Message string `json:"message"`
	}

	// HappyBirthday implement bot.Cron
	HappyBirthday struct {
		Day     int
		Month   int
		User    string
		Message string
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Cron used by the bot for run it
func NewCron(conf *Configuration) bot.Cron {
	return &HappyBirthday{
		Day:     conf.Day,
		Month:   conf.Month,
		User:    conf.User,
		Message: conf.Message,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *HappyBirthday) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

// Run the feature, triggered by event new message
func (f *HappyBirthday) Run(ctx *bot.Context) error {
	n := time.Now()

	yearLastBirthday := n.Year()

	if (n.Month() < time.Month(f.Month)) ||
		(n.Month() == time.Month(f.Month) && n.Day() < f.Day) {
		yearLastBirthday--
	}

	lastBirthday := time.Date(yearLastBirthday, time.Month(f.Month), f.Day, 0, 0, 0, 0, n.Location())

	days := int(n.Sub(lastBirthday).Hours() / 24)

	m := f.Message
	if days > 0 {
		m = fmt.Sprintf("%s +%d", m, days)
	}

	_, _, err := ctx.RTM.PostMessage(f.User, slack.MsgOptionText(m, false), slack.MsgOptionAsUser(true))

	return err
}
