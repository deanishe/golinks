# golinks

![](https://github.com/prologic/golinks/workflows/Coverage/badge.svg)
![](https://github.com/prologic/golinks/workflows/Docker/badge.svg)
![](https://github.com/prologic/golinks/workflows/Go/badge.svg)
![](https://github.com/prologic/golinks/workflows/ReviewDog/badge.svg)

[![CodeCov](https://codecov.io/gh/prologic/golinks/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/golinks)
[![Go Report Card](https://goreportcard.com/badge/prologic/golinks)](https://goreportcard.com/report/prologic/golinks)
[![codebeat badge](https://codebeat.co/badges/15fba8a5-3044-4f40-936f-9e0f5d5d1fd9)](https://codebeat.co/projects/github-com-prologic-bitcask-master)
[![GoDoc](https://godoc.org/github.com/prologic/golinks?status.svg)](https://godoc.org/github.com/prologic/golinks)
[![GitHub license](https://img.shields.io/github/license/prologic/golinks.svg)](https://github.com/prologic/golinks)

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

### Docker

To startup a container with an image from docker hub, run the following:

```
docker run -it -d -p 8000:8000 prologic/golinks
```

Or use the following docker-compose configuration:

```
golinks:
    image: 'prologic/golinks:latest'
    container_name: golinks
    ports:
        - "8000:8000"
```

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

## Configuration

golinks comes with sensible defaults, so it will run out-of-the box without any configuration (just run `golinks` and it will be available at `http://localhost:8000`, and save your custom bookmarks to `search.db` in the working directory), but there are several knobs you can tweak.

All configuration options can be specified via command-line flag, environment variable or a configuration file.

The available CLI flags are (as shown via `golinks -h`):

|    Flag    |                                 Default                                 |                                      Description                                      |
|------------|-------------------------------------------------------------------------|---------------------------------------------------------------------------------------|
| `-bind`    | `0:0:0:0:8000`                                                          | IP and port to bind server to.                                                        |
| `-fqdn`    | `localhost:8000`                                                        | Web address that corresponds to bind address.                                            |
| `-dbpath`  | `search.db`                                                             | Database to save your custom bookmarks to.                                            |
| `-suggest` | `https://suggestqueries.google.com/complete/search?client=firefox&q=%s` | URL of autosuggest service to retrieve search suggestions from.                       |
| `-title`   | `Search`                                                                | The OpenSearch service title (i.e. what your browser will call golinks' search).      |
| `-url`     | `https://www.google.com/search?q=%s&btnK`                               | The URL golinks will redirect searches to by default (if no custom bookmark matches). |
| `-config`  |                                                                         | Path to the optional configuration file (see below).                                   |
| `-h`       |                                                                         | Show CLI help and exit.                                                                        |
| `-v`       |                                                                         | Show golinks version number and exit.                                                 |

### Environment variables

All the above flags can also be specified via environment variable with the same name as the flag, but in uppercase. So `BIND=127.0.0.1:8081 FQDN=localhost:8081 golinks` is equivalent to `golinks -bind 127.0.0.1:8081 -fqdn localhost:8081`.

### Configuration file

golinks can also use a configuration file (specified via `-config /path/to/file` or `CONFIG=/path/to/file`) with a very simple format. Each line of the file should contain one option, specified in the same way as the corresponding CLI flag but without the leading `-`. Empty lines and lines beginning with `#` are ignored:

```
# with a space
bind 127.0.0.1:8081

# with equals
fqdn=localhost:8081
```

### Example

So, assuming your name is "Dave", a Linux user and fan of DuckDuckGo, you might run golinks like this:

```
golinks -config ~/.config/golinks/config.cfg
```

where `/home/dave/.config/golinks/config.cfg` looks like this:

```
bind 127.0.0.1:3456
fqdn localhost:3456
dbpath /home/dave/.config/golinks/search.db
title Dave's Search
url https://duckduckgo.com/?q=%s
suggest https://duckduckgo.com/ac/?type=list&q=%s
```

## Stargazers over time

[![Stargazers over time](https://starcharts.herokuapp.com/prologic/golinks.svg)](https://starcharts.herokuapp.com/prologic/golinks)

## Support

Support the ongoing development of Bitcask!

**Sponsor**

- Become a [Sponsor](https://www.patreon.com/prologic)

## Contributing

golinks is considered "production" software and is used daily. If you find this interresting or useful please fork and contribute back via pull-requests! If you find bugs or have ideas for new features, please file an issue!

## License

MIT
