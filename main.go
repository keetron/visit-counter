package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
)

var httpAddr = flag.String("http", ":8080", "Bind address for the HTTP server")

// todo replace these containers with mutex for concurrency reasons
var (
	idUrl          []string
	distinctVisits = make(map[string]int)
)

type Visit struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type Counter struct {
	Url    string `json:"url"`
	Visits int    `json:"visits"`
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mux := serveMux()
	return http.ListenAndServe(*httpAddr, mux)
}

func serveMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/visit", visitHandler())
	return mux
}

func visitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			defer func() {
				_ = r.Body.Close()
			}()

			reqBody, _ := io.ReadAll(r.Body)
			var visit Visit
			err := json.Unmarshal(reqBody, &visit)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				// todo add some logging here
				return
			}

			d := visit.ID + visit.URL

			if newVisitor(d) {
				// todo fix race condition
				idUrl = append(idUrl, d)
				distinctVisits[visit.URL] = distinctVisits[visit.URL] + 1
			}

			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Query() != nil && r.Method == "GET" {
			q := r.URL.Query().Get("u")
			v := distinctVisits[q]
			c := Counter{
				Url:    q,
				Visits: v,
			}

			err := json.NewEncoder(w).Encode(c)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				// todo add some logging here
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func newVisitor(str string) bool {
	// todo this solution is not good, replace with mutex
	for _, v := range idUrl {
		if v == str {
			return false
		}
	}
	return true
}
