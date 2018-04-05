package main

import (
	"math/rand"
	"time"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Chance is used to calculate the chance to run this feature (Probability/Chance)
		Probability int `json:"probability"`

		// Chance is used to calculate the chance to run this feature (Probability/Chance)
		Chance int `json:"chance"`

		// Messages can be send by this feature
		Messages []string `json:"messages"`

		// TalkAfter x hour, for avoid spam
		TalkAfter int `json:"talk_after"`
	}

	// RandomDirectMessage implement bot.Cron
	RandomDirectMessage struct {
		probability int
		chance      int
		messages    []string
		talkAfter   time.Duration
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewCron return interface bot.Cron used by the bot
func NewCron(conf *Configuration) bot.Cron {
	return &RandomDirectMessage{
		probability: conf.Probability,
		chance:      conf.Chance,
		messages:    conf.Messages,
		talkAfter:   time.Hour * time.Duration(conf.TalkAfter),
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *RandomDirectMessage) Skip(ctx *bot.Context) (bool, string, error) {
	if r := rand.Intn(f.chance); r >= f.probability {
		ctx.Log.WithField("random", r).Debug("roll")
		return true, "random greater than probability", nil
	}

	if len(f.messages) == 0 {
		return true, "no message", nil
	}


	return false, "", nil
}

// Run the cron
func (f *RandomDirectMessage) Run(ctx *bot.Context) error {
	user, err := ctx.Bot.GetUserToTalk()

	if err != nil {
		return err
	}

	message := f.messages[rand.Intn(len(f.messages))]

	ctx.Log.WithFields(logrus.Fields{
		"send_to": user.Name,
		"message": message,
	}).Info("send direct message")

	_, _, err = ctx.RTM.PostMessage(user.ID, message, slack.PostMessageParameters{AsUser: true})

	if err == nil {
		return ctx.Bot.StopTalkToUser(user.ID, f.talkAfter)
	}

	return err
}
