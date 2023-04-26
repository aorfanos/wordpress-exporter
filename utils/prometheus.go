package utils

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type WordpressCollector struct {
	posts      *prometheus.Desc
	categories *prometheus.Desc
	tags       *prometheus.Desc
	pages      *prometheus.Desc
	comments   *prometheus.Desc
	media      *prometheus.Desc
	users      *prometheus.Desc
	taxonomies *prometheus.Desc
	themes     *prometheus.Desc
	plugins    *prometheus.Desc
	Wp		 *Wordpress
}


func (c *WordpressCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.posts
	ch <- c.categories
	ch <- c.tags
	ch <- c.pages
	ch <- c.comments
	ch <- c.media
	ch <- c.users
	ch <- c.taxonomies
	ch <- c.themes
	ch <- c.plugins
}


func (c *WordpressCollector) Collect(ch chan<- prometheus.Metric) {
	var err error

	c.Wp.categories, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/categories"))
	ErrCheck(err)
	c.Wp.posts, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/posts"))
	ErrCheck(err)
	c.Wp.tags, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/tags"))
	ErrCheck(err)
	c.Wp.pages, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/pages"))
	ErrCheck(err)
	c.Wp.comments, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/comments"))
	ErrCheck(err)
	c.Wp.media, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/media"))
	ErrCheck(err)
	c.Wp.users, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/users"))
	ErrCheck(err)
	c.Wp.taxonomies, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/taxonomies"))
	ErrCheck(err)
	c.Wp.themes, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/themes"))
	ErrCheck(err)
	c.Wp.plugins, err = CountJSONItems(c.FetchJSONFromEndpoint("/wp-json/wp/v2/plugins"))
	ErrCheck(err)

	ch <- prometheus.MustNewConstMetric(c.categories, prometheus.GaugeValue, float64(c.Wp.categories))
	ch <- prometheus.MustNewConstMetric(c.posts, prometheus.GaugeValue, float64(c.Wp.posts))
	ch <- prometheus.MustNewConstMetric(c.tags, prometheus.GaugeValue, float64(c.Wp.tags))
	ch <- prometheus.MustNewConstMetric(c.pages, prometheus.GaugeValue, float64(c.Wp.pages))
	ch <- prometheus.MustNewConstMetric(c.comments, prometheus.GaugeValue, float64(c.Wp.comments))
	ch <- prometheus.MustNewConstMetric(c.media, prometheus.GaugeValue, float64(c.Wp.media))
	ch <- prometheus.MustNewConstMetric(c.users, prometheus.GaugeValue, float64(c.Wp.users))
	ch <- prometheus.MustNewConstMetric(c.taxonomies, prometheus.GaugeValue, float64(c.Wp.taxonomies))
	ch <- prometheus.MustNewConstMetric(c.themes, prometheus.GaugeValue, float64(c.Wp.themes))
	ch <- prometheus.MustNewConstMetric(c.plugins, prometheus.GaugeValue, float64(c.Wp.plugins))
}

func NewWordpressCollector(w *Wordpress) *WordpressCollector {

	// debug
	fmt.Printf("NewWordpressCollector:\nSite: %v\nUse auth: %v\n", w.MonitoredWordpress, w.Auth.Use)

	return &WordpressCollector{
		Wp:  w,
		posts:      prometheus.NewDesc("wordpress_post_count", "WordPress posts count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		categories: prometheus.NewDesc("wordpress_category_count", "WordPress category count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		tags:       prometheus.NewDesc("wordpress_tag_count", "WordPress tags count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		pages:      prometheus.NewDesc("wordpress_page_count", "WordPress pages count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		comments:   prometheus.NewDesc("wordpress_comment_count", "WordPress comments count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		media:      prometheus.NewDesc("wordpress_media_count", "WordPress media files count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		users:      prometheus.NewDesc("wordpress_user_count", "WordPress users count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		taxonomies: prometheus.NewDesc("wordpress_taxonomies_count", "WordPress taxonomy count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		themes:     prometheus.NewDesc("wordpress_theme_count", "WordPress theme count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
		plugins:    prometheus.NewDesc("wordpress_plugin_count", "Wordpress plugin count", nil, prometheus.Labels{"instance": w.MonitoredWordpress}),
	}
}


