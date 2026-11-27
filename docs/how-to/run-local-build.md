# Run your own build

## Overview

This procedure examines a change to the Helm chart on your computer. It then puts your change on a test server. The change can be to the chart or to one of the images. Use this procedure when you prepare a pull request. [CONTRIBUTING.md](../../CONTRIBUTING.md) gives the other steps of a contribution.

## Before you start

- Clone your fork of the repository. Read [CONTRIBUTING.md](../../CONTRIBUTING.md#fork-and-clone).
- Install [Helm](https://helm.sh/docs/intro/install/) and [Docker](https://docs.docker.com/engine/install/) on your computer.
- To put changes on a test server, install OpenGovMail on the test server. Use the [installation guide](../installation-guide.md) and the `le-prod` issuer. Use a spare domain or a subdomain, for example `test.example.com`.

The test server must be a real server on the internet, not a cluster on your computer. Sign-in operates only with trusted Let's Encrypt certificates. Let's Encrypt gives certificates only to a server that it can connect to through the internet. The server must also be amd64 (x86_64), because the published images are for amd64 only.

## Examine a chart change

CI does these checks on each pull request that changes `charts/`. You can do them on all computers.

```bash
CHART=charts/opengovmail
SECRETS="--set global.ravenIdp.clientSecret=ci --set global.ravenIdp.directAuthSecret=ci --set thunderid.configuration.server.security.directAuthSecret=ci"
helm repo add nextcloud https://nextcloud.github.io/helm/
helm dependency build "$CHART"
helm lint "$CHART" -f "$CHART/values-example.yaml" $SECRETS
helm template ci "$CHART" -f "$CHART/values-example.yaml" $SECRETS > /dev/null
helm template ci "$CHART" --set global.domain=example.com $SECRETS > /dev/null
```

To read the output of your change, remove `> /dev/null`. Then add `--show-only templates/<file>.yaml`.

CI also installs the chart in a [kind](https://kind.sigs.k8s.io/) cluster and runs `helm test`. It uses `values-example.yaml` and the changes in `test/install/values.yaml`: a self-signed certificate issuer, placeholder secrets, and no Nextcloud or Tuwunel. The check shows that the pods start, that Raven gets a token from ThunderID, and that each service answers on its ports. It does not test sign-in from a mail app or mail delivery, because the cluster has no trusted certificates and no public address. If the check fails, its last step shows the pods, the events and the logs.

## Put a chart change on the test server

1. Push your branch.
2. On the test server, get your branch:

   ```bash
   git fetch origin
   git checkout <your-branch>
   helm dependency update ./charts/opengovmail
   ```

3. Do the upgrade command in [Upgrade OpenGovMail](upgrade.md) with your values file.

## Use your own image

The chart uses the published images. To use a changed image, do the steps that follow.

1. Build the image on the test server or on a different amd64 computer. For the Postfix image, use this command:

   ```bash
   docker build -t ghcr.io/<you>/opengovmail-smtp:dev -f services/smtp/Dockerfile services
   ```

   For the OpenDKIM image, use `services/dkim/Dockerfile`. You build Raven from the [raven repository](https://github.com/OpenGovMail/raven).

2. Make the image available to the cluster. Do one of these:
   - Push the image to a registry:

     ```bash
     docker push ghcr.io/<you>/opengovmail-smtp:dev
     ```

   - On a k3s server with one node, import the image into k3s:

     ```bash
     docker save ghcr.io/<you>/opengovmail-smtp:dev | sudo k3s ctr images import -
     ```

3. Set the chart to use the image:

   ```bash
   helm upgrade opengovmail ./charts/opengovmail -n opengovmail --reuse-values \
     --set postfix.image.repository=ghcr.io/<you>/opengovmail-smtp \
     --set postfix.image.tag=dev
   ```

   If you imported the image into k3s, also add `--set postfix.image.pullPolicy=Never`. For OpenDKIM, use `opendkim.image.*`. For Raven, use `raven.image.*`.

## Return to the published images

Do the upgrade command in [Upgrade OpenGovMail](upgrade.md) without `--reuse-values`. Your values file does not set the image. Thus, the chart uses its default images again.

## See also

- [CONTRIBUTING.md](../../CONTRIBUTING.md)
- [Overview of the architecture](../concepts/architecture.md)
- [Troubleshooting](../troubleshooting.md)
