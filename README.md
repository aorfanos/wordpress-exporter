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
  -it ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14 \
  -host http://example.com \
  -auth.basic false

# authenticated to wordpress api (return data from all endpoints)
docker run -d --publish 11011:11011 \
  -it ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14 \
  -auth.user wordpress-exporter \
  -auth.pass "Wdnh 7Wm0 UuxW 64DL y2lx r0It" \ # Application password for authenticated use
  -host http://example.com \
  -auth.basic true
```

- Put scrape configuration in your `prometheus.yml`:

```yml
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

## Environment Variables Support

Starting from this version, the exporter supports configuration via environment variables as a fallback when CLI arguments match their defaults. This is particularly useful for containerized deployments like Kubernetes, where passing CLI args can be cumbersome. CLI flags take precedence over environment variables, which in turn override built-in defaults.

**Backward Compatibility Note**: All existing CLI flags and their names/defaults remain unchanged. Environment variables are a new, optional feature that do not affect previous configurations—your current Docker runs with CLI args will continue to work exactly as before.

### Supported Environment Variables

| Variable | Description | Corresponding CLI Flag | Default |
|----------|-------------|-------------------------|---------|
| `WORDPRESS_EXPORTER_HOST` | WordPress site URL (required, e.g., `http://example.com`) | `-host` | "" (empty, fatal if not set) |
| `WORDPRESS_EXPORTER_PORT` | Port to expose metrics (integer) | `-port` | 11011 |
| `WORDPRESS_EXPORTER_AUTH_USER` | Basic auth username | `-auth.user` | "admin" |
| `WORDPRESS_EXPORTER_AUTH_PASS` | Basic auth password | `-auth.pass` | "admin" |
| `WORDPRESS_EXPORTER_AUTH_BASIC` | Enable basic auth (true/false) | `-auth.basic` | true |
| `WORDPRESS_EXPORTER_CONFIG_FILE` | Path to config file (unused currently) | `-config.file` | "wordpress-exporter.yml" |

### Usage with Docker (Environment Variables)

```console
# Run with environment variables (no CLI args needed)
docker run -d --publish 11011:11011 \
  -e WORDPRESS_EXPORTER_HOST=http://example.com \
  -e WORDPRESS_EXPORTER_AUTH_BASIC=true \
  -e WORDPRESS_EXPORTER_AUTH_USER=wordpress-exporter \
  -e WORDPRESS_EXPORTER_AUTH_PASS="Wdnh 7Wm0 UuxW 64DL y2lx r0It" \
  ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14
```

Mixing CLI and env vars is supported, with CLI taking priority.

## Multi-Host Configuration File Support

The exporter now supports monitoring multiple WordPress sites from a single instance using a YAML configuration file. If a valid config file is found, it takes precedence over single-host mode.

### Config File Format

Create a YAML file (default: `wordpress-exporter.yml`) with the following structure:

```yaml
wordpress-exporter:
  - https://wordpress.org
  - http://example.com
```

### Usage Modes

The exporter operates in two modes:

1. **Single Host Mode** (default): Uses CLI flags or environment variables to monitor one site
2. **Multi-Host Mode**: Uses config file to monitor multiple sites simultaneously

**Priority Order**: Config File → CLI Flags → Environment Variables → Defaults

Example multi-host usage:

```console
# Create config file
echo "wordpress-exporter:
  - https://wordpress.org  
  - https://example.com" > sites.yml

# Run with custom config file
docker run -d --publish 11011:11011 \
  -v $(pwd)/sites.yml:/app/sites.yml \
  ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14 \
  -config.file sites.yml \
  -auth.basic false
```

### Kubernetes Deployment Example

For Kubernetes, you can deploy the exporter using a Deployment with environment variables. Use Kubernetes Secrets for sensitive values like passwords.

**Single Host Example** `wordpress-exporter-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wordpress-exporter
spec:
  replicas: 1
  selector:
    matchLabels:
      app: wordpress-exporter
  template:
    metadata:
      labels:
        app: wordpress-exporter
    spec:
      containers:
      - name: exporter
        image: ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14
        ports:
        - containerPort: 11011
        env:
        - name: WORDPRESS_EXPORTER_HOST
          value: "http://example.com"
        - name: WORDPRESS_EXPORTER_PORT
          value: "11011"
        - name: WORDPRESS_EXPORTER_AUTH_BASIC
          value: "true"
        - name: WORDPRESS_EXPORTER_AUTH_USER
          valueFrom:
            secretKeyRef:
              name: wp-secrets
              key: username
        - name: WORDPRESS_EXPORTER_AUTH_PASS
          valueFrom:
            secretKeyRef:
              name: wp-secrets
              key: password
---
apiVersion: v1
kind: Service
metadata:
  name: wordpress-exporter
spec:
  selector:
    app: wordpress-exporter
  ports:
  - port: 11011
    targetPort: 11011
```

**Multi-Host Example** using ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: wp-exporter-config
data:
  wordpress-exporter.yml: |
    wordpress-exporter:
      - https://wordpress.org
      - http://example.com
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wordpress-exporter-multi
spec:
  replicas: 1
  selector:
    matchLabels:
      app: wordpress-exporter-multi
  template:
    metadata:
      labels:
        app: wordpress-exporter-multi
    spec:
      containers:
      - name: exporter
        image: ghcr.io/aorfanos/wordpress-exporter/wordpress-exporter:v0.0.14
        ports:
        - containerPort: 11011
        env:
        - name: WORDPRESS_EXPORTER_AUTH_BASIC
          value: "false"
        volumeMounts:
        - name: config-volume
          mountPath: /app/wordpress-exporter.yml
          subPath: wordpress-exporter.yml
      volumes:
      - name: config-volume
        configMap:
          name: wp-exporter-config
```

Apply with `kubectl apply -f wordpress-exporter-deployment.yaml`. Create a Secret for credentials: `kubectl create secret generic wp-secrets --from-literal=username=wordpress-exporter --from-literal=password="your-app-password"`.

Update your Prometheus config to scrape the Service.

## Todo

- [x] Provide config from file to monitor multiple hosts with one exporter
- [ ] Support native WordPress cookie authentication
