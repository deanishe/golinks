package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func TestZeroBookmark(t *testing.T) {
	assert := assert.New(t)

	bookmark := Bookmark{}
	assert.Equal(bookmark.Name(), "")
	assert.Equal(bookmark.URL(), "")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	bookmark.Exec(w, r, "")
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})
}

func TestBookmarkWithQuery(t *testing.T) {
	assert := assert.New(t)

	bookmark := Bookmark{
		name: "g",
		url:  "https://www.google.com/search?q=%s&btnK",
	}
	assert.Equal(bookmark.Name(), "g")
	assert.Equal(bookmark.URL(), "https://www.google.com/search?q=%s&btnK")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	q := "foo bar"
	bookmark.Exec(w, r, q)
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		fmt.Sprintf(
			"https://www.google.com/search?q=%s&btnK",
			q,
		),
	)
}

func TestBookmarkWithoutQuery(t *testing.T) {
	assert := assert.New(t)

	bookmark := Bookmark{
		name: "g",
		url:  "https://www.google.com/",
	}
	assert.Equal(bookmark.Name(), "g")
	assert.Equal(bookmark.URL(), "https://www.google.com/")

	r, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	bookmark.Exec(w, r, "")
	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(
		w.Header().Get("Location"),
		"https://www.google.com/",
	)
}

func TestLoadBookmarks(t *testing.T) {
	assert := assert.New(t)

	test := map[string]string{
		"one":   "https://one.com",
		"two":   "https://two.com",
		"three": "https://three.com",
		"four":  "https://four.com",
	}

	data, err := toml.Marshal(test)
	assert.Nil(err)
	assert.Nil(ioutil.WriteFile("test.toml", data, 0666))

	bookmarks, err = NewBookmarks("test.toml")
	assert.Nil(err)

	for k, v := range test {
		bm, ok := bookmarks.Get(k)
		assert.True(ok)
		assert.Equal(k, bm.Name())
		assert.Equal(v, bm.URL())
	}
}

func TestSaveBookmarks(t *testing.T) {
	assert := assert.New(t)

	test := map[string]string{
		"one":   "https://one.com",
		"two":   "https://two.com",
		"three": "https://three.com",
		"four":  "https://four.com",
	}

	// clear file
	assert.Nil(ioutil.WriteFile("test.toml", []byte(``), 0666))

	var err error
	bookmarks, err = NewBookmarks("test.toml")
	assert.Nil(err)

	for k, v := range test {
		assert.Nil(bookmarks.Add(k, v))
	}

	var (
		data []byte
		bk   map[string]string
	)
	data, err = ioutil.ReadFile("test.toml")
	assert.Nil(err)

	err = toml.Unmarshal(data, &bk)
	assert.Nil(err)
	assert.Equal(test, bk)
}

func TestAddDefaultBookmarks(t *testing.T) {
	assert := assert.New(t)

	var err error
	bookmarks, err = NewBookmarks("")
	assert.Nil(err)
	assert.Nil(bookmarks.addDefaults())

	for k, v := range DefaultBookmarks {
		bookmark, ok := bookmarks.Get(k)
		assert.True(ok)
		assert.Equal(k, bookmark.Name())
		assert.Equal(v, bookmark.URL())
	}
}
