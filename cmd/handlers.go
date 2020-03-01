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
	"strconv"
	"time"

	_ "github.com/destari/pingtrack/cmd/internal/statik"
)

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResults, _ := json.Marshal(config)
	w.Write(jsonResults)
}

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

func HostDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	duration := int64(14400)
	end := time.Now().Unix()
	start := end - duration

	if vars["hostname"] != "" {
		w.WriteHeader(http.StatusOK)

		newStart := r.URL.Query().Get("start")
		newDuration := r.URL.Query().Get("duration")

		if newStart != "" {
			s, _ := strconv.Atoi(newStart)
			start = int64(s)

			end = start + duration
		}

		if newDuration != "" {
			e, _ := strconv.Atoi(newDuration)
			duration = int64(e)

			end = start + duration
		}

		data := StoreRetrieve(vars["hostname"], start)
		jsonResults, _ := json.Marshal(data)
		w.Write(jsonResults)
	} else {
		w.WriteHeader(422)
		//jsonResults, _ := json.Marshal([""])
		//w.Write(jsonResults)
	}
}


func HostsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	removeHost := ""

	if vars["hostname"] != "" {
		removeHost = vars["hostname"]
	}

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		var payload = map[string]string{
			"hostname": "",
		}

		if err := json.Unmarshal(body, &payload); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}

		config.Hosts = append(config.Hosts, payload["hostname"])
	} else if r.Method == "DELETE" {

		tmpHosts := []string{}
		for _, h := range config.Hosts {
			if h != removeHost {
				tmpHosts = append(tmpHosts, h)
			}
		}
		config.Hosts = tmpHosts
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	/*
	for k := range data.Results {
		keys = append(keys, k)
	}
	*/

	jsonResults, _ := json.Marshal(config.Hosts)
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
