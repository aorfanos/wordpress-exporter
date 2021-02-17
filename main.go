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

var (
	portNum      = flag.Int("port", 11011, "The port to expose metrics to")
	configFile   = flag.String("config.file", "wordpress-exporter.yml", "Configure which WordPress sites to monitor")
	authUsername = flag.String("auth.user", "admin", "User to use with basic auth")
	authPassword = flag.String("auth.pass", "admin", "Password to use with basic auth")
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

func FetchJSONFromEndpoint(APIEndpoint string, auth bool) []byte {
	APIBase := "https://aorfanos.com"
	HTTPClient := &http.Client{}
	fetchURL := fmt.Sprintf("%s%s", APIBase, APIEndpoint)
	request, err := http.NewRequest("GET", fetchURL, nil)
	errCheck(err)
	if auth {
		request.Header.Add("Authentication", BasicAuth(*authUsername, *authPassword))
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
	authHeaderValue := fmt.Sprintf("Basic ", base64.StdEncoding.EncodeToString([]byte(authString)))
	return authHeaderValue
}

func (c *WordpressCollector) Collect(ch chan<- prometheus.Metric) {
	var _wordpress = new(Wordpress)
	_wordpress.categories = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories", false))
	_wordpress.posts = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/posts", false))
	_wordpress.tags = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/tags", false))
	_wordpress.pages = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/pages", false))
	_wordpress.comments = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/comments", false))
	_wordpress.media = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/media", false))
	_wordpress.users = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/users", false))
	_wordpress.themes = CountJSONItems(FetchJSONFromEndpoint("/wp-json/wp/v2/categories", false))

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
