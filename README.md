# OpenGovMail
**_Modern Collaborative Email Platform_**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![CI](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/build-and-push-images.yml)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/build-and-push-images.yml)
[![Security Scan](https://img.shields.io/github/actions/workflow/status/OpenGovMail/OpenGovMail/trivy-fs.yml?label=security)](https://github.com/OpenGovMail/OpenGovMail/actions/workflows/trivy-fs.yml)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/16893/badge)](https://www.bestpractices.dev/projects/16893)
[![Last Commit](https://img.shields.io/github/last-commit/OpenGovMail/OpenGovMail)](https://github.com/OpenGovMail/OpenGovMail/commits/main)
[![Slack](https://img.shields.io/badge/Slack-join%20chat-4A154B?logo=slack&logoColor=white)](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA)

<p align="center">
  <a href="#project-description">Description</a> •
  <a href="#why-opengovmail">Why OpenGovMail?</a> •
  <a href="#who-this-project-is-for">Who it's for</a> •
  <a href="#project-dependencies">Dependencies</a> •
  <a href="#instructions-for-using-opengovmail">Using OpenGovMail</a> •
  <a href="#contributing-guidelines">Contributing</a> •
  <a href="#additional-documentation">Documentation</a> •
  <a href="#how-to-get-help">Help</a> •
  <a href="#terms-of-use">Terms of use</a>
</p>

## Project description

OpenGovMail is an email and communication system built to work at government scale. It aims to make email faster, smarter and easier to manage while keeping it secure and reliable. The platform grows in two stages: version 1.0 delivers reliable, standards-compliant email, and version 2.0 adds collaboration around it: calendars, shared drives, document editing, meetings and chat.

OpenGovMail runs on your own Kubernetes cluster from one Helm chart, so mail, files and accounts stay on your servers. Each user has one identity, and their email address and their account are the same thing. You can bring your own ThunderID or use the bundled one to organize your users and map your organization's hierarchy.

[Overview of the architecture](docs/concepts/architecture.md) shows how the parts fit together.

## Why OpenGovMail?

OpenGovMail uses the standard email architecture, so it works with existing mail apps and any other mail server. Unlike hosted services:

- Mail, Matrix chat, storage and accounts run in one installation on servers you control, so data stays with you by design, not by contract.
- Every account lives in [ThunderID](https://github.com/thunder-id/thunderid), which checks every sign-in and delivery. One account covers mail, chat, files and calendars, so one change applies everywhere when someone joins, moves or leaves.
- Mail follows the organization units, groups and roles in ThunderID: a unit can have its own mail domain, a group address reaches its current members, and each role has its own address.
- A role address such as `director@` needs no shared password. Each holder opens it with their own sign-in, in mail apps that sign in through ThunderID, so it passes to a successor or stand-in with the role.
- Each part of a message, such as an attachment, is stored once however many people receive it, and each mailbox has its own database.
- Connections use TLS, mail is authenticated with SPF, DKIM and DMARC, and Rspamd filters spam and malware through a ClamAV server you run.
- It runs as standard containers from one Helm chart, built on Postfix, Rspamd, SeaweedFS and ThunderID plus Raven, the mail server from the Lanka Software Foundation.

The [complete feature list](docs/FEATURE_LIST.md) shows which features are done and which are planned.

## Who this project is for

- Governments, public bodies and other organizations that want to run their own email instead of renting it.
- System administrators who install and run it. You should be comfortable with Linux, DNS and basic Kubernetes.
- Developers who want to change or extend it. Start with [CONTRIBUTING.md](CONTRIBUTING.md).

## Project dependencies

To run OpenGovMail you need:

- A Linux server on amd64 (x86_64) with a static public IP address and root access. The published images are built for amd64 only.
- At least 4 GB of memory, and 2 CPUs if you turn on collaboration.
- A domain whose DNS you control, and a PTR record for the server's IP address.
- A Kubernetes cluster. A single-node [k3s](https://k3s.io/) is the tested setup.
- [Helm](https://helm.sh/docs/intro/install/), [kubectl](https://kubernetes.io/docs/tasks/tools/) and [Git](https://git-scm.com/downloads/linux).

[Technology stack](docs/reference/technology-stack.md) lists every component and its version.

## Instructions for using OpenGovMail

### Install

- [Quickstart](docs/quickstart.md): from a fresh server to your first sent email.
- [Installation guide](docs/installation-guide.md): every step and option, including upgrades and removal.

### Configure

- [Helm values](docs/reference/helm-values.md): the settings you are most likely to change.
- [Add users](docs/how-to/add-users.md) in the ThunderID console.
- [Add the collaboration services](docs/how-to/enable-collaboration.md): Nextcloud, Collabora, Talk and Matrix.

### Run

- [Upgrade OpenGovMail](docs/how-to/upgrade.md) to a new release.
- [User guide](docs/user-guide/README.md): connect mail, calendar and file apps.

### Troubleshoot

[Troubleshooting](docs/troubleshooting.md) covers certificates that aren't issued, DKIM failures, mail marked as spam and pods that don't start.

## Contributing guidelines

Contributions of any size are welcome. [CONTRIBUTING.md](CONTRIBUTING.md) explains how to set up a development environment, keep a fork up to date and open a pull request. Everyone taking part follows the [code of conduct](CODE_OF_CONDUCT.md).

To report a security problem, follow [SECURITY.md](SECURITY.md) instead of opening an issue.

## Additional documentation

- [Documentation index](docs/README.md)
- [Release notes](CHANGELOG.md), one entry per version, following [Semantic Versioning](https://semver.org/)
- [Glossary](docs/glossary.md)

## How to get help

- Ask in the [OpenGovMail Slack](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA).
- Open a [GitHub issue](https://github.com/OpenGovMail/OpenGovMail/issues) for a bug or a feature request.

## Terms of use

OpenGovMail is distributed under the Apache 2.0 License. See [LICENSE](LICENSE).
