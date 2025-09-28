package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"os"

	"github.com/aorfanos/wordpress-exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	UserAgent = "prometheus-wordpress-exporter"
)

var (
	portNum          = flag.Int("port", 11011, "The port to expose metrics to")
	configFile       = flag.String("config.file", "wordpress-exporter.yml", "Configure which WordPress sites to monitor")
	monitorWordPress = flag.String("host", "", "Which host to monitor, with format <schema>://<host or FQDN>")
	useAuth          = flag.Bool("auth.basic", true, "Whether to use basic authentication (true|false)")
	authUsername     = flag.String("auth.user", "admin", "User to use with basic auth")
	authPassword     = flag.String("auth.pass", "admin", "Password to use with basic auth")
)

func init() {
	flag.Parse()

	// Environment variable fallbacks if CLI flags match defaults
	defaultHost := ""
	defaultUser := "admin"
	defaultPass := "admin"
	defaultUseAuth := true
	defaultPort := 11011

	if *monitorWordPress == defaultHost {
		if env := os.Getenv("WORDPRESS_EXPORTER_HOST"); env != "" {
			*monitorWordPress = env
		}
	}

	if *authUsername == defaultUser {
		if env := os.Getenv("WORDPRESS_EXPORTER_AUTH_USER"); env != "" {
			*authUsername = env
		}
	}

	if *authPassword == defaultPass {
		if env := os.Getenv("WORDPRESS_EXPORTER_AUTH_PASS"); env != "" {
			*authPassword = env
		}
	}

	if *useAuth == defaultUseAuth {
		if envStr := os.Getenv("WORDPRESS_EXPORTER_AUTH_BASIC"); envStr != "" {
			if parsed, err := strconv.ParseBool(envStr); err == nil {
				*useAuth = parsed
			} else {
				log.Printf("Invalid WORDPRESS_EXPORTER_AUTH_BASIC env var: %v, using default %t", err, defaultUseAuth)
			}
		}
	}

	if *portNum == defaultPort {
		if envStr := os.Getenv("WORDPRESS_EXPORTER_PORT"); envStr != "" {
			if parsed, err := strconv.Atoi(envStr); err == nil {
				*portNum = parsed
			} else {
				log.Printf("Invalid WORDPRESS_EXPORTER_PORT env var: %v, using default %d", err, defaultPort)
			}
		}
	}

	if *configFile == "wordpress-exporter.yml" {
		if env := os.Getenv("WORDPRESS_EXPORTER_CONFIG_FILE"); env != "" {
			*configFile = env
		}
	}

	// Try to load config file for multi-host support
	hosts, err := utils.LoadConfig(*configFile)
	if err != nil {
		// Config file doesn't exist or has errors, fall back to single host mode
		log.Printf("Config file not found or invalid (%v), using single host mode", err)

		// Validate required host for single host mode
		if *monitorWordPress == "" {
			log.Fatal("WordPress host must be provided via --host CLI flag or WORDPRESS_EXPORTER_HOST env var")
		}

		// Single host mode
		wp := utils.NewWordpress(*monitorWordPress, UserAgent, *authUsername, *authPassword, *useAuth)
		prometheus.MustRegister(utils.NewWordpressCollector(wp))
		fmt.Printf("Configured for single host: %s\n", *monitorWordPress)
	} else {
		// Multi-host mode using config file
		if len(hosts) == 0 {
			log.Fatal("Config file exists but contains no hosts")
		}

		fmt.Printf("Configured for multi-host mode with %d hosts from config file\n", len(hosts))
		for i, host := range hosts {
			wp := utils.NewWordpress(host, UserAgent, *authUsername, *authPassword, *useAuth)
			prometheus.MustRegister(utils.NewWordpressCollector(wp))
			fmt.Printf("  Host %d: %s\n", i+1, host)
		}
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Started WordPress exporter on port %d\n", *portNum)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}
