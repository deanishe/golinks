# golinks

[![Build Status](https://cloud.drone.io/api/badges/prologic/golinks/status.svg)](https://cloud.drone.io/prologic/golinks)
[![CodeCov](https://codecov.io/gh/prologic/golinks/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/golinks)
[![Go Report Card](https://goreportcard.com/badge/prologic/golinks)](https://goreportcard.com/report/prologic/golinks)
[![GoDoc](https://godoc.org/github.com/prologic/golinks?status.svg)](https://godoc.org/github.com/prologic/golinks) 

golinks is a web app that allows you to create smart bookmarks, commands and aliases by pointing your web browser's default search engine at a running instance. Similar to bunny1 or yubnub.

## Installation

### Source

```bash
$ go get github.com/prologic/golinks
```

### OS X Homebrew

There is a formula provided that you can tap and install from
[prologic/homebrew-golinks](https://github.com/prologic/homebrew-golinks):

```bash
$ brew tap prologic/golinks
$ brew install golinks
```

**NB:** This installs the latest released binary; so if you want a more
recent unreleased version from master you'll have to clone the repository
and build yourself.

## Usage

Run golinks:

```bash
$ golinks -bind 127.0.0.1:8000 -fqdn localhost:8000
```

Set your browser's default golinks engine to http://localhost:8000/?q=%s

- [Instructions for Chrome](https://support.google.com/chrome/answer/95426)
- [Instructions for Firefox](https://support.mozilla.org/en-US/kb/add-or-remove-search-engine-firefox#w_add-a-search-engine-from-the-address-bar)

Then type `help` to view the main help page, `g foo bar` to perform a [Google](https://google.com) search for "foo bar" or `list` to list all available commands.


### Custom bookmarks

To add a bookmark (or overwrite an existing one), enter `add [name] [url]` as your search query, where `name` is the shortcut for the bookmark and `url` the URL:

```
add imdb https://www.imdb.com
```

Now you can just enter `imdb` in your search bar to go straight to imdb.com.

You can also add `%s` to your URL, which will be replaced with your search query:

```
add ddg https://duckduckgo.com/?q=%s
```

Now you can use `ddg [query]` to search via DuckDuckGo, e.g. `ddg free stuff` to find yourself some free stuff.

To remove a search, use `remove [name]`, so `remove ddg` will remove the above search.

### Other commands

Use `list` to see all your bookmarks and commands (golinks comes with several useful built-ins) and `help` to view the online help page.

## Contributing

golinks is considered "production" software and is used daily. If you find this interresting or useful please fork and contribute back via pull-requests! If you find bugs or have ideas for new features, please file an issue!

## License

MIT
