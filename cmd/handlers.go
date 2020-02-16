package main
//go:generate statik -f -src=../web/public/ -dest=./internal/

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/destari/pingtrack/cmd/internal/statik"
)

func DataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if vars["hostname"] != "" {
		jsonResults, _ := json.Marshal(data.Results[vars["hostname"]])
		w.Write(jsonResults)
	} else {
		jsonResults, _ := json.Marshal(data)
		w.Write(jsonResults)
	}
}

func HostsHandler(w http.ResponseWriter, r *http.Request) {
	keys := []string{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	for k := range data.Results {
		keys = append(keys, k)
	}

	jsonResults, _ := json.Marshal(keys)
	w.Write(jsonResults)
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.FileServer(statikFS).ServeHTTP(w, r)
}
