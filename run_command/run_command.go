package main

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/friendly-bot/slack-bot"
	"io/ioutil"
)

const endpoint = "https://slack.com/api/chat.command"

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		Token   string `json:"token"`
		Channel string `json:"channel"`
		Command string `json:"command"`
		Text    string `json:"text"`
	}

	// RunCommand implement bot.Feature
	RunCommand struct {
		token   string
		channel string
		command string
		text    string
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Feature used by the bot for run it
func NewCron(conf *Configuration) bot.Cron {
	return &RunCommand{
		token:   conf.Token,
		channel: conf.Channel,
		command: conf.Command,
		text:    conf.Text,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *RunCommand) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

// Run the feature, triggered by event new message
func (f *RunCommand) Run(ctx *bot.Context) error {
	values := url.Values{
		"token":   {f.token},
		"channel": {f.channel},
		"command": {f.command},
		"text":    {f.text},
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	c := http.Client{Timeout: 5 * time.Second}

	_, err = c.Do(req)

	return err
}
