package main
//go:generate statik -f -src=../web/public/ -dest=./internal/

import (
	"fmt"
	"log"
	"net"
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

	/*
	pinger.OnRecv = func(pkt *ping.Packet) {
		res.IP = pkt.IPAddr.String()
	}
	*/

	pinger.Count = 4
	pinger.Interval = time.Millisecond*100
	pinger.Timeout = time.Second*1
	pinger.SetPrivileged(false)

	pinger.OnFinish = func(stats *ping.Statistics) {
		res.PacketLoss = stats.PacketLoss
		//res.IP = stats.IPAddr.String()
		res.MinRtt = stats.MinRtt
		res.AvgRtt = stats.AvgRtt
		res.MaxRtt = stats.MaxRtt
		config.PingCount += pinger.Count // stats
	}

	pinger.Run()
	return res
}

func pinger(queue chan string, output chan Results) {
	for {
		select {
		case host := <- queue:
			res := pingHost(host)
			output <- res
		}
	}
}

func fillQueue(queue chan string) {
	fmt.Println("Periodic queue filler started: Items: ", len(config.Hosts))
	tick := time.Tick(time.Duration(config.EchoTimes) * time.Second)

	for {
		select {
		case <-tick:
			for _, testHost := range config.Hosts {
				queue <- testHost
			}
		}
	}
}

func resultsReader(resq chan Results) {
	fmt.Println("Starting results reader")
	for {
		select {
		case result := <- resq:
			StoreResult(result)
		default:
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func Hosts(cidr string) ([]string, error) {
	ips := []string{}

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		ips = append(ips, cidr)
		return ips, nil
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

type Config struct {
	Hosts     []string
	EchoTimes int
	PingCount int
	Threads   int
}

var config Config

var data Data
var serveHost string
var servePort string
var databaseFile string
var ttl int
var rootCmd = &cobra.Command{Use: "pingtrack"}

func init() {
	var cmdPrint = &cobra.Command{
		Use:   "hosts [list of hosts / CIDRs]",
		Short: "Pings and tracks the list of hosts.",
		Long: `Pings and tracks the list of hosts. Use commas to separate IPs, CIDRs (192.168.1.0/24)`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hosts: " + strings.Join(args, ", "))
			//hosts = args
			var tmpHosts []string
			for _, h := range args {
				ips, _ := Hosts(h)
				tmpHosts = append(tmpHosts, ips...)
			}
			config.Hosts = tmpHosts
		},
	}

	rootCmd.PersistentFlags().IntVarP(&config.EchoTimes, "interval", "i", 10, "Time in seconds between pings")
	rootCmd.PersistentFlags().IntVarP(&ttl, "ttl", "e", 30, "Time in days to expire data")
	rootCmd.PersistentFlags().IntVarP(&config.Threads, "threads", "t", 50, "Number of threads to run in parallel")
	rootCmd.PersistentFlags().StringVarP(&serveHost, "bindhost", "H", "127.0.0.1", "Local host/IP to bind web server to")
	rootCmd.PersistentFlags().StringVarP(&servePort, "bindport", "p", "8080", "Port to bind web server to")
	rootCmd.PersistentFlags().StringVarP(&databaseFile, "database", "D", "pingtrack.db", "Database file name (:memory: for in-memory only")

	rootCmd.AddCommand(cmdPrint)
}

func main() {
	config.Hosts = []string{}

	rootCmd.Execute()

	if len(config.Hosts) == 0 {
		fmt.Println("NO HOSTS!")
		return
	}

	err := OpenStore(databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer CloseStore()

	data.Time = time.Now().Unix()
	data.Results = make(map[string][]Results)

	c := make(chan string)
	resq := make(chan Results)

	r := mux.NewRouter()
	r.HandleFunc("/api/config/", ConfigHandler)
	r.HandleFunc("/api/hosts/", HostsHandler)
	r.HandleFunc("/api/hosts/{hostname}", HostsHandler).Methods("DELETE")
	r.HandleFunc("/api/data/{hostname}", HostDataHandler)
	r.HandleFunc("/api/data/", DataHandler)

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
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go resultsReader(resq)
	go fillQueue(c)

	fmt.Println("Starting pinger threads..")
	for t := 0; t < config.Threads; t++ {
		go pinger(c, resq)
	}

	fmt.Println("Starting HTTP server: " + serveHost+":"+servePort)
	go log.Fatal(srv.ListenAndServe())
}