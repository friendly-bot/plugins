package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

const endpoint = "https://%s.slack.com/api/chat.attachmentAction"

type (
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

func NewOnMessage(cfg *viper.Viper) (api.OnMessage, error) {
	return &AutoAttachmentAction{
		textAttachment: cfg.GetString("text_attachment"),
		token:          cfg.GetString("token"),
		team:           cfg.GetString("team"),
	}, nil
}

func (a AutoAttachmentAction) OnMessage(msg *slack.MessageEvent, ctx api.Context) error {
	if len(msg.Attachments) == 0 ||
		!strings.Contains(msg.Attachments[0].Text, a.textAttachment) ||
		len(msg.Attachments[0].Actions) == 0 {

		ctx.Logger.Debug("requirement missing")
		return nil
	}

	p, err := json.Marshal(payload{
		Actions:      []interface{}{msg.Attachments[0].Actions[0]},
		AttachmentID: strconv.Itoa(msg.Attachments[0].ID),
		CallbackID:   msg.Attachments[0].CallbackID,
		ChannelID:    msg.Channel,
		MessageTS:    msg.Timestamp,
	})

	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	values := url.Values{
		"service_id": {msg.BotID},
		"token":      {a.token},
		"payload":    {fmt.Sprintf("%s", p)},
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, a.team), strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := http.Client{
		Timeout: time.Second * 10,
	}

	_, err = c.Do(req)

	return err
}
