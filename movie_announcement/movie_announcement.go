package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cyruzin/golang-tmdb"
	"github.com/friendly-bot/friendly-bot/api"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

var colors = []string{"#FECC00", "#6197D0"}

type MovieAnnouncement struct {
	channel    string
	client     *tmdb.Client
	region     string
	language   string
	sizePoster string
}

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	c, err := tmdb.Init(cfg.GetString("api_key"))
	if err != nil {
		return nil, err
	}

	c.SetClientAutoRetry()

	return &MovieAnnouncement{
		channel:    cfg.GetString("channel"),
		client:     c,
		region:     cfg.GetString("region"),
		language:   cfg.GetString("language"),
		sizePoster: cfg.GetString("size_poster"),
	}, err
}

func (p MovieAnnouncement) Run(ctx api.Context) error {
	movies, err := p.retrieveMoviesNowPlaying(ctx)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")

	var aa []slack.Attachment
	var count int

	for _, movie := range movies.Results {
		// fuck movie without overview
		if movie.ReleaseDate != today || movie.Overview == "" {
			continue
		}

		ctx.Logger.WithField("id", movie.ID).Debug("request movie details")
		details, err := p.client.GetMovieDetails(int(movie.ID), map[string]string{"language": p.language})
		if err != nil {
			return fmt.Errorf("GetMovieDetails: (id: %d) %w", movie.ID, err)
		}

		var genres string
		for _, g := range details.Genres {
			genres = fmt.Sprintf("%s, %s", genres, g.Name)
		}

		genres = strings.TrimPrefix(genres, ", ")

		a := slack.Attachment{
			Title:    movie.Title,
			Text:     fmt.Sprintf("%s\n:invisible:", movie.Overview),
			ImageURL: tmdb.GetImageURL(movie.PosterPath, p.sizePoster),
			Color:    colors[count%len(colors)],
			Fields: []slack.AttachmentField{
				{Title: "Genres", Value: genres, Short: true},
			},
		}

		if details.Runtime > 0 {
			a.Fields = append(a.Fields, slack.AttachmentField{Title: "Duration", Value: fmt.Sprintf("%d min", details.Runtime), Short: true})
		}

		aa = append(aa, a)
		count++

	}

	if len(aa) > 0 {
		_, _, err = ctx.RTM.PostMessage(p.channel,
			slack.MsgOptionText("Sorties de la semaine", false),
			slack.MsgOptionAttachments(aa...),
		)
	}

	return err
}

func (p MovieAnnouncement) retrieveMoviesNowPlaying(ctx api.Context) (*tmdb.MovieNowPlayingResults, error) {
	opts := map[string]string{"region": p.region, "language": p.language}
	movies := &tmdb.MovieNowPlayingResults{}

	for page, totalPages := 1, 1; page <= totalPages; page++ {
		ctx.Logger.WithField("page", page).Debug("request now playing")

		opts["page"] = strconv.Itoa(page)
		ms, err := p.client.GetMovieNowPlaying(opts)
		if err != nil {
			return nil, fmt.Errorf("GetMovieNowPlaying (page: %d): %w", page, err)
		}
		totalPages = int(ms.TotalPages)

		movies.Results = append(movies.Results, ms.Results...)
	}

	return movies, nil
}
