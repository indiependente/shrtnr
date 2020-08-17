package models

// URLShortened represents the short version of a URL.
type URLShortened struct {
	URL  string `json:"url"`
	Slug string `json:"slug"`
	Hits int    `json:"hits"`
}
