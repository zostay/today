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

const rootURL = `https://openscripture.today`

type Client struct {
	svc *text.Service
}

func New() (*Client, error) {
	res, err := esv.NewFromEnvironment()
	if err != nil {
		return nil, err
	}

	svc := text.NewService(res)
	return &Client{
		svc: svc,
	}, nil
}

func (c *Client) TodayVerse() (*Verse, error) {
	ru, err := url.Parse(rootURL)
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

	return c.svc.Verse(verse.Reference)
}

func (c *Client) TodayHTML() (template.HTML, error) {
	verse, err := c.TodayVerse()
	if err != nil {
		return "", err
	}

	return c.svc.VerseHTML(verse.Reference)
}
