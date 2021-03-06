[![build and push](https://github.com/aorfanos/wordpress-exporter/actions/workflows/build-and-deploy.yaml/badge.svg)](https://github.com/aorfanos/wordpress-exporter/actions/workflows/build-and-deploy.yaml)

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

The outcome will look something like this:

![Grafana WordPress dashboard](https://i.imgur.com/e5A6UnM.png)

## Metrics

Warning: some endpoints (e.g. plugins, settings) require authentication. 
If you want to gather metrics for those endpoints you need to enable basic authentication on your
WordPress installation (via e.g. a plugin or custom code).

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

## Labels

- `exported_instance`: Site being monitored

## Todo 

- Provide config from file to monitor multiple hosts with one exporter
- Support native WordPress cookie authentication