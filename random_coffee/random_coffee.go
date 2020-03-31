package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type (
	RandomCoffee struct {
		groupOf          int
		channel          string
		maxNumberOfGroup int
		header           string
		footer           string
		separator        string
		prefix           string
	}

	group struct {
		users     []string
		separator string
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	return &RandomCoffee{
		groupOf:          cfg.GetInt("group_of"),
		channel:          cfg.GetString("channel"),
		maxNumberOfGroup: cfg.GetInt("max_number_of_group"),
		header:           cfg.GetString("header"),
		footer:           cfg.GetString("footer"),
		separator:        cfg.GetString("separator"),
		prefix:           cfg.GetString("prefix"),
	}, nil
}

func (p RandomCoffee) Run(ctx api.Context) error {
	// TODO handle cursor?
	users, _, err := ctx.RTM.GetUsersInConversation(&slack.GetUsersInConversationParameters{ChannelID: p.channel, Limit: 1000})
	if err != nil {
		return err
	}
	users = removeBot(ctx, users)

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	groups := make([]group, 0, p.maxNumberOfGroup)

	for i := 0; i < p.maxNumberOfGroup; i++ {
		if p.groupOf*i+i >= len(users) {
			ctx.Logger.Info("can't create more groups, not enough users")
			break
		}

		groups = append(groups, group{separator: p.separator, users: users[p.groupOf*i : p.groupOf*(i+1)]})
	}

	if len(groups) == 0 {
		ctx.Logger.Warn("zero group created")
		return nil
	}

	msg := p.header
	for _, g := range groups {
		msg = fmt.Sprintf("%s\n%s %s", msg, p.prefix, g)
	}
	msg = fmt.Sprintf("%s%s", msg, p.footer)

	_, _, err = ctx.RTM.PostMessage(p.channel, slack.MsgOptionText(msg, false))

	return err
}

func removeBot(ctx api.Context, users []string) []string {
	filtered := make([]string, 0, len(users))

	for _, u := range users {
		i, err := ctx.RTM.GetUserInfo(u)
		if err != nil {
			ctx.Logger.WithFields(logrus.Fields{"user": u, "context": "get_user_info"}).Error(err)
			continue
		}

		if !i.IsBot {
			filtered = append(filtered, u)
		}
	}

	return filtered
}

func (g group) String() string {
	for i := range g.users {
		g.users[i] = fmt.Sprintf("<@%s>", g.users[i])
	}

	return strings.Join(g.users, g.separator)
}
