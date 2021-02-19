package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
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
	prometheus.MustRegister(newWordpressCollector())
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Started WordPress exporter for %s\n", *monitorWordPress)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}

type Wordpress struct {
	posts, categories, tags, pages, comments, media, users, adminUsers, themes, plugins int
}

type WordpressCollector struct {
	posts      *prometheus.Desc
	categories *prometheus.Desc
	tags       *prometheus.Desc
	pages      *prometheus.Desc
	comments   *prometheus.Desc
	media      *prometheus.Desc
	users      *prometheus.Desc
	adminUsers *prometheus.Desc
	themes     *prometheus.Desc
	plugins    *prometheus.Desc
}

type ConfigFile struct {
	MonitoredWordpress []string `yaml:"wordpress-exporter"`
}

type Settings struct {
	title          string `json:"title"`
	language       string `json:"language"`
	ping_status    string `json:"ping_status"`
	comment_status string `json:"comment_status"`
}

func (c *WordpressCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.posts
	ch <- c.categories
	ch <- c.tags
	ch <- c.pages
	ch <- c.comments
	ch <- c.media
	ch <- c.users
	ch <- c.adminUsers
	ch <- c.themes
	ch <- c.plugins
}

func newWordpressCollector() *WordpressCollector {
	// var labels map[[]byte]interface{}
	// settingsClient := http.Client{}
	// request, err := http.NewRequest("GET", *monitorWordPress+"/wp-json/wp/v2/settings", nil)
	// errCheck(err)
	// request.Header.Set("User-Agent", USER_AGENT)
	// result, err := settingsClient.Do(request)
	// errCheck(err)
	// resultBody, err := ioutil.ReadAll(result.Body)
	// errCheck(err)
	// json.Unmarshal(resultBody, &labels)

	return &WordpressCollector{
		posts:      prometheus.NewDesc("wordpress_post_count", "WordPress posts count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		categories: prometheus.NewDesc("wordpress_category_count", "WordPress category count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		tags:       prometheus.NewDesc("wordpress_tag_count", "WordPress tags count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		pages:      prometheus.NewDesc("wordpress_page_count", "WordPress pages count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		comments:   prometheus.NewDesc("wordpress_comment_count", "WordPress comments count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		media:      prometheus.NewDesc("wordpress_media_count", "WordPress media files count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		users:      prometheus.NewDesc("wordpress_user_count", "WordPress users count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		adminUsers: prometheus.NewDesc("wordpress_admin_user_count", "WordPress administrator-level user count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		themes:     prometheus.NewDesc("wordpress_theme_count", "WordPress theme count", nil, prometheus.Labels{"instance": *monitorWordPress}),
		plugins:    prometheus.NewDesc("wordpress_plugin_count", "Wordpress plugin count", nil, prometheus.Labels{"instance": *monitorWordPress}),
	}
}

func (c *ConfigFile) ParseConf() *ConfigFile {
	yamlFile, err := ioutil.ReadFile(*configFile)
	errCheck(err)
	err = yaml.Unmarshal(yamlFile, c)
	errCheck(err)
	return c
}

func FetchJSONFromEndpoint(APIEndpoint string, auth bool) []byte {
	APIBase := *monitorWordPress
	HTTPClient := &http.Client{}
	fetchURL := fmt.Sprintf("%s%s", APIBase, APIEndpoint)
	request, err := http.NewRequest("GET", fetchURL, nil)
	request.Header.Set("User-Agent", UserAgent)
	errCheck(err)
	if auth {
		request.Header.Add("Authorization", "Basic "+BasicAuth(*authUsername, *authPassword))
	}
	response, err := HTTPClient.Do(request)
	errCheck(err)
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

// count items returned in JSON and return length
func CountJSONItems(JSONResponse []byte) int {
	var JSONObject interface{}
	json.Unmarshal(JSONResponse, &JSONObject)

	JSONObjectSlice, isOK := JSONObject.([]interface{})
	if !isOK {
		fmt.Println("Cannot convert the JSON object")
	}

	return len(JSONObjectSlice)
}

func BasicAuth(username, password string) string {
	authString := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(authString))
}

func (c *WordpressCollector) Collect(ch chan<- prometheus.Metric) {
	var _wordpress = new(Wordpress)
	_wordpress.categories = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories", *useAuth))
	_wordpress.posts = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/posts", *useAuth))
	_wordpress.tags = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/tags", *useAuth))
	_wordpress.pages = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/pages", *useAuth))
	_wordpress.comments = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/comments", *useAuth))
	_wordpress.media = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/media", *useAuth))
	_wordpress.users = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/users", *useAuth))
	_wordpress.themes = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories", *useAuth))
	_wordpress.plugins = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/plugins", *useAuth))

	ch <- prometheus.MustNewConstMetric(c.categories, prometheus.GaugeValue, float64(_wordpress.categories))
	ch <- prometheus.MustNewConstMetric(c.posts, prometheus.GaugeValue, float64(_wordpress.posts))
	ch <- prometheus.MustNewConstMetric(c.tags, prometheus.GaugeValue, float64(_wordpress.tags))
	ch <- prometheus.MustNewConstMetric(c.pages, prometheus.GaugeValue, float64(_wordpress.pages))
	ch <- prometheus.MustNewConstMetric(c.comments, prometheus.GaugeValue, float64(_wordpress.comments))
	ch <- prometheus.MustNewConstMetric(c.media, prometheus.GaugeValue, float64(_wordpress.media))
	ch <- prometheus.MustNewConstMetric(c.users, prometheus.GaugeValue, float64(_wordpress.users))
	ch <- prometheus.MustNewConstMetric(c.themes, prometheus.GaugeValue, float64(_wordpress.themes))
	ch <- prometheus.MustNewConstMetric(c.plugins, prometheus.GaugeValue, float64(_wordpress.plugins))
}

func errCheck(e error) {
	if e != nil {
		log.Println(e)
	}
}
