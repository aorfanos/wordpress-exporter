# Prometheus WordPress exporter

Exposes WordPress site metrics using the WordPress Rest API.

## Installation

- Install the exporter (it doesn't need to be at the same machine as your site, as long as it can reach it by network): 

```console
docker pull saikolab/wordpress-exporter

docker run -d --publish 11011:11011 \
  -it saikolab/wordpress-exporter \
  -auth.user admin \
  -auth.pass adminpassword \
  -host https://aorfanos.com
```

- Put scrape configuration in your `prometheus.yml`:

```
scrape_configs:
  - job_name: wordpress_exporter
    honor_timestamps: true
    scrape_interval: 60s
    scrape_timeout: 15s
    metrics_path: /metrics
    scheme: http
    static_configs:
      - targets: ["<exporter-IP>:11011"]
```

## Metrics

| Metric name              | Type  | Description                 |
|--------------------------|-------|-----------------------------|
| wordpress_post_count     | Gauge |    WordPress posts count    |
| wordpress_category_count | Gauge |   WordPress category count  |
| wordpress_tag_count      | Gauge |     WordPress tags count    |
| wordpress_page_count     | Gauge |    WordPress pages count    |
| wordpress_comment_count  | Gauge |   WordPress comments count  |
| wordpress_media_count    | Gauge | WordPress media files count |
| wordpress_user_count     | Gauge |    WordPress users count    |
| wordpress_theme_count    | Gauge |    WordPress theme count    |
| wordpress_plugin_count   | Gauge |    Wordpress plugin count   |

