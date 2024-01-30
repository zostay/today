package ost

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/photo"
	"github.com/zostay/today/pkg/photo/unsplash"
	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
)

const DefaultBaseURL = `https://openscripture.today`

type Client struct {
	Client       *http.Client
	TextService  *text.Service
	PhotoService *photo.Service
	BaseURL      string
}

func New(ctx context.Context) (*Client, error) {
	res, err := esv.NewFromEnvironment()
	if err != nil {
		return nil, err
	}
	txtSvc := text.NewService(res)

	src, err := unsplash.NewFromEnvironment(ctx)
	if err != nil {
		return nil, err
	}
	imgSvc := photo.NewService(src)

	return &Client{
		Client:       http.DefaultClient,
		BaseURL:      DefaultBaseURL,
		TextService:  txtSvc,
		PhotoService: imgSvc,
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

func (c *Client) TodayVerse(ctx context.Context, opts ...Option) (*text.Verse, error) {
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

	r, err := http.NewRequest("GET", verseYamlUrl, nil)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	var verse text.Verse
	dec := yaml.NewDecoder(res.Body)
	err = dec.Decode(&verse)
	return &verse, err
}

func (c *Client) Today(ctx context.Context, opts ...Option) (string, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.Verse(ctx, verse.Reference)
}

func (c *Client) TodayHTML(ctx context.Context, opts ...Option) (template.HTML, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseHTML(ctx, verse.Reference)
}

func (c *Client) TodayPhoto(ctx context.Context, opts ...Option) (*photo.Info, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	o := processOptions(opts)
	if !o.onTime.IsZero() {
		datePath := o.onTime.Format("2006/01/02")
		ru.Path = path.Join(ru.Path, "verses", datePath)
	}

	ru.Path = path.Join(ru.Path, "photo.yaml")
	photoYamlUrl := ru.String()

	r, err := http.NewRequest("GET", photoYamlUrl, nil)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	var photo photo.Info
	dec := yaml.NewDecoder(res.Body)
	err = dec.Decode(&photo.Meta)

	photo.Key, _ = c.PhotoService.CacheKey(photo.Meta.Link)
	return &photo, err
}
