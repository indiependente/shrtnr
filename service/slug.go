//go:generate mockgen -package service -source=slug.go -destination slug_mock.go

package service

import (
	"math/rand"
	"time"
)

// Slugger defines the behaviour of a component capable of returning a slug.
type Slugger interface {
	Slug() string
	Validate(string) bool
}

// FixedLenSlugger is a slugger implementation that returns slugs having fixed length.
type FixedLenSlugger struct {
	dictionary []rune
	length     int
}

// NewFixedLenSlugger returns a new FixedLenSlugger.
func NewFixedLenSlugger(l int) FixedLenSlugger {
	rand.Seed(time.Now().UnixNano())
	return FixedLenSlugger{
		dictionary: []rune("abcdefghijklmnopqrstuvwxyz"),
		length:     l,
	}
}

// Slug returns a slug having fixed length.
func (s FixedLenSlugger) Slug() string {
	b := make([]rune, s.length)
	for i := range b {
		b[i] = s.dictionary[rand.Intn(len(s.dictionary))]
	}
	return string(b)
}

func (s FixedLenSlugger) Validate(slug string) bool {
	if s.length == 0 {
		return false
	}
	if len(slug) != s.length {
		return false
	}
	for _, c := range slug {
		if !lettersMap[c] {
			return false
		}
	}
	return true
}
