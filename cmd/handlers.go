package main


import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"io"
	"io/ioutil"
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
	//keys := []string{}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		var newHost string
		if err := json.Unmarshal(body, &newHost); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		hosts = append(hosts, newHost)
	} else if r.Method == "DELETE" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		var removeHost string
		if err := json.Unmarshal(body, &removeHost); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		tmpHosts := []string{}
		for _, h := range hosts {
			if h != removeHost {
				tmpHosts = append(tmpHosts, h)
			}
		}
		hosts = tmpHosts
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	/*
	for k := range data.Results {
		keys = append(keys, k)
	}
	*/

	jsonResults, _ := json.Marshal(hosts)
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
