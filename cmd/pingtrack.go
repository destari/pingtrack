package main
//go:generate statik -f -src=../web/public/ -dest=./internal/

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

var data Data
var hosts []string
var echoTimes int
var serveHost string
var servePort string
var rootCmd = &cobra.Command{Use: "pingtrack"}

func init() {
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

	rootCmd.PersistentFlags().IntVarP(&echoTimes, "interval", "i", 10, "Time in seconds between pings")
	rootCmd.PersistentFlags().StringVarP(&serveHost, "bindhost", "H", "127.0.0.1", "Local host/IP to bind web server to")
	rootCmd.PersistentFlags().StringVarP(&servePort, "bindport", "p", "8080", "Port to bind web server to")

	rootCmd.AddCommand(cmdPrint)
}

func main() {
	hosts = []string{}




	//cmdPrint.Flags().IntVarP(&echoTimes, "interval", "i", 1, "Time in seconds between pings")



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
	r.HandleFunc("/api/hosts/", HostsHandler)
	r.HandleFunc("/api/hosts/{hostname}", HostsHandler).Methods("DELETE")
	//r.Methods("GET", "POST", "DELETE", "HEAD", "OPTIONS", "PUT")
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
		Addr:         serveHost+":"+servePort,
		// Good practice: enforce timeouts for servers you create!
		//WriteTimeout: 15 * time.Second,
		//ReadTimeout:  15 * time.Second,
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

	fmt.Println("Starting HTTP server: " + serveHost+":"+servePort)
	go log.Fatal(srv.ListenAndServe())
}