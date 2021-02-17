package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	portFlagPtr := flag.Int("port", 11011, "The port to expose metrics to")
)

func main() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(string(*portFlagPtr), nil))
}

type Wordpress struct {
	posts, categories, tags, pages, comments, media, users, admin_users, themes int 
}

type WordpressCollector struct {
	posts *prometheus.Desc
	categories *prometheus.Desc
	tags *prometheus.Desc
	pages *prometheus.Desc
	comments *prometheus.Desc
	media *prometheus.Desc
	users *prometheus.Desc
	admin_users *prometheus.Desc
	themes *prometheus.Desc
}

func (c *WordpressCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.posts
	ch <- c.categories
	ch <- c.tags
	ch <- c.pages
	ch <- c.comments
	ch <- c.media
	ch <- c.users
	ch <- c.admin_users
	ch <- c.themes
}

func newWordpressCollector() *WordpressCollector {
	return &WordpressCollector{
		posts: prometheus.NewDesc("wordpress_post_count", "WordPress posts count", nil, nil),
		categories: prometheus.NewDesc("wordpress_category_count", "WordPress category count", nil, nil),
		tags: prometheus.NewDesc("wordpress_tag_count", "WordPress tags count", nil, nil),
		pages: prometheus.NewDesc("wordpress_page_count", "WordPress pages count", nil, nil),
		comments: prometheus.NewDesc("wordpress_comment_count", "WordPress comments count", nil, nil),
		media: prometheus.NewDesc("wordpress_media_count", "WordPress media files count", nil, nil),
		users: prometheus.NewDesc("wordpress_user_count", "WordPress users count", nil, nil),
		admin_users: prometheus.NewDesc("wordpress_admin_user_count", "WordPress administrator-level user count", nil, nil),
		themes: prometheus.NewDesc("wordpress_theme_count", "WordPress theme count", nil, nil),
	}
}

func errCheck(e error) {
	if e != nil {
		log.Println(e)
	}
}

