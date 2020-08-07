package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/fsnotify/fsnotify"
	yaml "gopkg.in/yaml.v2"
)

// Bookmarks contains user bookmarks.
type Bookmarks struct {
	filename  string
	bookmarks map[string]string
	mu        *sync.RWMutex // lock around access to bookmarks

	saving  bool // tell watcher to ignore write events
	watcher *fsnotify.Watcher

	reloaded chan struct{} // test hook
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

	bm.bookmarks = map[string]string{}
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

	bm.saving = true
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

func (bm *Bookmarks) startWatching() error {
	var err error
	if bm.watcher, err = fsnotify.NewWatcher(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case ev, ok := <-bm.watcher.Events:
				if !ok {
					return
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					if bm.saving {
						bm.saving = false
						continue
					}

					if err := bm.load(); err != nil {
						log.Printf("error reloading %q: %v", bm.filename, err)
					} else {
						log.Println("reloaded", bm.filename)
					}

					// signal unit test
					if bm.reloaded != nil {
						bm.reloaded <- struct{}{}
					}
				}

			case err, ok := <-bm.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("error watching %q: %v", bm.filename, err)
			}
		}
	}()

	return bm.watcher.Add(bm.filename)
}

func (bm *Bookmarks) stopWatching() error {
	if bm.watcher != nil {
		return bm.watcher.Close()
	}
	return nil
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
