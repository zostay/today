package ost

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"time"

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

func (c *Client) TodayVerse(ctx context.Context, opts ...Option) (*Verse, error) {
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

	var verse Verse
	err = LoadVerseYaml(res.Body, &verse)
	if err != nil {
		return nil, err
	}

	return &verse, nil
}

func (c *Client) Today(ctx context.Context, opts ...Option) (string, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseText(ctx, verse.Reference)
}

func (c *Client) TodayHTML(ctx context.Context, opts ...Option) (template.HTML, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseHTML(ctx, verse.Reference)
}

func (c *Client) TodayPhoto(ctx context.Context, opts ...Option) (*Photo, error) {
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

	var ph Photo
	err = LoadPhotoYaml(res.Body, &ph)
	if err != nil {
		return nil, err
	}

	// TODO Pull the photo from OST instead when Pruned = false
	uph, err := c.PhotoService.Photo(ctx, ph.Link)
	if err != nil {
		return nil, err
	}

	ph.AddImage(photo.Original, uph.GetImage(photo.Original))

	return &ph, err
}
