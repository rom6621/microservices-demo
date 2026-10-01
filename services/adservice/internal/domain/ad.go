package domain

import "net/url"

type Ad struct {
	text        string
	redirectURL *url.URL
}

func NewAd(text string, redirectURL *url.URL) (*Ad, error) {
	return &Ad{
		text:        text,
		redirectURL: redirectURL,
	}, nil
}

func (a *Ad) Text() string          { return a.text }
func (a *Ad) RedirectURL() *url.URL { return a.redirectURL }
