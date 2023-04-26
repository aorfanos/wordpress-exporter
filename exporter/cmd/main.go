package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aorfanos/wordpress-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	UserAgent = "prometheus-wordpress-exporter"
)

var (
	portNum          = flag.Int("port", 11012, "The port to expose metrics to")
	configFile       = flag.String("config.file", "wordpress-exporter.yml", "Configure which WordPress sites to monitor")
	monitorWordPress = flag.String("host", "", "Which host to monitor, with format <schema>://<host or FQDN>")
	useAuth          = flag.Bool("auth.basic", true, "Whether to use basic authentication (true|false)")
	authUsername     = flag.String("auth.user", "admin", "User to use with basic auth")
	authPassword     = flag.String("auth.pass", "admin", "Password to use with basic auth")
)

func init() {
	flag.Parse()
	wp := utils.NewWordpress(*monitorWordPress, UserAgent, *authUsername, *authPassword, *useAuth)
	prometheus.MustRegister(utils.NewWordpressCollector(wp))
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Started WordPress exporter for %s\n", *monitorWordPress)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}
