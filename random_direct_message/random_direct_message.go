package main

import (
	"math/rand"
	"time"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

var Config Configuration

type (
	Configuration struct {
		Probability int      `json:"probability"`
		Chance      int      `json:"chance"`
		Messages    []string `json:"messages"`
	}

	RandomDirectMessage struct {
		probability int
		chance      int
		messages    []string
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewCron(probability, chance int, messages []string) *RandomDirectMessage {
	return &RandomDirectMessage{
		probability: probability,
		chance:      chance,
		messages:    messages,
	}
}

func (f *RandomDirectMessage) Skip(ctx *bot.Context) (bool, string, error) {
	if r := rand.Intn(f.chance); r >= f.probability {
		ctx.Log.WithField("random", r).Debug("roll")
		return true, "random greater than probability", nil
	}

	return false, "", nil
}

func (f *RandomDirectMessage) Run(ctx *bot.Context) error {
	users, err := ctx.Bot.GetActiveUsers()

	if err != nil {
		return err
	}

	ctx.Log.WithFields(logrus.Fields{
		"active_users":  len(users),
		"count_message": len(f.messages),
	}).Debug("count")

	if len(f.messages) == 0 || len(users) == 0 {
		ctx.Log.Warn("empty list")
		return nil
	}

	user := users[rand.Intn(len(users))]
	message := f.messages[rand.Intn(len(f.messages))]

	ctx.Log.WithFields(logrus.Fields{
		"send_to": user.Name,
		"message": message,
	}).Info("send direct message")

	_, _, err = ctx.RTM.PostMessage(user.ID, message, slack.PostMessageParameters{AsUser: true})

	return err
}
