# Installation guide

## Introduction

This guide installs OpenGovMail with its Helm chart, explains each setting the install needs, and covers upgrades and removal. For the shortest path to a working server, follow the [quickstart](quickstart.md) instead.

## Installation types

OpenGovMail is installed with Helm on Kubernetes. The chart is the same in every case. What changes is the cluster underneath.

| Type | Cluster | Use it for |
|---|---|---|
| Production, single node | k3s on one amd64 server | Running real mail. This is the tested setup. |
| Production, other clusters | Any Kubernetes with an ingress controller and cert-manager | Running real mail on infrastructure you already have. Needs the changes in [Other clusters](#other-clusters). |
| Local development | k3d on an amd64 machine | Trying out changes. Mail to and from the internet doesn't work. See [Run your own build](how-to/run-local-build.md). |

## Overview

An install has five stages:

1. Point DNS at the server and open its firewall ports.
2. Install cert-manager and the Let's Encrypt issuers, once per cluster.
3. Write a values file for your domain.
4. Install the chart.
5. Publish the DKIM record that OpenDKIM generates.

## System requirements

| Requirement | Minimum |
|---|---|
| CPU architecture | amd64 (x86_64). The published OpenGovMail images are built for amd64 only. |
| CPU | 1 CPU for email. 2 CPUs with collaboration turned on. |
| Memory | 4 GB |
| Network | A static public IP address, with inbound and outbound port 25 allowed |
| Kubernetes | k3s, or another cluster with an ingress controller. Postfix and Raven bind host ports, so mail runs on a single node. |
| Tools | [Helm](https://helm.sh/docs/intro/install/), [kubectl](https://kubernetes.io/docs/tasks/tools/), [Git](https://git-scm.com/), `openssl` |
| cert-manager | v1.20.0, installed in Step 3 |

With collaboration turned on and only 1 CPU, the pods' CPU requests fill the node and some pods, including cert-manager's, never get scheduled.

## Before you begin

- You need a domain whose DNS you control, and access to your hosting provider's reverse DNS (PTR) settings.
- Some cloud providers block port 25 by default. Ask yours to open it before you start.
- In this guide, `example.com` stands for your domain and `12.34.56.78` for the server's public IP address. `example-com` is the same domain with dots replaced by dashes, as used in Secret names.
- [Glossary](glossary.md) explains the terms used here.

## Installation steps

### Step 1: Add DNS records

| Type | Name | Value |
|---|---|---|
| A | mail | 12.34.56.78 |
| A | example.com | 12.34.56.78 |
| MX | example.com | mail.example.com |
| TXT | example.com | "v=spf1 ip4:12.34.56.78 ~all" |
| TXT | _dmarc | "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com" |
| PTR | 12.34.56.78 | mail.example.com |

PTR records are usually set through your hosting provider. If you plan to turn on collaboration, add an A record for each of its hostnames too. [Hostnames and ports](reference/hostnames-and-ports.md) lists them.

You add the DKIM record after the install, in Step 7, because the key is generated on first start.

### Step 2: Open firewall ports

Open these inbound ports on the server:

| Port | Used by |
|---|---|
| 25/tcp | Mail from other servers |
| 80/tcp | Let's Encrypt HTTP-01 challenges |
| 443/tcp | Web services behind the ingress |
| 587/tcp | Mail clients sending mail |
| 993/tcp | Mail clients reading mail (IMAPS) |
| 8090/tcp | ThunderID console, sign-in and token endpoints |

Collaboration needs more ports. See [Hostnames and ports](reference/hostnames-and-ports.md).

### Step 3: Install Kubernetes, Helm and cert-manager

1. Install k3s and give your user access to the cluster:

   ```bash
   curl -sfL https://get.k3s.io | sh -
   mkdir -p ~/.kube
   sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
   sudo chown "$USER" ~/.kube/config
   ```

2. Install Helm:

   ```bash
   curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
   ```

3. Install cert-manager. It is cluster infrastructure, so install it once, outside OpenGovMail:

   ```bash
   helm install cert-manager oci://quay.io/jetstack/charts/cert-manager \
     --version v1.20.0 \
     --namespace cert-manager \
     --create-namespace \
     --set crds.enabled=true
   kubectl get pods -n cert-manager
   ```

   All pods should be `Running`.

### Step 4: Create the certificate issuers

Clone the repository and create two Let's Encrypt ClusterIssuers:

```bash
git clone https://github.com/OpenGovMail/OpenGovMail.git
cd OpenGovMail
kubectl apply -f infra/clusterissuers-http01.yaml
kubectl get clusterissuer
```

Both should show `READY` as `True`.

- `le-staging` issues untrusted certificates with no rate limits, for testing.
- `le-prod` issues trusted certificates, limited to 5 per set of names per week.

Both use the HTTP-01 challenge: Let's Encrypt checks each hostname over port 80, so every hostname must already point at the server. The issuers send challenges through Traefik, the ingress controller k3s ships with.

### Step 5: Write a values file

Copy [values-example.yaml](../charts/opengovmail/values-example.yaml) to a file outside the repository, such as `~/my-values.yaml`, and change it for your domain:

- Set `global.domain` to your domain. Everything else, such as the mail hostname and certificate names, comes from it. The install fails with `global.domain is required` without it.
- Set `thunderid.configuration.server.publicUrl` and `thunderid.configuration.gateClient.hostname` to the same domain. The ThunderID chart can't read `global.domain`, so it repeats here, and the install fails if they disagree.
- Change `thunderid.declarativeResources.secret.name` to `example-com-tls` for your domain.
- Remove the `nextcloud` and `tuwunel` blocks for an email-only install, or change their hostnames to yours. [Turn on collaboration](how-to/enable-collaboration.md) explains them.

`global.tls.issuer` chooses the issuer. The example uses `le-prod`. To test DNS and the challenge first without using up the production limit, start with `le-staging` and change to `le-prod` with an upgrade once certificates are issued. Staging certificates aren't trusted, so sign-in and Nextcloud don't work until you switch.

Keep secrets out of this file. They go on the command line in the next step. [Helm values](reference/helm-values.md) describes the other settings.

### Step 6: Install the chart

Fetch the ThunderID and Nextcloud charts, which aren't stored in Git:

```bash
helm dependency update ./charts/opengovmail
```

Then install. This command creates Raven's ThunderID credentials, so run it only once per namespace:

```bash
DIRECT_AUTH_SECRET=$(openssl rand -hex 32)
helm install opengovmail ./charts/opengovmail \
  -f ~/my-values.yaml \
  --set-string global.thunderAdmin.password='<your-admin-password>' \
  --set global.ravenIdp.clientSecret=$(openssl rand -hex 32) \
  --set global.ravenIdp.directAuthSecret=$DIRECT_AUTH_SECRET \
  --set thunderid.configuration.server.security.directAuthSecret=$DIRECT_AUTH_SECRET \
  --namespace opengovmail \
  --create-namespace
```

| Value | What it is |
|---|---|
| `global.thunderAdmin.password` | Password for the ThunderID `admin` user. If you leave it out, the chart generates one. Read it back with `kubectl get secret opengovmail-thunder-admin -n opengovmail -o jsonpath='{.data.password}' \| base64 -d` |
| `global.ravenIdp.clientSecret` | Raven's secret for signing in to ThunderID. Set it on the first install only. Later upgrades reuse the value in the cluster. |
| `global.ravenIdp.directAuthSecret` and `thunderid.configuration.server.security.directAuthSecret` | The secret that protects ThunderID's password check. Both must be the same value. |

The identity running Helm needs `get` on Secrets and `list` on Deployments in the namespace, because the chart reads them while rendering. The built-in `admin` and `edit` roles have both.

### Step 7: Add the DKIM record

When the pods are running, print the DKIM record:

```bash
kubectl exec -n opengovmail statefulset/opengovmail-opendkim -c opendkim -- \
  cat /etc/dkimkeys/example.com/mail.txt
```

Add it as a TXT record named `mail._domainkey`:

| Type | Name | Value |
|---|---|---|
| TXT | mail._domainkey | "v=DKIM1; h=sha256; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ..." |

> [!Important]
> The key is printed in several quoted parts. Join them into one value, without the quotes and the spaces between them.

## Verify installation

1. Check the pods and certificates:

   ```bash
   kubectl get pods -n opengovmail
   kubectl get certificate -n opengovmail
   ```

   Every pod should be `Running` or `Completed`. `mail-example-com-tls` and `example-com-tls` should show `READY` as `True`. If a certificate stays not ready, run `kubectl describe certificate <name> -n opengovmail`.

2. Open `https://example.com:8090/console` and sign in as `admin`.
3. [Add a user](how-to/add-users.md) and [connect a mail client](user-guide/mail-client.md).
4. Check the setup from outside:
   - [MXToolbox](https://mxtoolbox.com/SuperTool.aspx) checks your DNS records, blocklists and mail server.
   - [mail-tester](https://www.mail-tester.com/) scores a test email out of 10 for deliverability.

## Post installation

### Configuration options

- [Helm values](reference/helm-values.md) lists the settings you are most likely to change.
- [Turn on collaboration](how-to/enable-collaboration.md) adds Nextcloud, Collabora, Talk and Matrix.
- [Add users](how-to/add-users.md) in the ThunderID console.

### Other clusters

The chart is tested on single-node k3s. On another cluster:

- Change `ingressClassName` in `infra/clusterissuers-http01.yaml`, and the ingress class in your values, from `traefik` to your ingress controller.
- Postfix binds ports 25 and 587, and Raven ports 143 and 993, as host ports. Run them on one node whose public IP the MX record points at.
- Raven reaches ThunderID through the server's public name and port 8090. The cluster must route traffic from a pod to the node's own external IP (hairpin NAT). k3s does this.

### Upgrade options

[Upgrade OpenGovMail](how-to/upgrade.md) explains how to move to a new chart version. [CHANGELOG.md](../CHANGELOG.md) lists what changed in each version, and any steps an upgrade needs.

### Downgrade options

Helm can roll back to an earlier revision of the release:

```bash
helm history opengovmail -n opengovmail
helm rollback opengovmail <revision> -n opengovmail
```

A rollback restores the Kubernetes resources, not data. Check [CHANGELOG.md](../CHANGELOG.md) before going back past a release that changed stored data.

### Uninstallation options

```bash
helm uninstall opengovmail -n opengovmail
```

> [!Warning]
> Uninstalling deletes the volumes the chart creates directly. That includes Raven's, which holds every mailbox. Back up anything you need first; `kubectl get pvc -n opengovmail` lists the volumes. Volumes that belong to StatefulSets, such as SeaweedFS (attachments) and OpenDKIM (signing keys), stay until you delete them.

To delete everything that remains, delete the namespace:

```bash
kubectl delete namespace opengovmail
```

cert-manager and the issuers are shared by the cluster. Remove them separately if nothing else uses them.

## Troubleshooting

[Troubleshooting](troubleshooting.md) covers the common problems: a certificate that isn't issued, ThunderID's setup job not finishing, DKIM failures and mail going to spam.

## Next steps

- [Add users](how-to/add-users.md)
- [Turn on collaboration](how-to/enable-collaboration.md)
- [User guide](user-guide/README.md) for the people who will use it

## Product version history

[CHANGELOG.md](../CHANGELOG.md) records every release, following [Semantic Versioning](https://semver.org/). The chart version and the OpenGovMail image tags move together.

## Definition of terms

See the [glossary](glossary.md).
