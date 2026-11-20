# OpenGovMail Helm chart

This directory holds the OpenGovMail Helm chart. To install it, follow the [quickstart](../docs/quickstart.md) or the [installation guide](../docs/installation-guide.md). This page is for people changing the chart.

## Structure

- `opengovmail/`: the umbrella chart. Its `templates/` hold resources that span services, such as certificates, the ThunderID bootstrap data and the Nextcloud integration.
- `opengovmail/charts/<service>/`: one subchart per service: `opendkim`, `postfix`, `raven`, `redis`, `rspamd`, `seaweedfs`, `tuwunel` and `unbound`.
- ThunderID and Nextcloud are upstream charts, listed as dependencies in [opengovmail/Chart.yaml](opengovmail/Chart.yaml) and fetched with `helm dependency update`. They aren't stored in Git.
- [opengovmail/values-example.yaml](opengovmail/values-example.yaml) is a complete values file for a deployment. CI lints and renders it.

## Conventions

- Every subchart has a `values.yaml` with safe defaults, its `templates/`, and a `README.md` with install and operations notes.
- Deployment values, such as the domain and hostnames, belong in the deployment's own values file, not in the chart.
- Service names are fixed in `global.serviceNames` and don't depend on the release name, so services can find each other.
- Settings that two services must agree on, such as the S3 credentials, live under `global`.

## Adding a service

1. Create `opengovmail/charts/<service>/` following an existing subchart, such as `opendkim`.
2. Add it to `dependencies` in [opengovmail/Chart.yaml](opengovmail/Chart.yaml) with `condition: <service>.enabled`.
3. Add its `enabled` flag and any shared settings to [opengovmail/values.yaml](opengovmail/values.yaml).
4. Document it in [Technology stack](../docs/reference/technology-stack.md) and, if users set it, [Helm values](../docs/reference/helm-values.md).

## ThunderID

ThunderID comes from `oci://ghcr.io/thunder-id/helm-charts/thunderid` and is configured under the `thunderid:` key. It runs as a single pod on SQLite, with a setup job that creates its database.

### Bootstrap

The setup job applies ThunderID's own defaults (the `default` organization unit, the `Person` user type, flows, the Administrator role and the console app), then the umbrella's additions from [opengovmail/files/thunderid-bootstrap/](opengovmail/files/thunderid-bootstrap/):

- Raven System: the machine-to-machine application Raven signs in as.
- Raven System Role: gives that application the `system` permission it needs to read organization units, users and groups.

The umbrella renders these into the `thunder-bootstrap` ConfigMap. The setup job runs on install only, so changing the bootstrap data, or moving from the old asgardeo Thunder chart, means installing into a fresh namespace. The render fails if the namespace still runs the old Thunder chart.

Mail clients sign in with a password. OAuth sign-in for mail clients isn't set up yet.

### Credentials

| Credential | Where it comes from |
|---|---|
| Admin password (`opengovmail-thunder-admin`, username `admin`) | Set it with `global.thunderAdmin.password`. If left blank, the chart generates one with `lookup`, which only works for a real `helm upgrade --install`. A `helm template \| kubectl apply` (GitOps) render can't see the cluster and mints a new password on every render. |
| `global.ravenIdp.clientSecret` | Supply on first install. Reused from the cluster afterwards. A different value later fails the render, because the setup job never registers it again. |
| `global.ravenIdp.directAuthSecret` | Supply on first install, with the same value in `thunderid.configuration.server.security.directAuthSecret`. On upgrades, pass ThunderID's copy again. |

### Permissions

The chart reads the cluster while rendering: it gets the `opengovmail-thunder-admin` and `opengovmail-raven-idp` Secrets and lists Deployments in the release namespace, to refuse an in-place upgrade from the old Thunder chart. The identity running Helm needs `get` on Secrets and `list` on Deployments there. The built-in `admin` and `edit` roles include both.

### Networking

Raven reaches ThunderID on port `8090` using the public name, which matches the JWT issuer. [opengovmail/templates/thunder-lb.yaml](opengovmail/templates/thunder-lb.yaml) exposes `8090` on the node for the console and sign-in pages. `thunderid.ingress` is turned off.
