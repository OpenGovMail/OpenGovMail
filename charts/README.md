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

From the repo root (`global.domain` is mandatory — set it in `values.yaml` or via `--set`):

```bash
helm upgrade --install opengovmail ./charts/opengovmail \
  --set global.domain=yourdomain.com \
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

## Thunder identity server

Thunder is pulled in as an upstream OCI dependency
(`oci://ghcr.io/asgardeo/helm-charts/thunder`, version `0.32.0`) and configured
under the `thunder:` key in the umbrella values. The defaults in
[opengovmail/values.yaml](opengovmail/values.yaml) reproduce the docker-compose `thunder`
setup: the `0.32.0` image, a shared SQLite database seeded from the image (via an
init container, replacing the compose `thunder-db-init`), a one-time setup job
(compose `thunder-setup`), a single pod, and the bootstrap scripts from
[scripts/thunder](../scripts/thunder).

### Bootstrap ConfigMap (required before install)

Thunder's setup job runs as a Helm **pre-install hook**, so the bootstrap
ConfigMap it references must already exist in the namespace. Create it from the
canonical scripts (they are not duplicated into the chart):

```bash
scripts/thunder/create-bootstrap-configmap.sh opengovmail          # <namespace> [configmap-name]
```

This mounts `01-default-resources.sh` and `02-sample-resources.sh` via `subPath`,
preserving the image's default bootstrap scripts (including `common.sh`, which
`01`/`02` source). Re-run it to push script changes; it is idempotent.

### Admin credentials (set before install)

The admin **username** is `admin` (in `thunder.setup.env`). The **password**
comes from the `opengovmail-thunder-admin` Secret, which the umbrella chart
itself creates (`templates/thunder-admin-secret.yaml`) — set it explicitly
via `global.thunderAdmin.password` in a values file or `--set` before
installing:

```yaml
global:
  thunderAdmin:
    password: "<your-password>"
```

Set it explicitly rather than leaving it blank: if unset, the chart
auto-generates a password using `lookup` against the live cluster at render
time, which only works for a real `helm upgrade --install`. A
`helm template | kubectl apply` (GitOps) render can't see the cluster and
would mint a fresh random password on every render, silently rotating the
credential. Setting the password explicitly avoids this.

If you do leave it blank, the generated value can still be read back with:

```bash
kubectl get secret opengovmail-thunder-admin -n opengovmail \
  -o jsonpath='{.data.password}' | base64 -d
```

### Install

```bash
# 1. Bootstrap ConfigMap first (pre-install hook dependency)
scripts/thunder/create-bootstrap-configmap.sh opengovmail

# 2. Install the umbrella (Thunder enabled by default)
helm dependency update ./charts/opengovmail
helm upgrade --install opengovmail ./charts/opengovmail \
  --set global.domain=yourdomain.com \
  --set global.thunderAdmin.password=<your-password> \
  --namespace opengovmail --create-namespace
```

Thunder is reached in-cluster on its Service at port `8090` (e.g. raven ->
`thunder:8090`), matching the compose network. External ingress is disabled by
default; enable `thunder.ingress` and point it at your domain if you need the
console/gate UIs exposed. The dev overlay
([values-dev.yaml](opengovmail/values-dev.yaml)) disables Thunder for a postfix-only
bring-up.
---