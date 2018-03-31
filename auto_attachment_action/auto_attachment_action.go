package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/friendly-bot/slack-bot"
)

const endpoint = "https://%s.slack.com/api/chat.attachmentAction"

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		TextAttachment string `json:"text_attachment"`
		Token          string `json:"token"`
		Team           string `json:"team"`
	}

	// AutoAttachmentAction implement bot.Feature
	AutoAttachmentAction struct {
		textAttachment string
		token          string
		team           string
	}

	payload struct {
		Actions      []interface{} `json:"actions"`
		AttachmentID string        `json:"attachment_id"`
		CallbackID   string        `json:"callback_id"`
		ChannelID    string        `json:"channel_id"`
		MessageTS    string        `json:"message_ts"`
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Feature used by the bot for run it
func NewFeature(conf *Configuration) bot.Feature {
	return &AutoAttachmentAction{
		textAttachment: conf.TextAttachment,
		token:          conf.Token,
		team:           conf.Team,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *AutoAttachmentAction) Skip(ctx *bot.Context) (bool, string, error) {
	if len(ctx.MsgEvent.Attachments) == 0 {
		return true, "no attachment", nil
	}

	if !strings.Contains(ctx.MsgEvent.Attachments[0].Text, f.textAttachment) {
		return true, "attachment text doesn't match", nil
	}

	if len(ctx.MsgEvent.Attachments[0].Actions) == 0 {
		return true, "no action available", nil
	}

	return false, "", nil
}

// Run the feature, triggered by event new message
func (f *AutoAttachmentAction) Run(ctx *bot.Context) error {
	p, err := json.Marshal(payload{
		Actions:      []interface{}{ctx.MsgEvent.Attachments[0].Actions[0]},
		AttachmentID: strconv.Itoa(ctx.MsgEvent.Attachments[0].ID),
		CallbackID:   ctx.MsgEvent.Attachments[0].CallbackID,
		ChannelID:    ctx.MsgEvent.Channel,
		MessageTS:    ctx.MsgEvent.Timestamp,
	})

	if err != nil {
		return err
	}

	values := url.Values{
		"service_id": {ctx.MsgEvent.BotID},
		"token":      {f.token},
		"payload":    {fmt.Sprintf("%s", p)},
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, f.team), strings.NewReader(values.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	c := http.Client{}

	_, err = c.Do(req)

	return err
}
