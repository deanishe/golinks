package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prologic/bitcask"
	"github.com/stretchr/testify/assert"
)

// Foo ...
type Foo struct{}

// Name ...
func (f Foo) Name() string {
	return "foo"
}

// Desc ...
func (f Foo) Desc() string {
	return "foo bar"
}

// Exec ...
func (f Foo) Exec(w http.ResponseWriter, r *http.Request, args []string) error {
	w.Write([]byte(fmt.Sprintf("foo bar")))
	return nil
}

func TestNewCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Foo{}
	assert.Equal(cmd.Name(), "foo")
	assert.Equal(cmd.Desc(), "foo bar")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=foo", nil)

	args := []string{}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	body := w.Body.String()
	assert.Equal(body, "foo bar")
}

func TestPingCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Ping{}
	assert.Equal(cmd.Name(), "ping")
	assert.Contains(cmd.Desc(), "ping")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=ping", nil)

	args := []string{}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	body := w.Body.String()
	assert.Regexp("^pong [0-9]+$", body)
}

func TestListCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	err := EnsureDefaultBookmarks()
	assert.Nil(err)

	cmd := List{}
	assert.Equal(cmd.Name(), "list")
	assert.Contains(cmd.Desc(), "list")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=list", nil)

	args := []string{}
	err = cmd.Exec(w, r, args)
	assert.Nil(err)

	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})
	assert.Equal(w.Header().Get("Location"), "/list")
}

func TestHelpCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Help{}
	assert.Equal(cmd.Name(), "help")
	assert.Contains(cmd.Desc(), "help")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=help", nil)

	args := []string{}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	assert.Condition(func() bool {
		return w.Code >= http.StatusMultipleChoices &&
			w.Code <= http.StatusTemporaryRedirect
	})

	assert.Equal(w.Header().Get("Location"), "/help")
}

func TestTimeCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Time{}
	assert.Equal(cmd.Name(), "time")
	assert.Contains(cmd.Desc(), "time")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=time", nil)

	args := []string{}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	body := w.Body.String()
	ts, err := time.Parse("15:04:05", body)
	assert.Equal(ts.Hour(), time.Now().Hour())
	assert.Equal(ts.Minute(), time.Now().Minute())
}

func TestDateCommand(t *testing.T) {
	assert := assert.New(t)

	cmd := Date{}
	assert.Equal(cmd.Name(), "date")
	assert.Contains(cmd.Desc(), "date")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=date", nil)

	args := []string{}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	body := w.Body.String()
	ts, err := time.Parse(http.TimeFormat, body)
	assert.Equal(ts.Year(), time.Now().Year())
	assert.Equal(ts.Month(), time.Now().Month())
	assert.Equal(ts.Day(), time.Now().Day())
	assert.Equal(ts.Hour(), time.Now().Hour())
	assert.Equal(ts.Minute(), time.Now().Minute())
}

func TestAddCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	cmd := Add{}
	assert.Equal(cmd.Name(), "add")
	assert.Contains(cmd.Desc(), "add")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=add", nil)

	args := []string{"g", "https://www.google.com/search?q=%s&btnK"}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	bookmark, ok := LookupBookmark("g")
	assert.True(ok)

	assert.Equal(bookmark.Name(), "g")
	assert.Equal(bookmark.URL(), "https://www.google.com/search?q=%s&btnK")
}

func TestRemoveCommand(t *testing.T) {
	assert := assert.New(t)

	db, _ = bitcask.Open("test.db")
	defer db.Close()

	cmd := Remove{}
	assert.Equal(cmd.Name(), "remove")
	assert.Contains(cmd.Desc(), "remove")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "?q=remove", nil)

	args := []string{"imdb"}
	err := cmd.Exec(w, r, args)
	assert.Nil(err)

	bookmark, ok := LookupBookmark("imdb")
	assert.False(ok)

	assert.Equal("", bookmark.Name())
	assert.Equal("", bookmark.URL())
}
