# OpenGovMail
**_Modern Collaborative Email Platform_**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![CI](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/build-and-push-images.yml)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/build-and-push-images.yml)
[![Security Scan](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/trivy-fs.yml?label=security)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/trivy-fs.yml)
[![Last Commit](https://img.shields.io/github/last-commit/OpenGovMail/OpenGovMail)](https://github.com/OpenGovMail/OpenGovMail/commits/main)
[![Slack](https://img.shields.io/badge/Slack-join%20chat-4A154B?logo=slack&logoColor=white)](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA)

<p align="center">
  <a href="#project-description">Description</a> •
  <a href="#who-this-project-is-for">Who it is for</a> •
  <a href="#key-features">Key features</a> •
  <a href="#project-dependencies">Dependencies</a> •
  <a href="#instructions-for-using-opengovmail">Using OpenGovMail</a> •
  <a href="#contributing-guidelines">Contributing</a> •
  <a href="#additional-documentation">Documentation</a> •
  <a href="#how-to-get-help">Help</a> •
  <a href="#terms-of-use">Terms of use</a>
</p>

## Project description

OpenGovMail is an email and communication system for organizations as large as a government. Its aim is email that is fast and easy to manage, and also safe and reliable. The platform has two stages. Version 1.0 supplies reliable email that obeys the standards. Version 2.0 adds collaboration services: calendars, shared drives, document editing, meetings and chat.

OpenGovMail operates on your own Kubernetes cluster from one Helm chart. Thus, your mail, files and accounts stay on your servers. Each user has one identity: the email address and the account are the same. You can use your own ThunderID or the ThunderID that comes with OpenGovMail. ThunderID keeps your users and the structure of your organization.

The [overview of the architecture](docs/concepts/architecture.md) shows how the parts connect.

## Who this project is for

- Governments, public bodies and other organizations that want to operate their own email.
- System administrators who install and operate OpenGovMail. They must know Linux, DNS and the basics of Kubernetes.
- Developers who change or extend OpenGovMail. Start with [CONTRIBUTING.md](CONTRIBUTING.md).

## Key features

| Feature | Status |
|---------|--------|
| No licence fees. You pay only for the servers it runs on. | Completed |
| Runs on modest hardware: 4 GB of memory for email. | Completed |
| Mail, attachments and accounts stay on your own servers. | Completed |
| Attachments are kept in object storage, apart from the mailboxes. | Completed |
| Group addresses deliver to every member of a group. | Completed |
| Role addresses, such as `director@`, belong to the organization rather than a person. | Completed |
| Accounts come from ThunderID, either the bundled one or one you already run. | Completed |
| Spam filtering with Rspamd. You add it in the Helm values. | Completed |
| Calendars, shared drives, document editing, video meetings and chat. You add each one separately. | Completed |
| Malware scanning. | Planned |
| Configurable backups. | Planned |
| Email workflows that route a message through set steps, with tracking. | Planned |
| A native email client with IMAP and JMAP support. | Planned |
| Expiry dates and retention policies for email. | Planned |

**[View Complete Feature List →](docs/FEATURE_LIST.md)**

## Project dependencies

To operate OpenGovMail, you must have:

- A Linux server on amd64 (x86_64) with a static public IP address and root access. The published images are for amd64 only.
- A minimum of 4 GB of memory. With the collaboration services, a minimum of 2 CPUs.
- A domain with DNS records that you control, and a PTR record for the IP address of the server.
- A Kubernetes cluster. We test OpenGovMail on [k3s](https://k3s.io/) with one node.
- [Helm](https://helm.sh/docs/intro/install/), [kubectl](https://kubernetes.io/docs/tasks/tools/) and [Git](https://git-scm.com/downloads/linux).

The [technology stack](docs/reference/technology-stack.md) gives each component and its version.

## Instructions for using OpenGovMail

### Install

- [Quickstart](docs/quickstart.md): from a new server to your first email.
- [Installation guide](docs/installation-guide.md): all the steps and options. It also tells you how to upgrade and remove OpenGovMail.

### Configure

- [Helm values](docs/reference/helm-values.md): the settings that you will possibly change.
- [Add users](docs/how-to/add-users.md) in the ThunderID console.
- [Add the collaboration services](docs/how-to/enable-collaboration.md): Nextcloud, Collabora, Talk and Matrix.

### Operate

- [Upgrade OpenGovMail](docs/how-to/upgrade.md) to a new release.
- [User guide](docs/user-guide/README.md): connect mail, calendar and files apps.

### Troubleshoot

[Troubleshooting](docs/troubleshooting.md) gives solutions for certificates that do not become ready, DKIM failures, mail in spam folders and pods that do not start.

## Contributing guidelines

We welcome all contributions, large and small. [CONTRIBUTING.md](CONTRIBUTING.md) tells you how to prepare a development environment, keep a fork up to date and open a pull request. All contributors obey the [code of conduct](CODE_OF_CONDUCT.md).

To report a security problem, do not open an issue. Use the procedure in [SECURITY.md](SECURITY.md).

## Additional documentation

- [Documentation index](docs/README.md)
- [Release notes](CHANGELOG.md). There is one entry for each version. The versions follow [Semantic Versioning](https://semver.org/).
- [Glossary](docs/glossary.md)

## How to get help

- Ask in the [OpenGovMail Slack](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA).
- To report a bug or to ask for a feature, open a [GitHub issue](https://github.com/OpenGovMail/OpenGovMail/issues).

## Terms of use

OpenGovMail has the Apache 2.0 License. Read [LICENSE](LICENSE).
