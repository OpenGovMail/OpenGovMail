# Run your own build

## Overview

This guide checks a change to the Helm chart on your own machine, then runs your change, whether to the chart or to one of the images, on a test server. Use it while working on a pull request. [CONTRIBUTING.md](../../CONTRIBUTING.md) covers the rest of the contribution workflow.

## Before you start

- Clone your fork of the repository. See [CONTRIBUTING.md](../../CONTRIBUTING.md#fork-and-clone).
- Install [Helm](https://helm.sh/docs/intro/install/) and [Docker](https://docs.docker.com/engine/install/) on your machine.
- For running changes, you need a test server with OpenGovMail installed, following the [installation guide](../installation-guide.md). Use a spare domain or a subdomain such as `test.example.com`, and the `le-prod` issuer.

The test server must be a real server rather than a cluster on your laptop. OpenGovMail needs trusted Let's Encrypt certificates for sign-in to work, and Let's Encrypt can only issue them to a server it can reach over the internet. The server must also be amd64 (x86_64), because the published images are built for amd64 only.

## Check a chart change

These are the same checks CI runs on every pull request that changes `charts/`. They work on any machine.

```bash
CHART=charts/opengovmail
SECRETS="--set global.ravenIdp.clientSecret=ci --set global.ravenIdp.directAuthSecret=ci --set thunderid.configuration.server.security.directAuthSecret=ci"
helm repo add nextcloud https://nextcloud.github.io/helm/
helm dependency build "$CHART"
helm lint "$CHART" -f "$CHART/values-example.yaml" $SECRETS
helm template ci "$CHART" -f "$CHART/values-example.yaml" $SECRETS > /dev/null
helm template ci "$CHART" --set global.domain=example.com $SECRETS > /dev/null
```

To read what your change renders, drop `> /dev/null` and add `--show-only templates/<file>.yaml`.

## Run a chart change on the test server

1. Push your branch, then on the test server check it out:

   ```bash
   git fetch origin
   git checkout <your-branch>
   helm dependency update ./charts/opengovmail
   ```

2. Run the upgrade command from [Upgrade OpenGovMail](upgrade.md), with your values file.

## Run your own image

The chart pulls the published images. To run a changed image, build it, make it available to the cluster, and point the chart at it.

1. Build the image on the test server, or on another amd64 machine. For the Postfix image:

   ```bash
   docker build -t ghcr.io/<you>/opengovmail-smtp:dev -f services/smtp/Dockerfile services
   ```

   The OpenDKIM image is built the same way from `services/dkim/Dockerfile`. Raven is built from the [raven repository](https://github.com/OpenGovMail/raven).

2. Make the image available to the cluster. Either push it to a registry:

   ```bash
   docker push ghcr.io/<you>/opengovmail-smtp:dev
   ```

   or, on a single-node k3s server, import it straight into k3s:

   ```bash
   docker save ghcr.io/<you>/opengovmail-smtp:dev | sudo k3s ctr images import -
   ```

3. Point the chart at it:

   ```bash
   helm upgrade opengovmail ./charts/opengovmail -n opengovmail --reuse-values \
     --set postfix.image.repository=ghcr.io/<you>/opengovmail-smtp \
     --set postfix.image.tag=dev
   ```

   If you imported the image into k3s, also add `--set postfix.image.pullPolicy=Never`. For OpenDKIM use `opendkim.image.*`, and for Raven `raven.image.*`.

## Go back to the published images

Run the upgrade command from [Upgrade OpenGovMail](upgrade.md) without `--reuse-values`. Your values file doesn't set the image, so the chart's defaults return.

## See also

- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [Overview of the architecture](../concepts/architecture.md)
- [Troubleshooting](../troubleshooting.md)
