[![Snyk scans](https://github.com/aorfanos/wordpress-exporter/actions/workflows/security-scanning.yaml/badge.svg)](https://github.com/aorfanos/wordpress-exporter/actions/workflows/security-scanning.yaml)
[![Build and Push Go Code to Github Container Registry](https://github.com/aorfanos/wordpress-exporter/actions/workflows/build.yaml/badge.svg)](https://github.com/aorfanos/wordpress-exporter/actions/workflows/build.yaml)
# Prometheus WordPress exporter

Exposes WordPress site metrics using the WordPress Rest API.

## Installation

- Install the exporter (it doesn't need to be at the same machine as your site, as long as it can reach it by network): 

```console
docker pull ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter

# run plain, without authentication
docker run -d --publish 11011:11011 \
  -it ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.8 \
  -host http://example.com \
  -auth.basic false

# authenticated to wordpress api (return data from all endpoints)
docker run -d --publish 11011:11011 \
  -it ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.8 \
  -auth.user wordpress-exporter \
  -auth.pass "Wdnh 7Wm0 UuxW 64DL y2lx r0It" \ # Application password for authenticated use
  -host http://example.com \
  -auth.basic true
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

### Authenticating to WordPress API

Since version 5.6, WordPress has introduced the ability to use [Application Passwords](https://developer.wordpress.org/rest-api/using-the-rest-api/authentication/#basic-authentication-with-application-passwords) for performing basic authentication. This new feature is now the recommended way to securely access the REST API, making it easier for developers to build applications and services that integrate with WordPress while maintaining a high level of security.

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
| wordpress_taxonomy_count | Gauge |    WordPress taxonomy count |
| wordpress_theme_count    | Gauge |    WordPress theme count    |
| wordpress_plugin_count   | Gauge |    Wordpress plugin count   |

## Labels

- `exported_instance`: Site being monitored

## Todo 

- Provide config from file to monitor multiple hosts with one exporter
- Support native WordPress cookie authentication
