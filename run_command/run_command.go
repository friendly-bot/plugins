package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/spf13/viper"
)

const endpoint = "https://slack.com/api/chat.command"

type RunCommand struct {
	token   string
	channel string
	command string
	text    string
}

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	return &RunCommand{
		token:   cfg.GetString("token"),
		channel: cfg.GetString("channel"),
		command: cfg.GetString("command"),
		text:    cfg.GetString("text"),
	}, nil
}

func (p RunCommand) Run(_ api.Context) error {
	values := url.Values{
		"token":   {p.token},
		"channel": {p.channel},
		"command": {p.command},
		"text":    {p.text},
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := http.Client{Timeout: 5 * time.Second}

	_, err = c.Do(req)

	return err
}
