# OpenGovMail Helm Charts

This directory contains Helm charts for OpenGovMail services.

## Structure

- `opengovmail/`: Umbrella chart that aggregates service subcharts.
- `opengovmail/charts/opendkim/`: OpenDKIM service chart (first migrated service).

## Conventions

- Service charts live under `opengovmail/charts/<service-name>/`.
- Every service chart should include:
  - `values.yaml` with safe defaults.
  - `templates/` for Kubernetes resources.
  - `README.md` with install and operations guidance.
- Environment overlays should be kept in umbrella chart values files (`values-dev.yaml`, `values-prod.yaml`).

## Next Services

When adding future services (smtp, rspamd, raven, etc.), repeat the OpenDKIM chart pattern and add them as dependencies in the umbrella chart.

## Install cert-manager (once per cluster)

cert-manager is cluster infrastructure. Install it once, outside of OpenGovMail.

```bash
helm install cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.20.0 \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
```

Verify:

```bash
kubectl get pods -n cert-manager
# All pods should be Running
```

---

## Step 2 — Bootstrap ClusterIssuers (once per cluster)

The bootstrap script creates:
- A Cloudflare API token Secret in the `cert-manager` namespace
- `le-staging` ClusterIssuer (Let's Encrypt staging — untrusted, no rate limits)
- `le-prod` ClusterIssuer (Let's Encrypt prod — trusted, rate limited)

You will need:
- A Cloudflare API token scoped to `Zone:Read` and `DNS:Edit` for your domain
- An email address for Let's Encrypt notifications

Run from the repo root:

```bash
bash infra/bootstrap.sh
```

Verify:

```bash
kubectl get clusterissuer
# Both le-staging and le-prod should show READY=True
```

---

## Step 3 — Configure values.yaml

`global.domain` is the **one required input** — everything (mail hostname, cert
Secret, service config) derives from it. Installing without it fails fast:
`global.domain is required`.

```yaml
global:
  domain: yourdomain.com   # REQUIRED
  tls:
    enabled: true
    issuer: le-staging     # use le-staging first, switch to le-prod once verified
    renewBefore: "720h"    # renew 30 days before expiry
    # domains: []          # optional — defaults to [ mail.<domain>, <domain> ].
    #                        Override only to add extra SANs.
```

With `domains` left empty, certificates are minted for:
- `mail.yourdomain.com` (+ `*.mail.yourdomain.com`) — the Secret postfix mounts
- `yourdomain.com` (+ `*.yourdomain.com`)

### Staging vs Production

Always test with `le-staging` first. Staging issues real but **untrusted** certificates —
browsers will show a warning, but the full DNS challenge and issuance flow is verified.

Once the staging certificate shows `READY=True`, switch to `le-prod`:

```yaml
tls:
  issuer: le-prod
```

Let's Encrypt prod has a rate limit of **5 certificates per domain per week**.
Burning through this with misconfigured charts is a common mistake — staging prevents it.

---

## Step 4 — Install OpenGovMail

From the repo root. `global.domain` is mandatory, and a first install also needs
the ThunderID admin password and raven's two ThunderID secrets (see
[ThunderID identity server](#thunderid-identity-server)):

```bash
DIRECT_AUTH_SECRET=$(openssl rand -hex 32)
helm upgrade --install opengovmail ./charts/opengovmail \
  --set global.domain=yourdomain.com \
  --set global.thunderAdmin.password=<your-password> \
  --set global.ravenIdp.clientSecret=$(openssl rand -hex 32) \
  --set global.ravenIdp.directAuthSecret=$DIRECT_AUTH_SECRET \
  --set thunderid.configuration.server.security.directAuthSecret=$DIRECT_AUTH_SECRET \
  --set thunderid.configuration.server.publicUrl=https://yourdomain.com:8090 \
  --set thunderid.configuration.gateClient.hostname=yourdomain.com \
  --namespace opengovmail \
  --create-namespace
```

With a values overlay (e.g. dev — TLS/SASL off, standalone postfix, domain preset):

```bash
helm upgrade --install opengovmail ./charts/opengovmail \
  -f charts/opengovmail/values-dev.yaml \
  --namespace opengovmail \
  --create-namespace
```

Verify certificates:

```bash
kubectl get certificate -n opengovmail
# mail-<domain>-tls and <domain>-tls should show READY=True
```

If not ready yet, check progress:

```bash
kubectl describe certificate <name> -n opengovmail
```
---

## ThunderID identity server

ThunderID is pulled in as an upstream OCI dependency
(`oci://ghcr.io/thunder-id/helm-charts/thunderid`, see
[opengovmail/Chart.yaml](opengovmail/Chart.yaml) for the version) and configured under
the `thunderid:` key in the umbrella values. It runs as a single pod on SQLite, with
a one-time setup job that bootstraps the database. Run
`helm dependency update ./charts/opengovmail` before the first install.

### Bootstrap

The setup job applies ThunderID's own defaults (the `default` organization unit,
the `Person` user type, flows, the Administrator role and the console app), then
the umbrella's additions from
[opengovmail/files/thunderid-bootstrap/](opengovmail/files/thunderid-bootstrap/):

- **Raven System**: the machine-to-machine application raven authenticates as
- **Raven System Role**: grants that application the `system` permission it needs
  to read organization units, users and groups

The umbrella renders these into the `thunder-bootstrap` ConfigMap itself. Nothing
needs to be created by hand before installing.

The setup job runs **on install only**. Changing the bootstrap data, or moving from
the old asgardeo Thunder chart, means installing into a fresh namespace. The render
fails if the namespace still runs the old Thunder chart.

Mail clients sign in with a password. OAuth sign-in for mail clients is not
bootstrapped yet.

### Credentials

| Credential | Where it comes from |
|---|---|
| Admin password (`opengovmail-thunder-admin`, username `admin`) | **Set it explicitly** with `global.thunderAdmin.password`. If left blank, the chart generates one using `lookup`, which only works for a real `helm upgrade --install`. A `helm template \| kubectl apply` (GitOps) render can't see the cluster and would mint a new password on every render. Read a generated value with `kubectl get secret opengovmail-thunder-admin -n opengovmail -o jsonpath='{.data.password}' \| base64 -d` |
| `global.ravenIdp.clientSecret` | **Supply on first install.** Reused from the cluster afterwards. Supplying a different value later fails the render, because the setup job never re-registers it |
| `global.ravenIdp.directAuthSecret` | **Supply on first install**, with the same value in `thunderid.configuration.server.security.directAuthSecret`. On upgrades, pass the ThunderID copy again: `kubectl get secret opengovmail-raven-idp -n opengovmail -o jsonpath='{.data.directAuthSecret}' \| base64 -d` |

### Networking

Raven reaches ThunderID on port `8090` using the public name, which matches the JWT
issuer. [opengovmail/templates/thunder-lb.yaml](opengovmail/templates/thunder-lb.yaml)
exposes `8090` on the node for the console and gate UIs. `thunderid.ingress` is
disabled. The dev overlay ([values-dev.yaml](opengovmail/values-dev.yaml)) disables
ThunderID and raven for a postfix-only bring-up.
---