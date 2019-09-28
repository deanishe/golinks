package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"

	// Logging
	"github.com/unrolled/logger"

	// Stats/Metrics
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/exp"
	"github.com/thoas/stats"

	rice "github.com/GeertJohan/go.rice"
	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
)

// Counters ...
type Counters struct {
	r metrics.Registry
}

func NewCounters() *Counters {
	counters := &Counters{
		r: metrics.NewRegistry(),
	}
	return counters
}

func (c *Counters) Inc(name string) {
	metrics.GetOrRegisterCounter(name, c.r).Inc(1)
}

func (c *Counters) Dec(name string) {
	metrics.GetOrRegisterCounter(name, c.r).Dec(1)
}

func (c *Counters) IncBy(name string, n int64) {
	metrics.GetOrRegisterCounter(name, c.r).Inc(n)
}

func (c *Counters) DecBy(name string, n int64) {
	metrics.GetOrRegisterCounter(name, c.r).Dec(n)
}

// Server ...
type Server struct {
	bind      string
	config    Config
	templates *Templates
	router    *httprouter.Router

	// Logger
	logger *logger.Logger

	// Stats/Metrics
	counters *Counters
	stats    *stats.Stats
}

func (s *Server) render(name string, w http.ResponseWriter, ctx interface{}) {
	buf, err := s.templates.Exec(name, ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if buf != nil {
		_, err = buf.WriteTo(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// IndexHandler ...
func (s *Server) IndexHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var (
			q    string
			cmd  string
			args []string
		)

		s.counters.Inc("n_index")

		// Query ?q=
		q = r.URL.Query().Get("q")

		// Form name=q
		if q == "" {
			q = r.FormValue("q")
		}

		if q != "" {
			tokens := strings.Split(q, " ")
			if len(tokens) > 0 {
				cmd, args = tokens[0], tokens[1:]
			}
		} else {
			cmd = p.ByName("command")
			args = strings.Split(p.ByName("args"), "/")
		}

		if cmd == "" {
			s.render("index", w, nil)
		} else {
			if command := LookupCommand(cmd); command != nil {
				err := command.Exec(w, r, args)
				if err != nil {
					http.Error(
						w,
						fmt.Sprintf(
							"Error processing command %s: %s",
							command.Name(), err,
						),
						http.StatusInternalServerError,
					)
				}
			} else if bookmark, ok := LookupBookmark(cmd); ok {
				q := strings.Join(args, " ")
				bookmark.Exec(w, r, q)
			} else {
				if s.config.URL != "" {
					url := s.config.URL
					if q != "" {
						url = fmt.Sprintf(url, q)
					}
					http.Redirect(w, r, url, http.StatusFound)
				} else {
					http.Error(
						w,
						fmt.Sprintf("Invalid Command: %v", cmd),
						http.StatusBadRequest,
					)
				}
			}
		}
	}
}

// HelpHandler ...
func (s *Server) HelpHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.counters.Inc("n_help")

		s.render("help", w, nil)
	}
}

// ListHandler ...
func (s *Server) ListHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.counters.Inc("n_list")

		var (
			bk  []Bookmark
			cmd []Command
		)

		prefix := []byte("bookmark_")
		err := db.Scan(prefix, func(key []byte) error {
			val, err := db.Get(key)
			if err != nil {
				return err
			}
			name := strings.TrimPrefix(string(key), "bookmark_")
			bk = append(bk, Bookmark{name, string(val)})
			return nil
		})
		if err != nil {
			log.Printf("error reading list of bookmarks: %s", err)
		}

		var names []string
		for k := range commands {
			names = append(names, k)
		}
		sort.Sort(sort.StringSlice(names))
		for _, name := range names {
			cmd = append(cmd, commands[name])
		}

		data := map[string]interface{}{
			"Bookmarks": bk,
			"Commands":  cmd,
		}
		s.render("list", w, data)
	}
}

// OpenSearchHandler ...
func (s *Server) OpenSearchHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s.counters.Inc("n_opensearch")

		w.Header().Set("Content-Type", "text/xml")
		w.Write(
			[]byte(fmt.Sprintf(
				OpenSearchTemplate,
				s.config.Title,
				s.config.FQDN,
			)),
		)
	}
}

// StatsHandler ...
func (s *Server) StatsHandler() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		bs, err := json.Marshal(s.stats.Data())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(bs)
	}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	log.Fatal(
		http.ListenAndServe(
			s.bind,
			s.logger.Handler(
				s.stats.Handler(
					gziphandler.GzipHandler(
						s.router,
					),
				),
			),
		),
	)
}

func (s *Server) initRoutes() {
	s.router.Handler("GET", "/debug/metrics", exp.ExpHandler(s.counters.r))
	s.router.GET("/debug/stats", s.StatsHandler())

	s.router.GET("/", s.IndexHandler())
	s.router.POST("/", s.IndexHandler())
	s.router.GET("/help", s.HelpHandler())
	s.router.GET("/list", s.ListHandler())
	s.router.GET("/opensearch.xml", s.OpenSearchHandler())
}

// NewServer ...
func NewServer(bind string, config Config) *Server {
	server := &Server{
		bind:      bind,
		config:    config,
		router:    httprouter.New(),
		templates: NewTemplates("base"),

		// Logger
		logger: logger.New(logger.Options{
			Prefix:               "golinks",
			RemoteAddressHeaders: []string{"X-Forwarded-For"},
			OutputFlags:          log.LstdFlags,
		}),

		// Stats/Metrics
		counters: NewCounters(),
		stats:    stats.New(),
	}

	// Templates
	box := rice.MustFindBox("templates")

	indexTemplate := template.New("index")
	template.Must(indexTemplate.Parse(box.MustString("index.html")))
	template.Must(indexTemplate.Parse(box.MustString("base.html")))

	helpTemplate := template.New("help")
	template.Must(helpTemplate.Parse(box.MustString("help.html")))
	template.Must(helpTemplate.Parse(box.MustString("base.html")))

	listTemplate := template.New("list")
	template.Must(listTemplate.Parse(box.MustString("list.html")))
	template.Must(listTemplate.Parse(box.MustString("base.html")))

	server.templates.Add("index", indexTemplate)
	server.templates.Add("help", helpTemplate)
	server.templates.Add("list", listTemplate)

	server.initRoutes()

	return server
}
