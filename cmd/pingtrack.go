package main

//go:generate statik -f -src=../web/public/ -dest=./internal/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/sparrc/go-ping"
	"github.com/spf13/cobra"

	_ "github.com/destari/pingtrack/cmd/internal/statik"
)

type Results struct {
	Host string
	IP string
	MinRtt time.Duration
	AvgRtt time.Duration
	MaxRtt time.Duration
	PacketLoss float64
	Time time.Time
	EpochTime int64
}

type Data struct {
	Time int64
	Results map[string][]Results
}

func pingHost(host string) Results {

	res := Results{Host: host, Time: time.Now(), EpochTime: time.Now().Unix()}

	pinger, err := ping.NewPinger(host)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return res
	}

	// listen for ctrl-C signal
	catch := make(chan os.Signal, 1)
	signal.Notify(catch, os.Interrupt)
	go func() {
		for range catch {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *ping.Packet) {
		res.IP = pkt.IPAddr.String()
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		res.PacketLoss = stats.PacketLoss
		res.IP = stats.IPAddr.String()
		res.MinRtt = stats.MinRtt
		res.AvgRtt = stats.AvgRtt
		res.MaxRtt = stats.MaxRtt
	}

	pinger.Count = 4
	pinger.Interval = time.Millisecond*100
	pinger.Timeout = time.Second*2
	pinger.SetPrivileged(false)

	pinger.Run()
	return res
}

func pinger(queue chan string, output chan Results) {
	fmt.Println("Starting pinger..")
	for {
		select {
		case host := <- queue:
			//fmt.Println("Pinging")
			res := pingHost(host)
			output <- res
		}
	}
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

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
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

	// check whether a file exists at the given path
	/*
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	 */

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	// otherwise, use http.FileServer to serve the static dir
	//http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
	http.FileServer(statikFS).ServeHTTP(w, r)
}

var data Data

func main() {
	echoTimes := 10
	hosts := []string{}

	var cmdPrint = &cobra.Command{
		Use:   "hosts [hosts to ping]",
		Short: "Pings and tracks the list of hosts.",
		Long: `Pings and tracks the list of hosts.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hosts: " + strings.Join(args, ", "))
			hosts = args
		},
	}

	cmdPrint.Flags().IntVarP(&echoTimes, "interval", "i", 1, "Time in seconds between pings")

	var rootCmd = &cobra.Command{Use: "pingtrack"}
	rootCmd.AddCommand(cmdPrint)
	rootCmd.Execute()


	tick := time.Tick(time.Duration(echoTimes) * time.Second)

	if len(hosts) == 0 {
		fmt.Println("NO HOSTS!")
		return
	}


	data.Time = time.Now().Unix()
	data.Results = make(map[string][]Results)

	c := make(chan string)
	resq := make(chan Results)

	r := mux.NewRouter()
	r.HandleFunc("/api/data/", DataHandler)
	r.HandleFunc("/api/data/{hostname}/", DataHandler)

	spa := spaHandler{staticPath: "../web/public", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	headersOk := handlers.AllowedHeaders([]string{"*"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	handler := handlers.CORS(originsOk, headersOk, methodsOk)(r)

	http.Handle("/", handler)

	srv := &http.Server{
		Handler:      handler,
		Addr:         "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		fmt.Println("Starting results reader")
		for {
			select {
			case result := <- resq:
				data.Results[result.Host] = append(data.Results[result.Host], result)
			default:
				time.Sleep(1 * time.Second)
			}
		}
	}()

	/*
	fmt.Println("Filling initial queue.")
	for _, testHost := range hosts {
		c <- testHost
	}
	*/

	go func() {
		fmt.Println("Periodic queue")
		for {
			select {
			case <-tick:
				for _, testHost := range hosts {
					c <- testHost
				}
			}
		}
	}()

	go pinger(c, resq)

	fmt.Println("Starting HTTP server..")
	go log.Fatal(srv.ListenAndServe())
}