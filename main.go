package main

import (
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

var (
	portNum    = flag.Int("port", 11011, "The port to expose metrics to")
	configFile = flag.String("config.file", "wordpress-exporter.yml", "Configure which WordPress sites to monitor")
)

func init() {
	flag.Parse()
	prometheus.MustRegister(newWordpressCollector())
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portNum), nil))
}

type Wordpress struct {
	posts, categories, tags, pages, comments, media, users, adminUsers, themes int
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
}

type ConfigFile struct {
	MonitoredWordpress []string `yaml:"wordpress-exporter"`
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
}

func newWordpressCollector() *WordpressCollector {
	return &WordpressCollector{
		posts:      prometheus.NewDesc("wordpress_post_count", "WordPress posts count", nil, nil),
		categories: prometheus.NewDesc("wordpress_category_count", "WordPress category count", nil, nil),
		tags:       prometheus.NewDesc("wordpress_tag_count", "WordPress tags count", nil, nil),
		pages:      prometheus.NewDesc("wordpress_page_count", "WordPress pages count", nil, nil),
		comments:   prometheus.NewDesc("wordpress_comment_count", "WordPress comments count", nil, nil),
		media:      prometheus.NewDesc("wordpress_media_count", "WordPress media files count", nil, nil),
		users:      prometheus.NewDesc("wordpress_user_count", "WordPress users count", nil, nil),
		adminUsers: prometheus.NewDesc("wordpress_admin_user_count", "WordPress administrator-level user count", nil, nil),
		themes:     prometheus.NewDesc("wordpress_theme_count", "WordPress theme count", nil, nil),
	}
}

func (c *ConfigFile) ParseConf() *ConfigFile {
	yamlFile, err := ioutil.ReadFile(*configFile)
	errCheck(err)
	err = yaml.Unmarshal(yamlFile, c)
	errCheck(err)
	return c
}

func FetchJSONFromEndpoint(APIEndpoint string) []byte {
	APIBase := "https://aorfanos.com"
	fetchURL := fmt.Sprintf("%s%s", APIBase, APIEndpoint)
	response, err := http.Get(fetchURL)
	errCheck(err)
	data, _ := ioutil.ReadAll(response.Body)
	return data
}

func CountJSONItems(JSONResponse []byte) int {
	var JSONObject interface{}
	json.Unmarshal(JSONResponse, &JSONObject)

	JSONObjectSlice, isOK := JSONObject.([]interface{})
	if !isOK {
		fmt.Println("Cannot convert the JSON object")
	}

	return len(JSONObjectSlice)
}

func (c *WordpressCollector) Collect(ch chan<- prometheus.Metric) {
	var _wordpress = new(Wordpress)
	_wordpress.categories = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories"))
	_wordpress.posts = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/posts"))
	_wordpress.tags = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/tags"))
	_wordpress.pages = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/pages"))
	_wordpress.comments = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/comments"))
	_wordpress.media = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/media"))
	_wordpress.users = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/users"))
	_wordpress.themes = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories"))

	ch <- prometheus.MustNewConstMetric(c.categories, prometheus.GaugeValue, float64(_wordpress.categories))
	ch <- prometheus.MustNewConstMetric(c.posts, prometheus.GaugeValue, float64(_wordpress.posts))
	ch <- prometheus.MustNewConstMetric(c.tags, prometheus.GaugeValue, float64(_wordpress.tags))
	ch <- prometheus.MustNewConstMetric(c.pages, prometheus.GaugeValue, float64(_wordpress.pages))
	ch <- prometheus.MustNewConstMetric(c.comments, prometheus.GaugeValue, float64(_wordpress.comments))
	ch <- prometheus.MustNewConstMetric(c.media, prometheus.GaugeValue, float64(_wordpress.media))
	ch <- prometheus.MustNewConstMetric(c.users, prometheus.GaugeValue, float64(_wordpress.users))
	ch <- prometheus.MustNewConstMetric(c.themes, prometheus.GaugeValue, float64(_wordpress.themes))

}

func errCheck(e error) {
	if e != nil {
		log.Println(e)
	}
}
