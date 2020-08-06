package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"sync"

	yaml "gopkg.in/yaml.v2"
)

// Bookmarks contains user bookmarks.
type Bookmarks struct {
	filename  string
	bookmarks map[string]string
	mu        *sync.RWMutex
}

// NewBookmarks loads bookmarks from a file (if it exists).
func NewBookmarks(filename string) (*Bookmarks, error) {
	bm := &Bookmarks{
		filename:  filename,
		bookmarks: map[string]string{},
		mu:        &sync.RWMutex{},
	}

	if err := bm.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return bm, nil
}

// Add saves a new bookmark.
func (bm *Bookmarks) Add(name, url string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.bookmarks[name] = url
	return bm.save()
}

// Get retrieves a bookmark.
func (bm *Bookmarks) Get(name string) (bookmark Bookmark, ok bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	var url string
	if url, ok = bm.bookmarks[name]; ok {
		bookmark = Bookmark{name, url}
	}

	return
}

// Delete removes a bookmark.
func (bm *Bookmarks) Delete(name string) error {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if _, ok := bm.bookmarks[name]; !ok {
		return errors.New("unknown bookmark")
	}

	delete(bm.bookmarks, name)
	return bm.save()
}

// Len returns the number of bookmarks.
func (bm *Bookmarks) Len() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return len(bm.bookmarks)
}

// All returns all bookmarks in alphabetical order.
func (bm *Bookmarks) All() []Bookmark {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	var (
		bk    = make([]Bookmark, len(bm.bookmarks))
		names []string
	)

	for k := range bm.bookmarks {
		names = append(names, k)
	}
	sort.Strings(names)

	for i, k := range names {
		bk[i] = Bookmark{k, bm.bookmarks[k]}
	}

	return bk
}

// load bookmarks from file
func (bm *Bookmarks) load() error {
	if bm.filename == "" {
		return nil
	}

	bm.mu.Lock()
	defer bm.mu.Unlock()

	data, err := ioutil.ReadFile(bm.filename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, &bm.bookmarks); err != nil {
		return err
	}
	return nil
}

// save bookmarks to file
func (bm *Bookmarks) save() error {
	if bm.filename == "" {
		return nil
	}

	data, err := yaml.Marshal(bm.bookmarks)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(bm.filename, data, 0600)
}

// add default bookmarks
func (bm *Bookmarks) addDefaults() error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for k, v := range DefaultBookmarks {
		bm.bookmarks[k] = v
	}
	return bm.save()
}

// Bookmark ...
type Bookmark struct {
	name string
	url  string
}

// Name ...
func (b Bookmark) Name() string {
	return b.name
}

// URL ...
func (b Bookmark) URL() string {
	return b.url
}

// Exec ...
func (b Bookmark) Exec(w http.ResponseWriter, r *http.Request, q string) {
	url := b.url
	if q != "" {
		url = fmt.Sprintf(b.url, q)
	}
	http.Redirect(w, r, url, http.StatusFound)
}
