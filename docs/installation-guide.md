# Installation guide

## Introduction

This guide installs OpenGovMail with its Helm chart. It tells you about each setting that the installation uses. It also gives the procedures to upgrade and remove OpenGovMail. For the shortest procedure, use the [quickstart](quickstart.md).

## Installation types

You install OpenGovMail with Helm on Kubernetes. The chart is the same for each type of installation. Only the cluster is different.

| Type | Cluster | Use |
|---|---|---|
| Production, one node | k3s on one amd64 server | Real mail. We test this installation type. |
| Production, other clusters | A Kubernetes cluster with an ingress controller and cert-manager | Real mail on a cluster that you have already. Do the changes in [Other clusters](#other-clusters). |
| Development | k3s on an amd64 test server | Tests of your changes. Read [Run your own build](how-to/run-local-build.md). |

## Overview

The installation has five stages:

1. Point DNS to the server and open the firewall ports.
2. Install cert-manager and the Let's Encrypt issuers. Do this one time for each cluster.
3. Write a values file for your domain.
4. Install the chart.
5. Add the DKIM record that OpenDKIM makes.

## System requirements

| Requirement | Minimum |
|---|---|
| CPU architecture | amd64 (x86_64). The published OpenGovMail images are for amd64 only. |
| CPU | 1 CPU for email. 2 CPUs with the collaboration services. |
| Memory | 4 GB |
| Network | A static public IP address. Inbound and outbound connections on port 25 must be possible. |
| Kubernetes | k3s, or a different cluster with an ingress controller. Postfix and Raven use host ports, thus mail operates on one node only. |
| Tools | [Helm](https://helm.sh/docs/intro/install/), [kubectl](https://kubernetes.io/docs/tasks/tools/), [Git](https://git-scm.com/), `openssl` |
| cert-manager | v1.20.0. You install it in Step 3. |

Do not add the collaboration services to a server with 1 CPU. The pods then request more CPU than the node has. Kubernetes cannot start some pods, for example the cert-manager pods.

## Before you begin

- You must have a domain and access to its DNS records.
- You must have access to the reverse DNS (PTR) settings at your hosting provider.
- Some cloud providers close port 25. If your provider closes port 25, ask the provider to open it before you start.
- In this guide, `example.com` is your domain and `12.34.56.78` is the public IP address of your server. `example-com` is your domain with dashes in place of the dots. Secret names use this form.
- The [glossary](glossary.md) tells you about the terms in this guide.

## Installation steps

### Step 1: Add DNS records

1. Add these records:

   | Type | Name | Value |
   |---|---|---|
   | A | mail | 12.34.56.78 |
   | A | example.com | 12.34.56.78 |
   | MX | example.com | mail.example.com |
   | TXT | example.com | "v=spf1 ip4:12.34.56.78 ~all" |
   | TXT | _dmarc | "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com" |
   | PTR | 12.34.56.78 | mail.example.com |

   Usually, you set the PTR record at your hosting provider.

2. If you will add the collaboration services, add an A record for each of their hostnames. [Hostnames and ports](reference/hostnames-and-ports.md) gives the list.

You add the DKIM record in Step 7, after the installation. OpenDKIM makes the key when it starts for the first time.

### Step 2: Open firewall ports

Open these inbound ports on the server:

| Port | Use |
|---|---|
| 25/tcp | Mail from other servers |
| 80/tcp | Let's Encrypt HTTP-01 challenges |
| 443/tcp | Web services behind the ingress |
| 587/tcp | Mail apps that send mail |
| 993/tcp | Mail apps that read mail (IMAPS) |
| 8090/tcp | ThunderID console, sign-in and tokens |

The collaboration services use more ports. Read [Hostnames and ports](reference/hostnames-and-ports.md).

### Step 3: Install Kubernetes, Helm and cert-manager

1. Install k3s. Give your user access to the cluster:

   ```bash
   curl -sfL https://get.k3s.io | sh -
   mkdir -p ~/.kube
   sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
   sudo chown "$USER" ~/.kube/config
   export KUBECONFIG=~/.kube/config
   echo 'export KUBECONFIG=~/.kube/config' >> ~/.bashrc
   ```

   The `kubectl` command that k3s installs reads `/etc/rancher/k3s/k3s.yaml` by default, which only `root` can read. The two `KUBECONFIG` lines make `kubectl` use your copy in this session and in new sessions.

2. Install Helm:

   ```bash
   curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
   ```

3. Install cert-manager. Install it one time for the cluster, not as a part of OpenGovMail:

   ```bash
   helm install cert-manager oci://quay.io/jetstack/charts/cert-manager \
     --version v1.20.0 \
     --namespace cert-manager \
     --create-namespace \
     --set crds.enabled=true
   ```

4. Make sure that all the cert-manager pods show `Running`:

   ```bash
   kubectl get pods -n cert-manager
   ```

### Step 4: Create the certificate issuers

1. Get the repository:

   ```bash
   git clone https://github.com/OpenGovMail/OpenGovMail.git
   cd OpenGovMail
   ```

2. Create the two Let's Encrypt ClusterIssuers:

   ```bash
   kubectl apply -f infra/clusterissuers-http01.yaml
   ```

3. Make sure that the two issuers show `True` in the `READY` column:

   ```bash
   kubectl get clusterissuer
   ```

The two issuers are:

- `le-staging`: Certificates for tests. Browsers and mail apps do not trust them. There is no rate limit.
- `le-prod`: Trusted certificates. The limit is 5 certificates for each set of names in one week.

The two issuers use the HTTP-01 challenge. Let's Encrypt connects to each hostname on port 80. Thus, each hostname must point to the server before you install. The issuers send the challenges through Traefik. Traefik is the ingress controller that k3s installs.

### Step 5: Write a values file

1. Copy [values-example.yaml](../charts/opengovmail/values-example.yaml) to a file outside the repository, for example `~/my-values.yaml`.
2. Set `global.domain` to your domain. The chart uses this value for the mail hostname, the certificate names and other settings. If `global.domain` has no value, the installation stops with the error `global.domain is required`.
3. Set `thunderid.configuration.server.publicUrl` and `thunderid.configuration.gateClient.hostname` to the same domain. The ThunderID chart cannot read `global.domain`, thus you must write the domain again here. If the values are different, the installation stops.
4. Do not change `thunderid.declarativeResources.secret.name` from `thunderid-tls`. With this name, the chart gets the ThunderID certificate before the ThunderID setup job starts.
5. Do one of these:
   - For an installation with only email, remove the `nextcloud` and `tuwunel` blocks.
   - For the collaboration services, change their hostnames to your hostnames. [Add the collaboration services](how-to/enable-collaboration.md) tells you about them.

`global.tls.issuer` sets the issuer. The example uses `le-prod`. To test DNS and the challenge before you use the production limit, start with `le-staging`. When the certificates are ready, change to `le-prod` and do an upgrade. Sign-in and Nextcloud do not operate with `le-staging` certificates.

Do not put secrets in this file. You give the secrets on the command line in Step 6. [Helm values](reference/helm-values.md) tells you about the other settings.

### Step 6: Install the chart

1. Get the ThunderID and Nextcloud charts. Git does not keep these charts:

   ```bash
   helm dependency update ./charts/opengovmail
   ```

2. Install the chart. Replace `<your-admin-password>` with a password for the ThunderID administrator:

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

> [!Caution]
> Do this command only one time for each namespace. The command makes the ThunderID credentials for Raven. A second installation with new values does not operate.

| Value | Description |
|---|---|
| `global.thunderAdmin.password` | The password for the ThunderID `admin` user. If you do not give a password, the chart makes one. To show it, use `kubectl get secret opengovmail-thunder-admin -n opengovmail -o jsonpath='{.data.password}' \| base64 -d` |
| `global.ravenIdp.clientSecret` | The secret that Raven uses to sign in to ThunderID. Set it on the first installation only. Upgrades use the value that is in the cluster. |
| `global.ravenIdp.directAuthSecret` and `thunderid.configuration.server.security.directAuthSecret` | The secret that protects the password check of ThunderID. The two values must be the same. |

The user that runs Helm must have `get` permission on Secrets and `list` permission on Deployments in the namespace. The chart reads these items when it makes the Kubernetes resources. The `admin` and `edit` roles of Kubernetes have the two permissions.

### Step 7: Add the DKIM record

1. Make sure that the pods show `Running`.
2. Show the DKIM record:

   ```bash
   kubectl exec -n opengovmail statefulset/opengovmail-opendkim -c opendkim -- \
     cat /etc/dkimkeys/example.com/mail.txt
   ```

3. Add the record as a TXT record with the name `mail._domainkey`:

   | Type | Name | Value |
   |---|---|---|
   | TXT | mail._domainkey | "v=DKIM1; h=sha256; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQ..." |

> [!Important]
> Join the parts of the key into one value. The command shows the key in some parts, each in quotation marks. Remove the quotation marks and the spaces between the parts.

## Verify installation

1. Examine the pods and the certificates:

   ```bash
   kubectl get pods -n opengovmail
   kubectl get certificate -n opengovmail
   ```

2. Make sure that each pod shows `Running` or `Completed`.
3. Make sure that `mail-example-com-tls`, `example-com-tls` and `thunderid-tls` show `True` in the `READY` column. If a certificate does not become ready, use `kubectl describe certificate <name> -n opengovmail` to find the cause.
4. Open `https://example.com:8090/console`. Sign in as `admin`.
5. [Add a user](how-to/add-users.md). Then [connect a mail app](user-guide/mail-client.md).
6. Examine the server from outside your network:
   - [MXToolbox](https://mxtoolbox.com/SuperTool.aspx) examines your DNS records, blocklists and mail server.
   - [mail-tester](https://www.mail-tester.com/) gives a test email a score out of 10.

## Post installation

### Configuration options

- [Helm values](reference/helm-values.md) gives the settings that you will possibly change.
- [Add the collaboration services](how-to/enable-collaboration.md): Nextcloud, Collabora, Talk and Matrix.
- [Add users](how-to/add-users.md) in the ThunderID console.

### Other clusters

We test the chart on k3s with one node. On a different cluster, do these changes:

- In `infra/clusterissuers-http01.yaml` and in your values, change the ingress class from `traefik` to the class of your ingress controller.
- Postfix uses ports 25 and 587 as host ports. Raven uses ports 143 and 993 as host ports. Put them on one node. The MX record must point to the public IP address of this node.
- Raven connects to ThunderID through the public name of the server on port 8090. The cluster must send traffic from a pod to the external IP address of its node (hairpin NAT). k3s does this.

### Upgrade options

[Upgrade OpenGovMail](how-to/upgrade.md) gives the procedure to change to a new chart version. [CHANGELOG.md](../CHANGELOG.md) gives the changes in each version and the steps that an upgrade possibly needs.

### Downgrade options

Helm can return the release to an earlier revision:

```bash
helm history opengovmail -n opengovmail
helm rollback opengovmail <revision> -n opengovmail
```

A rollback restores the Kubernetes resources. It does not restore data. Before you return to an earlier release, read [CHANGELOG.md](../CHANGELOG.md). Look for releases that changed the stored data.

### Uninstallation options

> [!Warning]
> Back up the volumes before you uninstall OpenGovMail. The uninstallation deletes the volumes that the chart creates directly. These include the Raven volume, which has all the mailboxes. To show the volumes, use `kubectl get pvc -n opengovmail`.

1. Uninstall the chart:

   ```bash
   helm uninstall opengovmail -n opengovmail
   ```

   The volumes of StatefulSets stay, for example SeaweedFS (attachments) and OpenDKIM (signing keys).

2. To delete all the items that stay, delete the namespace:

   ```bash
   kubectl delete namespace opengovmail
   ```

The cert-manager installation and the issuers are for all of the cluster. If no other application uses them, remove them separately.

## Troubleshooting

[Troubleshooting](troubleshooting.md) gives the solutions to the usual problems. Some examples are:

- A certificate that does not become ready.
- A ThunderID setup job that does not complete.
- DKIM failures.
- Mail in spam folders.

## Next steps

- [Add users](how-to/add-users.md)
- [Add the collaboration services](how-to/enable-collaboration.md)
- Give the [user guide](user-guide/README.md) to your users

## Product version history

[CHANGELOG.md](../CHANGELOG.md) records each release. The versions follow [Semantic Versioning](https://semver.org/). The chart version and the OpenGovMail image tags change together.

## Definition of terms

Read the [glossary](glossary.md).
