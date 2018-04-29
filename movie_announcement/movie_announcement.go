package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"errors"

	"strings"

	"github.com/friendly-bot/slack-bot"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

const (
	// Endpoint for retrieve movie list
	Endpoint = "http://api.allocine.fr/rest/v3/movielist?"

	// ReleaseStateReprise -
	ReleaseStateReprise = "Reprise"

	//ReleaseVersionRestore -
	ReleaseVersionRestore = "Version restaurée"

	// LinkWeb is rel value for web page
	LinkWeb = "aco:web"

	// LinkShowtimes is rel value for showtime page
	LinkShowtimes = "aco:web_showtimes"

	// LinkTrailers is rel value for trailers
	LinkTrailers = "aco:web_trailers"
)

var colors = []string{"#FECC00", "#6197D0"}

type (
	// Configuration for the plugin, unmarshal by bot api
	Configuration struct {
		// Channel using for send message
		Channel string `json:"channel"`

		// PartnerKey of allocine api
		PartnerKey string `json:"partner_key"`

		// SecretKey of allocine api
		SecretKey string `json:"secret_key"`
	}

	// MovieAnnoucement implement bot.Cron
	MovieAnnoucement struct {
		channel    string
		partnerKey string
		secretKey  string
	}
)

type (
	// MovieList is struct return by allocine API
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

		// Statistics about the movie (vote, fan, rank, ...)
		Statistics Statistics `json:"statistics"`

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
		// Directors list separate by commas
		Directors string `json:"directors"`

		// Actors list separate by commas
		Actors string `json:"actors"`
	}

	// Statistics about the movie (vote, fan, rank, ...)
	Statistics struct {
		// TODO
	}

	// Link to another resource
	Link struct {
		// Rel name of the resource
		Rel string `json:"rel"`

		// Name readable for the resource
		Name string `json:"name"`

		// Href to the resource
		Href string `json:"href"`
	}

	// Links is alias for []Link
	Links []Link
)

// NewConfiguration return default configuration for this feature
func NewConfiguration() *Configuration {
	return &Configuration{}
}

// NewCron return interface bot.Cron used by the bot
func NewCron(conf *Configuration) bot.Cron {
	return &MovieAnnoucement{
		channel:    conf.Channel,
		partnerKey: conf.PartnerKey,
		secretKey:  conf.SecretKey,
	}
}

// Skip the run depend on the context, return bool (need to be skipped), string (reason of the skip), and an error if any
func (f *MovieAnnoucement) Skip(ctx *bot.Context) (bool, string, error) {
	return false, "", nil
}

// Run the cron
func (f *MovieAnnoucement) Run(ctx *bot.Context) error {
	movies, err := f.retrieveReleaseMoviesToday(ctx)

	if err != nil || len(movies) == 0 {
		return err
	}

	today := time.Now().Format("2006-01-02")

	var aa []slack.Attachment

	for i, movie := range movies {
		if movie.Release.ReleaseState.Value == ReleaseStateReprise || movie.Release.ReleaseVersion.Value == ReleaseVersionRestore {
			continue
		}

		if movie.Release.ReleaseDate != today {
			break
		}

		a := slack.Attachment{
			Title:    movie.Title,
			ImageURL: movie.Poster.Href,
			Color:    colors[i%2],
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
			// Footer: "TODO FOR RANK"
		}

		aa = append(aa, a)
	}

	_, _, err = ctx.RTM.PostMessage(f.channel, "Sorties de la semaine", slack.PostMessageParameters{Attachments: aa})

	return err
}

func (f *MovieAnnoucement) retrieveReleaseMoviesToday(ctx *bot.Context) ([]Movie, error) {
	c := http.Client{}

	sed := time.Now().Format("20060102")

	q := fmt.Sprintf(
		"partner=%s&count=50&filter=%s&page=1&order=datedesc&format=json&sed=%s",
		f.partnerKey,
		"nowshowing",
		sed,
	)

	s := sha1.Sum([]byte(f.secretKey + q))

	sig := url.PathEscape(base64.RawStdEncoding.EncodeToString(s[:]))

	r, err := c.Get(Endpoint + q + "&sig=" + sig + "%3D")

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	var movies MovieList
	var bs []byte

	if bs, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		ctx.Log.WithFields(logrus.Fields{"status_code": r.StatusCode, "response": string(bs)}).Warn("allocine api error")
		return nil, errors.New("allocine api error")
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
