# Quickstart

## Overview

This quickstart takes a fresh Linux server to a working OpenGovMail email server and your first sent message. It installs email and identity only, on a single-node k3s cluster, with the settings most installs use. Collaboration, other clusters and every option are in the [installation guide](installation-guide.md).

Expect it to take about an hour, most of it waiting for DNS to update.

In the steps below, replace `example.com` with your domain and `12.34.56.78` with your server's IP address.

## Before you start

You need:

- A Linux server on amd64 (x86_64) with at least 4 GB of memory, a static public IP address and root access. The published images are built for amd64 only, so ARM servers don't work.
- A domain whose DNS records you can edit.
- Access to your hosting provider's settings for the server's reverse DNS (PTR) record.
- Ports 25, 80, 443, 587, 993 and 8090 open to the internet. Some cloud providers block port 25 until you ask them to open it.

## Part 1: Prepare the server

### Step 1: Add DNS records

Add these records in your DNS provider's control panel:

| Type | Name | Value |
|---|---|---|
| A | mail | 12.34.56.78 |
| A | example.com | 12.34.56.78 |
| MX | example.com | mail.example.com |
| TXT | example.com | "v=spf1 ip4:12.34.56.78 ~all" |
| TXT | _dmarc | "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com" |

Ask your hosting provider to set the PTR record for `12.34.56.78` to `mail.example.com`.

Check that the A records answer before you go on. Let's Encrypt can't issue certificates until they do:

```bash
dig +short mail.example.com
dig +short example.com
```

Both commands should print your server's IP address.

### Step 2: Install k3s and Helm

On the server:

```bash
curl -sfL https://get.k3s.io | sh -
mkdir -p ~/.kube
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown "$USER" ~/.kube/config
curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
kubectl get nodes
```

The last command should show one node with status `Ready`.

### Step 3: Get the chart

```bash
git clone https://github.com/OpenGovMail/OpenGovMail.git
cd OpenGovMail
```

## Part 2: Set up certificates

### Step 1: Install cert-manager

```bash
helm install cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.20.0 \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
```

### Step 2: Create the certificate issuers

```bash
kubectl apply -f infra/clusterissuers-http01.yaml
kubectl get clusterissuer
```

Wait until `le-staging` and `le-prod` both show `READY` as `True`.

## Part 3: Install OpenGovMail

### Step 1: Write a values file

Create `my-values.yaml` outside the repository, so it never gets committed. Replace `example.com`, and `example-com` in the Secret name:

```yaml
global:
  domain: example.com
  tls:
    issuer: le-prod

thunderid:
  configuration:
    server:
      publicUrl: https://example.com:8090
    gateClient:
      hostname: example.com
    tls:
      certFile: config/tls/server.cert
      keyFile: config/tls/server.key
  declarativeResources:
    enabled: true
    readOnly: true
    secret:
      name: example-com-tls
      items:
        - key: tls.crt
          mountPath: /opt/thunderid/config/tls/server.cert
        - key: tls.key
          mountPath: /opt/thunderid/config/tls/server.key
```

### Step 2: Install the chart

Pick an administrator password for ThunderID, then run:

```bash
helm dependency update ./charts/opengovmail
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

> [!Note]
> On chart 0.6.1 this command can wait forever on ThunderID's setup job. [Troubleshooting](troubleshooting.md#thunderid-setup-job-never-finishes-on-a-first-install) has the fix.

### Step 3: Check that it's running

```bash
kubectl get pods -n opengovmail
kubectl get certificate -n opengovmail
```

Every pod should be `Running` or `Completed`, and every certificate `READY`. This can take a few minutes.

## Part 4: Finish DNS

### Step 1: Add the DKIM record

Print the DKIM key that OpenDKIM generated on its first start:

```bash
kubectl exec -n opengovmail statefulset/opengovmail-opendkim -c opendkim -- \
  cat /etc/dkimkeys/example.com/mail.txt
```

Add it as a TXT record named `mail._domainkey`. The key is printed in several quoted parts. Join them into one value, without the quotes and the spaces between them.

## Part 5: Send your first email

### Step 1: Create a user

1. Open `https://example.com:8090/console` and sign in as `admin` with the password you chose.
2. Add a user to the `default` organization unit, with the `Person` user type, a username and a password.

[Add users](how-to/add-users.md) has the details.

### Step 2: Connect a mail client

Set up Thunderbird or another client with these settings, signing in with the user's full email address and password:

| Setting | Value |
|---|---|
| Incoming (IMAP) | `mail.example.com`, port 993, SSL/TLS |
| Outgoing (SMTP) | `mail.example.com`, port 587, STARTTLS |

### Step 3: Test delivery

Go to [mail-tester.com](https://www.mail-tester.com/), and send an email to the address it shows. It scores your setup out of 10 and lists any DNS record or header it found wrong.

## Next steps

- [Add more users](how-to/add-users.md) and [connect their apps](user-guide/README.md).
- [Turn on collaboration](how-to/enable-collaboration.md): calendars, files, documents, meetings and chat.
- Read the [installation guide](installation-guide.md) for upgrades, other settings and removal.
- Read the [overview of the architecture](concepts/architecture.md) to see how the parts fit together.
