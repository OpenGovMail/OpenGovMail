# Development Guide

## Local Development Mode

Build images locally to test your changes.

### Step 1: Enable Local Build
```bash
cp dev/examples/docker-compose.override.yml services/docker-compose.override.yml
```

### Step 2: Create Environment File
```bash
cp services/.env.example services/.env
nano services/.env
```

### Step 3: Build and Run
```bash
bash scripts/setup/setup.sh
bash scripts/service/start-opengovmail.sh
```

### Step 4: Switch Back to Production
```bash
rm services/docker-compose.override.yml
bash scripts/service/start-opengovmail.sh
```

## Running your own build on Kubernetes

The chart pulls the published images. To run your changes, build the image, make it available to the cluster, and point the chart at it. For the Postfix image:

```bash
docker build -t ghcr.io/<you>/opengovmail-smtp:dev -f services/smtp/Dockerfile services
docker push ghcr.io/<you>/opengovmail-smtp:dev
helm upgrade opengovmail ./charts/opengovmail -n opengovmail --reuse-values \
  --set postfix.image.repository=ghcr.io/<you>/opengovmail-smtp \
  --set postfix.image.tag=dev
```

The OpenDKIM image works the same way with `services/dkim/Dockerfile` and `opendkim.image.*`. Raven is built from the [raven repository](https://github.com/OpenGovMail/raven) and set with `raven.image.*`.

On a single-node k3s cluster you can skip the registry. Import the image straight into k3s and stop the cluster pulling it:

```bash
docker save ghcr.io/<you>/opengovmail-smtp:dev | sudo k3s ctr images import -
# then add to the helm upgrade above:
#   --set postfix.image.pullPolicy=Never
```

## Files

- `examples/docker-compose.override.yml` - Template for local builds
