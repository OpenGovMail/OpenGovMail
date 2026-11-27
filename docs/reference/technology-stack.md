# Technology stack

This page gives each component that the OpenGovMail Helm chart installs, and the platform for the chart. The versions are the versions in chart 0.6.1. The [overview of the architecture](../concepts/architecture.md) shows how the components connect. [CHANGELOG.md](../../CHANGELOG.md) gives the changes in each release.

## Platform

| Requirement | Version | Notes |
|---|---|---|
| CPU architecture | amd64 (x86_64) | The Postfix, OpenDKIM and Raven images are for amd64 only. |
| Kubernetes | k3s, current stable release | We test with one node. Other clusters must have an ingress controller and hairpin NAT. Read [Other clusters](../installation-guide.md#other-clusters). |
| Ingress controller | Traefik, from k3s | Set in `infra/clusterissuers-http01.yaml` and the Nextcloud ingress values. |
| Helm | 3 | 3.17 or later for `--take-ownership`. One [troubleshooting](../troubleshooting.md#the-thunderid-setup-job-does-not-complete-on-a-first-installation) procedure uses it. |
| cert-manager | 1.20.0 | Installed once per cluster. |

## Email and identity

These are on by default unless the table says otherwise.

| Component | Version | Role | Turned on with |
|---|---|---|---|
| [Postfix](https://www.postfix.org/) | image 0.6.1 | Sends and receives mail | `postfix.enabled` |
| [OpenDKIM](http://www.opendkim.org/) | image 0.6.1 | Signs outgoing mail | `opendkim.enabled` |
| [Raven](https://github.com/OpenGovMail/raven) | 0.9.1 | IMAP, LMTP delivery and SASL authentication | `raven.enabled` |
| [ThunderID](https://github.com/thunder-id/thunderid) | 1.0.2 | Users, groups, organization units and sign-in | `thunderid.enabled` |
| [SeaweedFS](https://github.com/seaweedfs/seaweedfs) | 4.51 | S3 object storage for attachments and files | `seaweedfs.enabled` |
| [Rspamd](https://rspamd.com/) | 4.1.1 | Spam filtering. Off by default. | `rspamd.enabled` |
| [Unbound](https://nlnetlabs.nl/projects/unbound/) | latest | DNS resolver for Rspamd's blocklist lookups. Off by default. | `unbound.enabled` |
| [Valkey](https://valkey.io/) | 9.1 | Rspamd's storage. Off by default. | `valkey.enabled` |

## Collaboration

By default, all these components are off. Read [Add the collaboration services](../how-to/enable-collaboration.md).

| Component | Version | Role |
|---|---|---|
| [Nextcloud](https://nextcloud.com/) | 34.0.4 (chart 9.3.0) | Calendars, files, webmail and Talk |
| [Collabora Online (CODE)](https://www.collaboraonline.com/code/) | 25.04.9.4.1 | Document editor |
| [coturn](https://github.com/coturn/coturn) | 4.18.0 | TURN server for Talk |
| [Tuwunel](https://github.com/matrix-construct/tuwunel) | 2.0.0 | Matrix homeserver |
| [LiveKit](https://livekit.io/) | 1.13.11 | Media server for Matrix calls |
| [lk-jwt-service](https://github.com/element-hq/lk-jwt-service) | 0.10.0 | Gives signed-in Matrix users a LiveKit token |
| [PostgreSQL](https://www.postgresql.org/) | 16 | Nextcloud's database |
| [Valkey](https://valkey.io/) | 9.1 | Nextcloud's cache |

## Where versions are set

| Component | File |
|---|---|
| Chart, Postfix and OpenDKIM images | [Chart.yaml](../../charts/opengovmail/Chart.yaml), and `image.tag` in each subchart's `values.yaml` |
| ThunderID and Nextcloud charts | `dependencies` in [Chart.yaml](../../charts/opengovmail/Chart.yaml) |
| Raven, SeaweedFS, Rspamd, Valkey, Unbound, Tuwunel | `image.tag` or `appVersion` in [charts/opengovmail/charts/](../../charts/opengovmail/charts) |
| Collabora, coturn, PostgreSQL, LiveKit, lk-jwt-service | [values.yaml](../../charts/opengovmail/values.yaml) |
