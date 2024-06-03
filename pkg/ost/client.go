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

// DefaultBaseURL is the default base URL for the openscripture.today API.
const DefaultBaseURL = `https://openscripture.today`

// Client is a client for the openscripture.today API.
type Client struct {
	Client       *http.Client
	TextService  *text.Service
	PhotoService *photo.Service
	BaseURL      string
}

// New creates a new openscripture.today API client. The context is used only to
// help load related client objects from the local environment and is not
// stored.
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

type indexOptions struct {
	period string
}

// IndexOption is a functional option for the openscripture.today API client.
type IndexOption func(*indexOptions)

// ForAllTime is the option that selects all verses for all time. This is the
// default.
func ForAllTime() IndexOption {
	return func(o *indexOptions) {
		o.period = ""
	}
}

// ForYear is the option that selects all verses for a given year.
func ForYear(year string) IndexOption {
	return func(o *indexOptions) {
		o.period = year
	}
}

// ForMonth is the option that selects all verses for a given month.
func ForMonth(year, month string) IndexOption {
	return func(o *indexOptions) {
		o.period = path.Join(year, month)
	}
}

// makeIndexOptions processes the options for the index request.
func makeIndexOptions(opts []IndexOption) *indexOptions {
	o := &indexOptions{
		period: "",
	}
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// VerseIndex returns an index for some segment of verses. Use one of the
// ost.IndexOption options to select the period of verses to index.
func (c *Client) VerseIndex(ctx context.Context, opts ...IndexOption) (*Index, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	o := makeIndexOptions(opts)
	if o.period == "" {
		ru.Path = path.Join(ru.Path, "verses", "index.yaml")
	} else {
		ru.Path = path.Join(ru.Path, "verses", o.period, "index.yaml")
	}
	indexYamlUrl := ru.String()

	r, err := http.NewRequest("GET", indexYamlUrl, nil)
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	var index Index
	err = LoadIndexYaml(res.Body, &index)
	if err != nil {
		return nil, err
	}

	return &index, nil
}

type dayOptions struct {
	onTime time.Time
}

// DayOption is a functional option for the openscripture.today API client.
type DayOption func(*dayOptions)

// On sets the time to use for the request. Only the date part of the time is used.
func On(t time.Time) DayOption {
	return func(o *dayOptions) {
		o.onTime = t
	}
}

func processDayOptions(opts []DayOption) *dayOptions {
	o := &dayOptions{}
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// TodayVerse returns the entire verse object for a given day. If no time is
// provided via the ost.On option, the current day is used.
func (c *Client) TodayVerse(ctx context.Context, opts ...DayOption) (*Verse, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	o := processDayOptions(opts)
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

// Today returns the text of the verse for the current day. If no time is
// provided via the ost.On option, the current day is used.
func (c *Client) Today(ctx context.Context, opts ...DayOption) (string, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseText(ctx, verse.Reference)
}

// TodayHTML returns the HTML of the verse for the current day. If no time is
// provided via the ost.On option, the current day is used.
func (c *Client) TodayHTML(ctx context.Context, opts ...DayOption) (template.HTML, error) {
	verse, err := c.TodayVerse(ctx, opts...)
	if err != nil {
		return "", err
	}

	return c.TextService.VerseHTML(ctx, verse.Reference)
}

// TodayPhoto returns the photo for the current day. If no time is provided via
// the ost.On option, the current day is used.
func (c *Client) TodayPhoto(ctx context.Context, opts ...DayOption) (*Photo, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	o := processDayOptions(opts)
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
