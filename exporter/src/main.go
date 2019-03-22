package main

import (
	"./growatt"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	totalEnergy = make(map[string]prometheus.Counter)
	currentPower = make(map[string]prometheus.Gauge)
)

func main() {
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	addr     := flag.String("a", ":5000", "The address to listen on for HTTP requests.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Example: %s -u john -p secret https://server.growatt.com/\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [http[s]://]hostname[:port]/[path/]\nOptions are:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 || *username == "" || *password == ""{
		flag.Usage()
		os.Exit(1)
	}
	u, err := url.ParseRequestURI(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	hostAddr, err := net.LookupHost(u.Host)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "Username: %s\nPassword: %s\n", *username, *password)
	fmt.Fprintf(os.Stdout, "Growatt API: %s\nGrowatt Server Address: %s\n", flag.Arg(0), hostAddr)
	fmt.Fprintf(os.Stdout, "\nRunning on %s ...\n", *addr)

	c := growatt.New(*username, *password, flag.Arg(0))
	initMetrics(c.PlantList())
	recordMetrics(*c)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*addr, nil)
}

func recordMetrics(c growatt.HttpClient) {
	go func() {
		for {
			for _, p := range c.PlantList() {
				tcv := readCounter(totalEnergy[p.PlantID])
				tes := strings.Split(p.TotalEnergy, " ")
				te, _ := strconv.ParseFloat(tes[0], 64)
				if tes[1] == "kWh" {
					te *= 1000
				}
				if tes[1] == "MWh" {
					te *= 1000000
				}
				if tcv < te {
					totalEnergy[p.PlantID].Add(te - tcv)
				}
				ces := strings.Split(p.CurrentPower, " ")
				ce, _ := strconv.ParseFloat(ces[0], 64)
				if ces[1] == "kW" {
					ce *= 1000
				}
				if ces[1] == "MW" {
					ce *= 1000000
				}
				currentPower[p.PlantID].Set(ce)

			}
			time.Sleep(1 * time.Minute)
		}
	}()
}

func initMetrics(pl []growatt.Plant){
	for _, p := range pl {
		totalEnergy[p.PlantID] = promauto.NewCounter(prometheus.CounterOpts{
			Name: "plant_total_energy",
			Help: "Total energy produced in the given plant in Wh",
			ConstLabels: map[string]string{
				"plant_id": p.PlantID,
				"plant_name": p.PlantName,
				"is_have_storage": p.IsHaveStorage,
			},
		})
		currentPower[p.PlantID] = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "plant_current_power",
			Help: "Current energy produced in the given plant in W",
			ConstLabels: map[string]string{
				"plant_id": p.PlantID,
				"plant_name": p.PlantName,
				"is_have_storage": p.IsHaveStorage,
			},
		})
	}
}

func readCounter(m prometheus.Counter) float64 {
	pb := &dto.Metric{}
	m.Write(pb)
	return pb.GetCounter().GetValue()
}
