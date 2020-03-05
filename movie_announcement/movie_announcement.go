package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/grokify/html-strip-tags-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	endpoint              = "http://api.allocine.fr/rest/v3/movielist?"
	ReleaseStateReprise   = "Reprise"
	ReleaseVersionRestore = "Version restaurée"
	LinkWeb               = "aco:web"
	LinkShowtimes         = "aco:web_showtimes"
	LinkTrailers          = "aco:web_trailers"

	proxy = "https://images.weserv.nl/?url=%s&h=%d"
)

var colors = []string{"#FECC00", "#6197D0"}

type MovieAnnouncement struct {
	channel     string
	partnerKey  string
	secretKey   string
	heightImage int
}

type (
	MovieList struct {
		// Feed content the response
		Feed struct {
			// Movies list
			Movies []Movie `json:"movie"`
		} `json:"feed"`
	}

	// Movie entity
	Movie struct {
		// Title of the movie
		Title string `json:"title"`

		// Genres of the movie
		Genres Genres `json:"genre"`

		// Release date and original release if any
		Release Release `json:"release"`

		// SynopsisShort short version of the synopsis
		SynopsisShort string `json:"synopsisShort"`

		// CastingShort short version of the casting
		CastingShort CastingShort `json:"castingShort"`

		// Poster information
		Poster Poster `json:"poster"`

		// Links to another resource linked with this movie
		Links Links `json:"link"`
	}

	// Genre entity
	Genre struct {
		// Name -
		Name string `json:"$"`
	}

	// Genres is alias for []Genre
	Genres []Genre

	// Poster information
	Poster struct {
		// Href to poster
		Href string `json:"href"`
	}

	// Release information
	Release struct {
		// ReleaseDate in theaters
		ReleaseDate string `json:"releaseDate"`

		// ReleaseState (if reprise)
		ReleaseState ReleaseState `json:"releaseState"`

		// ReleaseVersion (if restored)
		ReleaseVersion ReleaseVersion `json:"releaseVersion"`
	}

	// ReleaseState if is original release or reprise
	ReleaseState struct {
		// Value -
		Value string `json:"$"`
	}

	// ReleaseVersion if is original release or restored version
	ReleaseVersion struct {
		// Value -
		Value string `json:"$"`
	}

	// CastingShort embed directors and actors casting
	CastingShort struct {
		Directors string `json:"directors"`
		Actors    string `json:"actors"`
	}

	Link struct {
		Rel  string `json:"rel"`
		Name string `json:"name"`
		Href string `json:"href"`
	}

	Links []Link
)

func NewJob(cfg *viper.Viper) (api.Runner, error) {
	return &MovieAnnouncement{
		channel:     cfg.GetString("channel"),
		partnerKey:  cfg.GetString("partner_key"),
		secretKey:   cfg.GetString("secret_key"),
		heightImage: cfg.GetInt("height_image"),
	}, nil
}

func (p MovieAnnouncement) Run(ctx api.Context) error {
	movies, err := p.retrieveReleaseMoviesToday(ctx)

	if err != nil || len(movies) == 0 {
		return err
	}

	today := time.Now().Format("2006-01-02")

	var aa []slack.Attachment
	var count int

	for _, movie := range movies {
		if movie.Release.ReleaseState.Value == ReleaseStateReprise || movie.Release.ReleaseVersion.Value == ReleaseVersionRestore || movie.Release.ReleaseDate != today {
			continue
		}

		a := slack.Attachment{
			Title:    movie.Title,
			Text:     strip.StripTags(movie.SynopsisShort),
			ImageURL: generateURL(movie.Poster.Href, p.heightImage),
			Color:    colors[count%2],
			Fields: []slack.AttachmentField{
				{Title: "Directors", Value: movie.CastingShort.Directors, Short: true},
				{Title: "Genres", Value: movie.Genres.AllGenres(), Short: true},
				{Title: "Actors", Value: movie.CastingShort.Actors, Short: false},
			},
			Actions: []slack.AttachmentAction{
				{Text: "Allocine", Type: "button", URL: movie.Links.GetHrefFor(LinkWeb)},
				{Text: "Bande annonces", Type: "button", URL: movie.Links.GetHrefFor(LinkTrailers)},
				{Text: "Séances", Type: "button", URL: movie.Links.GetHrefFor(LinkShowtimes)},
			},
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

func (p MovieAnnouncement) retrieveReleaseMoviesToday(ctx api.Context) ([]Movie, error) {
	c := http.Client{
		Timeout: time.Second * 10,
	}

	sed := time.Now().Format("20060102")

	q := fmt.Sprintf(
		"partner=%s&count=50&filter=%s&page=1&order=datedesc&format=json&sed=%s",
		p.partnerKey, "nowshowing", sed,
	)

	s := sha1.Sum([]byte(p.secretKey + q))
	sig := url.PathEscape(base64.RawStdEncoding.EncodeToString(s[:]))
	r, err := c.Get(fmt.Sprintf("%s%s&sig=%s%%3D", endpoint, q, sig))

	if err != nil {
		return nil, err
	}

	defer func() { _ = r.Body.Close() }()

	var movies MovieList
	var bs []byte

	if bs, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		ctx.Logger.WithFields(logrus.Fields{"status_code": r.StatusCode, "response": string(bs)}).Warn("allocine api error")
		return nil, fmt.Errorf("allocine api error code %d", r.StatusCode)
	}

	if err = json.Unmarshal(bs, &movies); err != nil {
		return nil, err
	}

	return movies.Feed.Movies, nil
}

func (gs Genres) AllGenres() string {
	var genres string

	for _, g := range gs {
		genres = fmt.Sprintf("%s, %s", genres, g.Name)
	}

	return strings.TrimLeft(genres, ", ")
}

func (ls Links) GetHrefFor(rel string) string {
	for _, l := range ls {
		if l.Rel == rel {
			return l.Href
		}
	}

	return ""
}

func generateURL(url string, h int) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	return fmt.Sprintf(proxy, url, h)
}
