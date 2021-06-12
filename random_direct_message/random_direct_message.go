package main

import (
	"context"
	"fmt"
	"github.com/friendly-bot/friendly-bot/api"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"math/rand"
	"time"
)

const presenceActive = "active"

type RandomDirectMessage struct {
	messages  []string
	talkAfter time.Duration
}

func init() {
	rand.Seed(time.Now().Unix())
}

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	return &RandomDirectMessage{
		messages:  cfg.GetStringSlice("messages"),
		talkAfter: cfg.GetDuration("talk_after"),
	}, nil
}

func (p RandomDirectMessage) Run(ctx api.Context) (err error) {
	page := ctx.RTM.GetUsersPaginated(slack.GetUsersOptionPresence(true))
	users := make([]slack.User, 0, 200)

	for {
		page, err = page.Next(context.Background())
		if err != nil {
			break
		}

		for _, u := range page.Users {
			ctx.Logger.WithFields(logrus.Fields{
				"bot":              u.IsBot,
				"restricted":       u.IsRestricted,
				"ultra_restricted": u.IsUltraRestricted,
				"app_user":         u.IsAppUser,
				"user":             u.Name,
				"stranger":         u.IsStranger,
				"deleted":          u.Deleted,
				"id":               u.ID,
			}).Debug("user")

			if u.ID == "" || u.Deleted || u.IsStranger || u.IsBot || u.IsRestricted || u.IsUltraRestricted || u.IsAppUser || ctx.Cache.Exist(u.ID) {
				continue
			}

			users = append(users, u)
		}
	}

	ctx.Logger.WithField("count", len(users)).Info("eligible user")

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	var user slack.User

	for _, u := range users {
		p, err := ctx.RTM.GetUserPresence(u.ID)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{
				"context":   "get_user_presence",
				"user_name": u.Name,
				"user_id":   u.ID,
			}).Error(err)

			continue
		}

		if p.Presence != presenceActive {
			continue
		}

		user = u
		break
	}

	ctx.Cache.Set(user.ID, "1", p.talkAfter)

	message := p.messages[rand.Intn(len(p.messages))]

	ctx.Logger.WithFields(logrus.Fields{
		"send_to": user.Name,
		"message": message,
	}).Info("send direct message")

	_, _, err = ctx.RTM.PostMessage(user.ID,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		return fmt.Errorf("PostMessage: %w", err)
	}

	return nil
}
