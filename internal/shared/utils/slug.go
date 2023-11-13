package utils

import "github.com/gosimple/slug"

func GenerateSlug(text string) string {
	s := slug.Make(text)
	return s
}
