package chart

import (
	"testing"

	"github.com/crossdock/crossdock-go/require"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	f := NewFetcher("https://www.imdb.com/india/top-rated-indian-movies", 1)
	f.Init()
	err := f.Start()
	require.NoError(t, err)
	want := Chart{
		Movie{
			Title:       "Pather Panchali",
			ReleaseYear: 1955,
			Rating:      8.6,
			Summary:     "Impoverished priest Harihar Ray, dreaming of a better life for himself and his family, leaves his rural Bengal village in search of work.",
			Duration:    "2h 5min",
			Genre:       "Drama",
		},
	}
	// Should contain one movie
	require.Len(t, f.Chart, 1)
	// Elements should match
	assert.Equal(t, f.Chart, want)
}
