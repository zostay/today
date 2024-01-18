package ost

import (
	"html/template"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
)

const DefaultBaseURL = `https://openscripture.today`

type Client struct {
	Client  *http.Client
	Service *text.Service
	BaseURL string
}

func New() (*Client, error) {
	res, err := esv.NewFromEnvironment()
	if err != nil {
		return nil, err
	}

	svc := text.NewService(res)
	return &Client{
		Client:  http.DefaultClient,
		BaseURL: DefaultBaseURL,
		Service: svc,
	}, nil
}

type options struct {
	onTime time.Time
}

type Option func(*options)

func On(t time.Time) Option {
	return func(o *options) {
		o.onTime = t
	}
}

func processOptions(opts []Option) *options {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (c *Client) TodayVerse(opts ...Option) (*Verse, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	o := processOptions(opts)
	if !o.onTime.IsZero() {
		datePath := o.onTime.Format("2006/01/02")
		ru.Path = path.Join(ru.Path, "verses", datePath)
	}

	ru.Path = path.Join(ru.Path, "verse.yaml")
	verseYamlUrl := ru.String()
	res, err := http.DefaultClient.Get(verseYamlUrl)
	if err != nil {
		return nil, err
	}

	var verse Verse
	dec := yaml.NewDecoder(res.Body)
	err = dec.Decode(&verse)
	return &verse, err
}

func (c *Client) Today(opts ...Option) (string, error) {
	verse, err := c.TodayVerse(opts...)
	if err != nil {
		return "", err
	}

	return c.Service.Verse(verse.Reference)
}

func (c *Client) TodayHTML(opts ...Option) (template.HTML, error) {
	verse, err := c.TodayVerse(opts...)
	if err != nil {
		return "", err
	}

	return c.Service.VerseHTML(verse.Reference)
}
