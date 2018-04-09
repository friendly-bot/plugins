package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

const (
	SubTypeMessageDeleted = "message_deleted"
	SubTypeMessageChanged = "message_changed"

	keyMessage = "message_%s"
)

var ErrEmptyKey = errors.New("empty message")

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Message to send to the user, suffixed by the deleted message
		Message string `json:"message"`
	}

	// NotifyDeletedMessage implement bot.Feature
	NotifyDeletedMessage struct {
		message string
	}

	// Message cached in redis
	Message struct {
		Text string `json:"text"`
		User string `json:"user"`
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Cron used by the bot
func NewFeature(conf *Configuration) *NotifyDeletedMessage {
	return &NotifyDeletedMessage{message: conf.Message}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *NotifyDeletedMessage) Skip(ctx *bot.Context) (bool, string, error) {
	if ctx.MsgEvent.SubType == SubTypeMessageDeleted {
		return false, "", nil
	}

	if ctx.MsgEvent.Text != "" {
		return false, "", nil
	}

	if ctx.MsgEvent.SubMessage != nil && ctx.MsgEvent.SubMessage.Text != "" {
		return false, "", nil
	}

	return true, "neither a new message or deleted message", nil
}

// Run the feature
func (f *NotifyDeletedMessage) Run(ctx *bot.Context) error {
	if ctx.MsgEvent.SubType == SubTypeMessageDeleted {
		return f.notify(ctx)
	}

	return f.log(ctx)
}

func (f *NotifyDeletedMessage) notify(ctx *bot.Context) error {
	k := fmt.Sprintf(keyMessage, ctx.MsgEvent.DeletedTimestamp)

	ctx.Log.WithField("key", k).Info("notify")

	if !ctx.Bot.IsKeyExist(k) {
		ctx.Log.WithField("key", k).Warn("key not exist")
		return nil
	}

	var m Message

	s := ctx.Bot.GetCacheString(k)

	if s == "" {
		return ErrEmptyKey
	}

	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return err
	}

	_, _, err := ctx.RTM.PostMessage(m.User, fmt.Sprintf("%s\n>%s", f.message, m.Text), slack.PostMessageParameters{AsUser: true})

	return err
}

func (f *NotifyDeletedMessage) log(ctx *bot.Context) error {
	m := Message{
		Text: ctx.MsgEvent.Text,
		User: ctx.MsgEvent.User,
	}

	k := fmt.Sprintf(keyMessage, ctx.MsgEvent.Timestamp)

	ctx.Log.WithField("key", k).Info("log")

	if ctx.MsgEvent.SubType == SubTypeMessageChanged {
		m.User = ctx.MsgEvent.Msg.User
		m.Text = ctx.MsgEvent.Msg.Text

		k = fmt.Sprintf(keyMessage, ctx.MsgEvent.Msg.Timestamp)
	}

	bs, err := json.Marshal(m)

	if err != nil {
		return err
	}

	return ctx.Bot.SetCacheExpire(k, bs, time.Hour)
}
