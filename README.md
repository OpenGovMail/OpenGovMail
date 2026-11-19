# OpenGovMail
**_Modern Collaborative Email Platform_**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![CI](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/build-and-push-images.yml)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/build-and-push-images.yml)
[![Security Scan](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/trivy-fs.yml?label=security)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/trivy-fs.yml)
[![Last Commit](https://img.shields.io/github/last-commit/OpenGovMail/OpenGovMail)](https://github.com/OpenGovMail/OpenGovMail/commits/main)
[![Slack](https://img.shields.io/badge/Slack-join%20chat-4A154B?logo=slack&logoColor=white)](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA)

**OpenGovMail** aims to build a new kind of email and communication system that can work at a government scale. The goal is to make email faster, smarter, and easier to manage while keeping it secure and reliable. The platform will evolve in two stages: Version 1.0 delivers reliable, standards-compliant email, while Version 2.0 adds collaboration around it: calendars, shared drives, document editing, meetings and chat.

<p align="center">
  •   <a href="#why-OpenGovMail">Why OpenGovMail?</a> •
  <a href="#key-features">Key Features</a> •
  <a href="#getting-started">Getting Started</a> •
  <a href="#collaboration-optional">Collaboration</a> •
  <a href="#open-source-components">Open Source Components</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a> •
</p>

## Why OpenGovMail?
OpenGovMail is designed to be secure, reliable, and easy to manage. It runs entirely on your own hardware, giving you full control over your data and ensuring privacy. The system is lightweight and efficient, performing well even on minimal hardware, which makes it easy to deploy in a variety of environments. Each user has one identity, so their email address and their account are the same thing. You can bring your own ThunderID or use the bundled one to organize your users and map your organization hierarchically.

External firewalls are not required to filter emails, and attachments are stored separately in blob storage to save space and improve overall system performance. The Docker Compose setup also includes observability (Prometheus, Grafana and Loki), so administrators can monitor activity and catch issues early.

## Key Features

| Feature | Status |
|---------|--------|
| No licence fees. You pay only for the servers it runs on. | Completed |
| Runs on modest hardware: 4 GB of memory for email. | Completed |
| Mail, attachments and accounts stay on your own servers. | Completed |
| Spam filtering and malware scanning are built in (on Kubernetes, spam filtering is optional and malware scanning isn't included yet). | Completed |
| Attachments are kept in object storage, apart from the mailboxes. | Completed |
| Group addresses deliver to every member of a group. | Completed |
| Role addresses, such as `director@`, belong to the organization rather than a person. | Completed |
| Accounts come from ThunderID, either the bundled one or one you already run. | Completed |
| Configurable backups. | Planned |
| Email workflows that route a message through set steps, with tracking. | Planned |
| A native email client with IMAP and JMAP support. | Planned |
| Expiry dates and retention policies for email. | Planned |

**[View Complete Feature List →](docs/FEATURE_LIST.md)**

## Getting Started

There are two ways to install OpenGovMail:

- Kubernetes with Helm is the main path. It includes the optional [collaboration](#collaboration-optional) services.
- Docker Compose runs email only, on a single server.

Release notes for each version are in [CHANGELOG.md](CHANGELOG.md).

### Prerequisites
- A dedicated Linux server with a static public IP address. You also require root access and port access control.
- Domain with DNS control
- At least 4GB of memory, and 2 CPUs if you turn on [collaboration](#collaboration-optional)
- [Git](https://git-scm.com/downloads/linux), plus either a Kubernetes cluster with [Helm](https://helm.sh/docs/intro/install/) (a single-node [k3s](https://k3s.io/) works) or [Docker Engine](https://docs.docker.com/engine/install/)

### DNS setup
You own <a>example.com</a> and want to send an email as person@example.com.

You will need to add a few records to your DNS control panel.

> [!Note]
> Replace example.com and 12.34.56.78 in the below example with your domain and ip address.

| DNS Record | Name        | Value                                                  |
| ---------- | ----------- | ------------------------------------------------------ |
| A          | mail        | 12.34.56.78                                            |
| A          | example.com | 12.34.56.78                                            |
| MX         | example.com | mail.example.com                                       |
| TXT        | example.com | "v=spf1 ip4:12.34.56.78 ~all"                          |
| TXT        | _dmarc      | "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com" |
| PTR        | 12.34.56.78 | mail.example.com                                       |

> [!Tip]
> PTR records are usually set through your hosting provider. 

The DKIM record is added after the install, because the key is generated on first start.

### Install on Kubernetes

```bash
git clone https://github.com/OpenGovMail/OpenGovMail.git
cd OpenGovMail
```

Then follow [charts/README.md](charts/README.md). It covers cert-manager, the certificate issuers, the first install and upgrades. Start your values file from [values-example.yaml](charts/opengovmail/values-example.yaml).

When the pods are running, print the DKIM record:

```bash
kubectl exec -n opengovmail statefulset/opengovmail-opendkim -c opendkim -- cat /etc/dkimkeys/example.com/mail.txt
```

### Install with Docker Compose (email only)

```bash
git clone https://github.com/OpenGovMail/OpenGovMail.git
cd OpenGovMail
```

- Open [`opengovmail.yaml`](https://github.com/OpenGovMail/OpenGovMail/blob/main/conf/opengovmail.yaml) with a text editor.

- Enter your domain name.

- Create your own environment file from the template and fill in the values before running setup:

```bash
cp services/.env.example services/.env
nano services/.env
```

- Follow the manual setup instructions in the [SeaweedFS README](https://github.com/OpenGovMail/OpenGovMail/blob/main/services/seaweedfs/README.md) to create the `.env` and `s3-config.json` files.

- Run `bash scripts/setup/setup.sh` to set up the configs.

- Run `bash scripts/service/start-opengovmail.sh` to start the mail server. It prints the DKIM record at the end.

### DKIM record

Add the record printed by either install. It looks like this:

| DNS Record | Name            | Value                                                                                                                                                                                                                                                  |
| ---------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| TXT        | mail._domainkey | "v=DKIM1; h=sha256; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDYZd3CAas0+81zf13cvtO6o0+rlGx8ZobYQXRR9W8qcJOeO1SiQGx8F4/DjZE1ggujOaY1bkt8OnUg7vG7/bk5PNe05EHJrg344krodqCJrVI74ZzEB77Z1As395KX6/XqbQxBepQ8D5+RpGFOHitI443G/ZWgZ6BRyaaE6t3u0QIDAQAB" |

> [!Important] 
> Ensure that your dkim value is correctly formatted. The key is printed in several quoted parts. Join them into one value, without the quotes and spaces between them.

### Adding users

Users in OpenGovMail are provisioned through the [ThunderID](https://github.com/thunder-id/thunderid) console that ships with the platform. Sign in to the console at `https://<your-domain>:8090/console` with your admin credentials.

For the full walkthrough, including creating a user schema, assigning it to an organization unit, and adding users by password or invitation link, see [Adding Users](docs/Adding-Users.md).

### Testing your setup
- Now that you have a working email server, you can test your configuration using the following links/scripts.

  - [mxtoolbox](https://mxtoolbox.com/SuperTool.aspx)
    - Checks your DNS records, blocklists and mail server.
  - [mail-tester](https://www.mail-tester.com/)
    - Send it an email and it scores your server setup out of 10 for deliverability.
  

- You can also set up a Mail User Agent (MUA) like Thunderbird to send and receive emails. Follow the instructions in [Mail User Agent Setup](docs/Mail-User-Agent-Setup.md).

## Collaboration (optional)

The Helm chart can also run calendars, shared drives, document editing, video meetings and chat for the same users, who sign in with their mail account. All of it is off by default. [docs/Collaboration.md](docs/Collaboration.md) explains how to turn each part on and which hostnames and ports it needs.

## Open Source Components

### Email Flow in OpenGovMail

<table>
<tr>
<td width="50%">

**Figure 1: Outbound Email Flow**

<img src="https://github.com/user-attachments/assets/17938298-8f0c-4558-9331-6e96272cb0ba" alt="Inbound Email Flow" width="100%" />

</td>
<td width="50%">

**Figure 2: Inbound Email Flow**

<img src="https://github.com/user-attachments/assets/58e81d33-6a26-4bff-92fb-663737446740" alt="Outbound Email Flow" width="100%" />

</td>
</tr>
</table>

OpenGovMail is built using open-source software. [docs/Architecture.md](docs/Architecture.md) shows how the parts connect and lists the version of each.

- [Postfix](https://www.postfix.org/) - Handles sending and receiving email.
- [Raven](https://github.com/OpenGovMail/raven) - Handles SASL authentication, LMTP, and IMAP server for email retrieval.
- [ThunderID](https://github.com/thunder-id/thunderid) - Identity provider and user manager
- [OpenDKIM](http://www.opendkim.org/) - Signs outgoing mail.
- [Rspamd](https://rspamd.com/) - Spam filtering system.
- [ClamAV](https://docs.clamav.net/Introduction.html) - Virus scanning system.
- [SeaweedFS](https://github.com/seaweedfs/seaweedfs) - Object storage for attachments and files.
- [cert-manager](https://cert-manager.io/) - Gets Let's Encrypt certificates on Kubernetes.
- [Certbot](https://certbot.eff.org/) - Gets Let’s Encrypt certificates for Docker Compose.

## Contributing

Thank you for wanting to contribute to our project. Please see [CONTRIBUTING.md](https://github.com/OpenGovMail/OpenGovMail/blob/main/docs/CONTRIBUTING.md) for more details. To report a security problem, follow [SECURITY.md](docs/SECURITY.md) instead of opening an issue.

## License 

Distributed under the Apache 2.0 License. See [LICENSE](https://github.com/OpenGovMail/OpenGovMail/blob/main/LICENSE) for more information.
