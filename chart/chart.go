package chart

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Movie contains info about a movie.
type Movie struct {
	Title       string  `json:"title"`
	ReleaseYear int     `json:"movie_release_year"`
	Rating      float32 `json:"imdb_rating"`
	Summary     string  `json:"summary"`
	Duration    string  `json:"duration"`
	Genre       string  `json:"genre"`
}

// Chart is a collection of movies.
type Chart []Movie

// Fetcher is used to fetch a chart.
type Fetcher struct {
	c          *colly.Collector
	url        string
	itemsCount int
	Chart      Chart
}

// NewFetcher creates an instance of Fetcher.
func NewFetcher(url string, count int) *Fetcher {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
		colly.MaxDepth(3),
		colly.UserAgent("*"),
	)
	return &Fetcher{
		c:          c,
		url:        url,
		itemsCount: count,
		Chart:      make(Chart, 0, count),
	}
}

func (f *Fetcher) Init() {
	movieInfoCollector := f.c.Clone()
	moviesVisited := 0
	visitedURLs := map[string]struct{}{}

	f.c.OnHTML(".lister-list tr td a", func(e *colly.HTMLElement) {
		if moviesVisited < int(f.itemsCount) {
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))
			movieInfoCollector.Visit(nextPage)
			if _, ok := visitedURLs[nextPage]; !ok {
				moviesVisited++
			}
			visitedURLs[nextPage] = struct{}{}
		}
	})

	movieInfoCollector.OnHTML("div[id=main_top]", func(e *colly.HTMLElement) {
		titleAndYear := e.ChildText("div.title_wrapper > h1")
		summary := e.ChildText("div.summary_text")
		rating := e.ChildText("div.ratingValue > strong > span")
		subtext := strings.Split(e.ChildText("div.subtext"), "|")
		duration, genre := subtext[1], subtext[2]

		// Parse rating.
		var parsedRating float64
		var err error
		if parsedRating, err = strconv.ParseFloat(rating, 32); err != nil {
			return
		}

		// Parse release year.
		year, err := parseReleaseYear(titleAndYear)
		if err != nil {
			return
		}

		movie := &Movie{
			Title:       strings.Split(titleAndYear, "(")[0],
			ReleaseYear: year,
			Summary:     summary,
			Rating:      float32(parsedRating),
			Duration:    duration,
			Genre:       genre,
		}
		movie.sanitize()

		f.Chart = append(f.Chart, *movie)
	})
}

func (f *Fetcher) Start() error {
	return f.c.Visit(f.url)
}

// Removes leading/trailing spaces and special characters.
func (m *Movie) sanitize() {
	m.Title = strings.TrimSpace(m.Title)
	m.Summary = strings.TrimSpace(m.Summary)
	m.Duration = strings.TrimSpace(m.Duration)

	// Sanitize genres.
	genres := strings.Split(m.Genre, ",")
	for i, g := range genres {
		genres[i] = strings.TrimSpace(g)
	}
	m.Genre = strings.Join(genres, ",")
}

func parseReleaseYear(x string) (int, error) {
	yearString := strings.Trim(strings.Split(x, "(")[1], ")")
	year, err := strconv.ParseInt(yearString, 10, 32)
	return int(year), err
}

func (c Chart) PrintJSON() {
	chartJSON, err := json.Marshal(c)
	if err != nil {
		fmt.Printf("error marshalling chart: %s\n", err)
		return
	}
	fmt.Println(string(chartJSON))
}
