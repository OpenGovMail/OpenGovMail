# Technology stack

This page lists every component the OpenGovMail Helm chart installs, with the version chart 0.6.1 ships, and the platform it runs on. [Overview of the architecture](../concepts/architecture.md) explains how the components connect. [CHANGELOG.md](../../CHANGELOG.md) lists what changed in each release.

## Platform

| Requirement | Version | Notes |
|---|---|---|
| CPU architecture | amd64 (x86_64) | The Postfix, OpenDKIM and Raven images are built for amd64 only. |
| Kubernetes | k3s, current stable release | The tested setup is a single node. Other clusters need an ingress controller and hairpin NAT; see [Other clusters](../installation-guide.md#other-clusters). |
| Ingress controller | Traefik, as shipped with k3s | Set in `infra/clusterissuers-http01.yaml` and the Nextcloud ingress values. |
| Helm | 3 | 3.17 or later for `--take-ownership`, used in one [troubleshooting](../troubleshooting.md#thunderid-setup-job-never-finishes-on-a-first-install) fix. |
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
| [Redis](https://redis.io/) | 7 | Rspamd's storage. Off by default. | `redis.enabled` |
| [Unbound](https://nlnetlabs.nl/projects/unbound/) | latest | DNS resolver for Rspamd's blocklist lookups. Off by default. | `unbound.enabled` |

## Collaboration

All off by default. See [Turn on collaboration](../how-to/enable-collaboration.md).

| Component | Version | Role |
|---|---|---|
| [Nextcloud](https://nextcloud.com/) | 34.0.4 (chart 9.3.0) | Calendars, files, webmail and Talk |
| [Collabora Online (CODE)](https://www.collaboraonline.com/code/) | 25.04.9.4.1 | Document editor |
| [coturn](https://github.com/coturn/coturn) | 4.18.0 | TURN server for Talk |
| [Tuwunel](https://github.com/matrix-construct/tuwunel) | 2.0.0 | Matrix homeserver |
| [LiveKit](https://livekit.io/) | 1.13.11 | Media server for Matrix calls |
| [lk-jwt-service](https://github.com/element-hq/lk-jwt-service) | 0.10.0 | Gives signed-in Matrix users a LiveKit token |
| [PostgreSQL](https://www.postgresql.org/) | 16 | Nextcloud's database |
| [Redis](https://redis.io/) | 7 | Nextcloud's cache |

## Where versions are set

| Component | File |
|---|---|
| Chart, Postfix and OpenDKIM images | [Chart.yaml](../../charts/opengovmail/Chart.yaml), and `image.tag` in each subchart's `values.yaml` |
| ThunderID and Nextcloud charts | `dependencies` in [Chart.yaml](../../charts/opengovmail/Chart.yaml) |
| Raven, SeaweedFS, Rspamd, Redis, Unbound, Tuwunel | `image.tag` or `appVersion` in [charts/opengovmail/charts/](../../charts/opengovmail/charts) |
| Collabora, coturn, PostgreSQL, LiveKit, lk-jwt-service | [values.yaml](../../charts/opengovmail/values.yaml) |
