package ost

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/zostay/today/pkg/image"
	"github.com/zostay/today/pkg/image/unsplash"
	"github.com/zostay/today/pkg/text"
	"github.com/zostay/today/pkg/text/esv"
)

const DefaultBaseURL = `https://openscripture.today`

type Client struct {
	Client       *http.Client
	TextService  *text.Service
	PhotoService *image.Service
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
	imgSvc := image.NewService(src)

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

	return c.TextService.Verse(verse.Reference)
}

func (c *Client) TodayHTML(opts ...Option) (template.HTML, error) {
	verse, err := c.TodayVerse(opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseHTML(verse.Reference)
}

func (c *Client) TodayPhoto(opts ...Option) (*image.PhotoInfo, error) {
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
	res, err := http.DefaultClient.Get(photoYamlUrl)
	if err != nil {
		return nil, err
	}

	var photo image.PhotoInfo
	dec := yaml.NewDecoder(res.Body)
	err = dec.Decode(&photo)
	return &photo, err
}

func (c *Client) TodayPhotoWithDownload(
	ctx context.Context,
	opts ...Option,
) (*image.PhotoInfo, error) {
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
	res, err := http.DefaultClient.Get(photoYamlUrl)
	if err != nil {
		return nil, err
	}

	var photo image.PhotoInfo
	dec := yaml.NewDecoder(res.Body)
	err = dec.Decode(&photo)

	err = c.PhotoService.Download(ctx, &photo)
	if err != nil {
		return nil, err
	}

	return &photo, err
}
