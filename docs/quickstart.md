# Quickstart

## Overview

This quickstart installs OpenGovMail on a new Linux server. At the end, you send your first email.

The quickstart installs only email and identity. It uses a k3s cluster with one node and the most common settings. For collaboration, other clusters and all the options, read the [installation guide](installation-guide.md).

The procedure takes approximately one hour. Most of this time is for the DNS records to change.

In the steps that follow, replace `example.com` with your domain. Replace `12.34.56.78` with the IP address of your server.

## Before you start

You must have these items:

- A Linux server on amd64 (x86_64). The published images are for amd64 only. They do not operate on ARM servers.
- A minimum of 4 GB of memory on the server.
- A static public IP address and root access on the server.
- A domain. You must have access to its DNS records.
- Access to the reverse DNS (PTR) settings of your server at your hosting provider.
- Ports 25, 80, 443, 587, 993 and 8090 open to the internet. Some cloud providers close port 25. If your provider closes port 25, ask the provider to open it.

## Part 1: Prepare the server

### Step 1: Add DNS records

1. Add these records in the control panel of your DNS provider:

   | Type | Name | Value |
   |---|---|---|
   | A | mail | 12.34.56.78 |
   | A | example.com | 12.34.56.78 |
   | MX | example.com | mail.example.com |
   | TXT | example.com | "v=spf1 ip4:12.34.56.78 ~all" |
   | TXT | _dmarc | "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com" |

2. Ask your hosting provider to set the PTR record for `12.34.56.78` to `mail.example.com`.
3. Make sure that the A records show the IP address of your server:

   ```bash
   dig +short mail.example.com
   dig +short example.com
   ```

   The two commands must show the IP address of your server. If they do not, wait and do this step again. Let's Encrypt cannot give certificates for a hostname that does not show the IP address of the server.

### Step 2: Install k3s and Helm

1. On the server, run these commands:

   ```bash
   curl -sfL https://get.k3s.io | sh -
   mkdir -p ~/.kube
   sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
   sudo chown "$USER" ~/.kube/config
   curl -fsSL https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
   ```

2. Make sure that the cluster operates:

   ```bash
   kubectl get nodes
   ```

   The command must show one node with the status `Ready`.

### Step 3: Get the chart

```bash
git clone https://github.com/OpenGovMail/OpenGovMail.git
cd OpenGovMail
```

## Part 2: Prepare the certificates

### Step 1: Install cert-manager

```bash
helm install cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.20.0 \
  --namespace cert-manager \
  --create-namespace \
  --set crds.enabled=true
```

### Step 2: Create the certificate issuers

1. Create the issuers:

   ```bash
   kubectl apply -f infra/clusterissuers-http01.yaml
   ```

2. Make sure that the issuers are ready:

   ```bash
   kubectl get clusterissuer
   ```

   Wait until `le-staging` and `le-prod` show `True` in the `READY` column.

## Part 3: Install OpenGovMail

### Step 1: Write a values file

1. Create the file `my-values.yaml` in your home directory. Do not put it in the repository. This prevents a commit of the file by accident.
2. Copy this text into the file:

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
         name: thunderid-tls
         items:
           - key: tls.crt
             mountPath: /opt/thunderid/config/tls/server.cert
           - key: tls.key
             mountPath: /opt/thunderid/config/tls/server.key
   ```

3. Replace `example.com` with your domain. Do not change the Secret name `thunderid-tls`. With this name, the chart gets the ThunderID certificate before the ThunderID setup job starts.

### Step 2: Install the chart

1. Choose a password for the ThunderID administrator.
2. Run these commands. Replace `<your-admin-password>` with the password:

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

### Step 3: Make sure that OpenGovMail operates

1. Examine the pods and the certificates:

   ```bash
   kubectl get pods -n opengovmail
   kubectl get certificate -n opengovmail
   ```

2. Make sure that each pod shows `Running` or `Completed`, and that each certificate shows `True` in the `READY` column. This can take some minutes.

## Part 4: Complete the DNS records

### Step 1: Add the DKIM record

OpenDKIM makes a DKIM key when it starts for the first time.

1. Show the DKIM record:

   ```bash
   kubectl exec -n opengovmail statefulset/opengovmail-opendkim -c opendkim -- \
     cat /etc/dkimkeys/example.com/mail.txt
   ```

2. Add the record as a TXT record with the name `mail._domainkey`.

The command shows the key in some parts, each in quotation marks. Join the parts into one value. Remove the quotation marks and the spaces between the parts.

## Part 5: Send your first email

### Step 1: Create a user

1. Open `https://example.com:8090/console`.
2. Sign in as `admin` with the password that you chose.
3. Add a user to the `default` organization unit. Use the `Person` user type. Give a username and a password.

For more data, read [Add users](how-to/add-users.md).

### Step 2: Connect a mail app

Add the account to Thunderbird or to a different mail app. Use these settings. Sign in with the full email address and the password of the user.

| Setting | Value |
|---|---|
| Incoming (IMAP) | `mail.example.com`, port 993, SSL/TLS |
| Outgoing (SMTP) | `mail.example.com`, port 587, STARTTLS |

### Step 3: Do a test of the delivery

1. Go to [mail-tester.com](https://www.mail-tester.com/).
2. Send an email to the address that it shows.
3. Read the result. Mail-tester gives your server a score out of 10. It also shows each DNS record or header that is not correct.

## Next steps

- [Add more users](how-to/add-users.md). Then give them the [user guide](user-guide/README.md).
- [Add the collaboration services](how-to/enable-collaboration.md): calendars, files, documents, meetings and chat.
- Read the [installation guide](installation-guide.md) for upgrades, other settings and removal.
- Read the [overview of the architecture](concepts/architecture.md). It shows how the parts of OpenGovMail connect.
