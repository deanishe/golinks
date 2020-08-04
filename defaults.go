package main

const (
	// DefaultURL redirects to Google Search for unknown queries
	DefaultURL string = "https://www.google.com/search?q=%s&btnK"
	// DefaultSuggestURL provides search suggestions from Google
	DefaultSuggestURL string = "https://suggestqueries.google.com/complete/search?client=firefox&q=%s"
)

// DefaultBookmarks ...
var DefaultBookmarks = map[string]string{
	"g":    "https://www.google.com/search?q=%s&btnK",
	"gl":   "https://www.google.com/search?q=%s&btnI",
	"gh":   "https://github.com/search?q=%s&ref=opensearch",
	"go":   "https://golang.org/search?q=%s",
	"wp":   "http://en.wikipedia.org/?search=%s",
	"py":   "https://docs.python.org/2/search.html?q=%s",
	"py3":  "https://docs.python.org/3/search.html?q=%s",
	"yt":   "http://www.youtube.com/results?search_type=search_videos&search_sort=relevance&search_query=%s&search=Search",
	"gim":  "https://www.google.com/search?q=%s&um=1&ie=UTF-8&hl=en&tbm=isch",
	"gdef": "http://www.google.com/search?q=define%%3A+%s&hl=en&lr=&oi=definel&defl=all",
	"imdb": "http://www.imdb.com/find?q=%s",
	"gm":   "http://maps.google.com/maps?q=%s",
}
