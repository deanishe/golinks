package main

import (
	"fmt"
	"log"
	"os"

	"github.com/namsral/flag"
)

var (
	bookmarks *Bookmarks
	cfg       Config
)

func main() {
	var (
		version bool
		config  string
		dbpath  string
		bmpath  string
		bind    string
	)

	flag.BoolVar(&version, "v", false, "display version information")

	flag.StringVar(&dbpath, "export", "", "export bookmarks from a legacy database")

	flag.StringVar(&config, "config", "", "config file")
	flag.StringVar(&bmpath, "bookmarks", "bookmarks.toml", "path to bookmarks file")
	flag.StringVar(&cfg.Title, "title", "Search", "OpenSearch title")
	flag.StringVar(&bind, "bind", "0.0.0.0:8000", "[int]:<port> to bind to")
	flag.StringVar(&cfg.FQDN, "fqdn", "localhost:8000", "FQDN for public access")
	flag.StringVar(&cfg.URL, "url", DefaultURL, "default URL to redirect to")
	flag.StringVar(&cfg.SuggestURL, "suggest", DefaultSuggestURL,
		"default URL to retrieve search suggestions from")

	flag.Parse()

	if version {
		fmt.Println(FullVersion())
		os.Exit(0)
	}

	if dbpath != "" {
		if err := exportDatabase(dbpath); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var err error
	bookmarks, err = NewBookmarks(bmpath)
	if err != nil {
		log.Fatal(err)
	}

	if bookmarks.Len() == 0 {
		if err := bookmarks.addDefaults(); err != nil {
			log.Fatal(err)
		}
	}

	NewServer(bind, cfg).ListenAndServe()
}
