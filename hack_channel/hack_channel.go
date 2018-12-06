package main

import (
	"fmt"
	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
)

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Message send by the bot
		Message string `json:"message"`

		// ChannelKeyword trigger for @channel
		ChannelKeyword []string `json:"channel_keyword"`

		// HereKeyword trigger for @everyone
		HereKeyword []string `json:"here_keyword"`

		// OnPublic disable or enable mention channel or everyone on public channel
		OnPublic bool `json:"on_public"`

		// EnabledChannel is a list of channel where feature is enable
		EnabledChannel []string `json:"enabled_channel"`

		// EnabledChannel is a list of channel where feature is disable (override enabled_channel)
		DisabledChannel []string `json:"disabled_channel"`
	}

	// HackChannel implement bot.Feature
	HackChannel struct {
		message         string
		channelKeyword  []string
		hereKeyword     []string
		onPublic        bool
		enabledChannel  []string
		disabledChannel []string
	}
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewFeature return interface bot.Feature used by the bot for run it
func NewFeature(conf *Configuration) bot.Feature {
	return &HackChannel{
		message:         conf.Message,
		channelKeyword:  conf.ChannelKeyword,
		hereKeyword:     conf.HereKeyword,
		onPublic:        conf.OnPublic,
		enabledChannel:  conf.EnabledChannel,
		disabledChannel: conf.DisabledChannel,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *HackChannel) Skip(ctx *bot.Context) (bool, string, error) {
	sentence := fmt.Sprintf(" %s ", ctx.MsgEvent.Text)
	sentence = strings.Replace(sentence, "\n", " ", -1)

	if !contains(sentence, append(f.channelKeyword, f.hereKeyword...)) {
		return true, fmt.Sprintf("no %s or %s", f.channelKeyword, f.hereKeyword), nil
	}

	cID := ctx.MsgEvent.Channel

	if f.IsDisabled(cID) {
		return true, fmt.Sprintf("%s is disabled", cID), nil
	}

	if f.IsEnabled(cID) {
		return false, "", nil
	}

	isPublic := false

	// if n != cID, name was found, so is a public channel (otherwise is group, so private)
	if n := ctx.Bot.GetChannelNameByID(cID, false); n != cID {
		isPublic = true
	}

	if isPublic && !f.onPublic {
		return true, "public channel is not enabled", nil
	}

	// chan public - public_enable - enable - !disable -> true
	// chan public - public_enable - enable - disable -> false
	// chan public - public_enable - !enable - disable -> false
	// chan public - public_enable - !enable - !disable -> true

	// chan public - !public_enable - enable - !disable -> true
	// chan public - !public_enable - enable - disable -> false
	// chan public - !public_enable - !enable - disable -> false
	// chan public - !public_enable - !enable - !disable -> false

	// chan private - enable - !disable -> true
	// chan private - enable - disable -> false
	// chan private - !enable - disable -> false
	// chan private - !enable - !disable -> true

	return false, "", nil
}

// Run the feature, triggered by event new message
func (f *HackChannel) Run(ctx *bot.Context) error {
	// add extra space for matching with single word
	sentence := fmt.Sprintf(" %s ", ctx.MsgEvent.Text)
	sentence = strings.Replace(sentence, "\n", " ", -1)

	t := "<!here>"

	if contains(sentence, f.channelKeyword) {
		t = "<!channel>"
	}

	a := slack.Attachment{Text: fmt.Sprintf(f.message, t)}

	_, _, e := ctx.RTM.PostMessage(
		ctx.MsgEvent.Channel,
		slack.MsgOptionAttachments(a),
		//slack.MsgOptionTS(ctx.MsgEvent.Timestamp),
	)

	return e
}

func contains(sentence string, keywords []string) bool {
	for _, keyword := range keywords {
		// add extra space for react only on full word
		if strings.Contains(sentence, fmt.Sprintf(" %s ", keyword)) {
			return true
		}
	}

	return false
}

func (f *HackChannel) IsDisabled(chanID string) bool {
	for _, d := range f.disabledChannel {
		if d == chanID {
			return true
		}
	}

	return false
}

func (f *HackChannel) IsEnabled(chanID string) bool {
	for _, e := range f.enabledChannel {
		if e == chanID {
			return true
		}
	}

	return false
}
