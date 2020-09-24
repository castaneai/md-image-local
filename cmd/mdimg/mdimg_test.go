package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPattern(t *testing.T) {
	t.Run("ImageWithoutAlt", func(t *testing.T) {
		ss := pattern.FindAllString(`![](https://imgpath)`, -1)
		for _, s := range ss {
			md, err := extractImage(s)
			assert.NoError(t, err)
			assert.Equal(t, "", md.alt)
			assert.Equal(t, "https://imgpath", md.url)
		}

	})
	t.Run("ImageInLink", func(t *testing.T) {
		ss := pattern.FindAllString(`[![alt](https://imgpath)](https://link)`, -1)
		for _, s := range ss {
			md, err := extractImage(s)
			assert.NoError(t, err)
			assert.Equal(t, "alt", md.alt)
			assert.Equal(t, "https://imgpath", md.url)
		}
	})
}
