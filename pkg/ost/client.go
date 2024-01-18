package ost

import (
	"html/template"
	"net/http"
	"net/url"
	"path"

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

func (c *Client) TodayVerse() (*Verse, error) {
	ru, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
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

func (c *Client) Today() (string, error) {
	verse, err := c.TodayVerse()
	if err != nil {
		return "", err
	}

	return c.Service.Verse(verse.Reference)
}

func (c *Client) TodayHTML() (template.HTML, error) {
	verse, err := c.TodayVerse()
	if err != nil {
		return "", err
	}

	return c.Service.VerseHTML(verse.Reference)
}
