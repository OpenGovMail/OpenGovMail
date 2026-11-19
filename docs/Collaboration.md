# Collaboration

The Helm chart can run collaboration services next to email. Users sign in to all of them through ThunderID with their mail account. Everything here is off by default, and Docker Compose doesn't include any of it.

| Part | What it gives users | Turn it on with |
|---|---|---|
| Nextcloud | Calendars (CalDAV), shared drives, files and webmail | `nextcloud.enabled: true` |
| Collabora | Document editing inside Nextcloud | On with Nextcloud (`nextcloud.collabora.enabled`) |
| Nextcloud Talk | Video meetings, with a TURN server so calls work behind NAT | On with Nextcloud (`nextcloudDeps.turn.enabled`) |
| Matrix (Tuwunel) | Chat in Element or any Matrix client | `tuwunel.enabled: true` |
| Matrix calls | Voice and video calls in Element | `tuwunel.matrixRtc.enabled: true` |

[values-example.yaml](../charts/opengovmail/values-example.yaml) turns all of them on. Copy the `nextcloud` and `tuwunel` blocks and change the hostnames.

With everything on, the server needs at least 2 CPUs. On 1 CPU the pods' CPU requests use up the node, and some pods, including cert-manager's, never get scheduled.

## Hostnames

Each service has its own hostname, pointed at the server like `mail.<domain>`. cert-manager gets the certificates.

| Default name | Set with | Serves |
|---|---|---|
| none, must be set | `nextcloud.nextcloud.host` and `nextcloud.ingress.tls` | Nextcloud, and Talk's TURN server |
| none, must be set | `nextcloud.collabora.ingress`, `collabora.server_name` and `collabora.aliasgroups` | Collabora. Users never open it directly: Nextcloud loads the editor from it |
| `matrix.<domain>` | `tuwunel.host` | Matrix |
| `matrix-rtc.<domain>` | `tuwunel.matrixRtc.host` | Matrix calls (LiveKit and its token service) |

Matrix also needs `https://<domain>/.well-known/matrix/client`, which tells Element where the Matrix server is. Tuwunel serves it on `<domain>`. If `<domain>` is hosted somewhere else, such as a website, see [DNS in the Tuwunel chart](../charts/opengovmail/charts/tuwunel/README.md#dns).

## Ports

Media for meetings and calls doesn't go through the ingress. Open these on the node's firewall:

| Port | For |
|---|---|
| 3478/udp and 3478/tcp | Talk's TURN server |
| 49160-49200/udp | Talk's TURN relay (`nextcloudDeps.turn.minPort` and `maxPort`) |
| 7882/udp | Matrix call media |
| 7881/tcp | Matrix call media on networks that block UDP |

If the node only has a private IP, set `nextcloudDeps.turn.externalIP` and `tuwunel.matrixRtc.livekit.nodeIP` to its public address. The [Tuwunel chart's calls section](../charts/opengovmail/charts/tuwunel/README.md#calls) has the details.

## User types

Nextcloud and Matrix register themselves as applications in ThunderID on every install and upgrade. Only the user types they list may sign in, and ThunderID rejects Nextcloud's registration if those types lack the attributes it asks for. The defaults match ThunderID's `Person` type. If your users have another type, set:

- `nextcloudDeps.registerJob.allowedUserTypes` and `tuwunelDeps.registerJob.allowedUserTypes` to its name.
- `nextcloudDeps.claims` to its attributes for the display name and email address, for example `{displayName: full_name, email: external_email}`.

## Talk meeting size

Talk calls are peer to peer, which suits meetings of about four or five people. Larger meetings need Nextcloud's High Performance Backend, which the chart doesn't include yet.

## Client setup

- [Calendar Client Setup](Calendar-Client-Setup.md): Thunderbird, Outlook, phones and tablets.
- [Files Client Setup](Files-Client-Setup.md): the Nextcloud desktop and mobile apps.
- Upstream user guides: [Nextcloud Calendar](https://docs.nextcloud.com/server/latest/user_manual/en/groupware/calendar.html), [Nextcloud Files](https://docs.nextcloud.com/server/latest/user_manual/en/files/index.html), [Nextcloud Talk](https://docs.nextcloud.com/server/latest/user_manual/en/talk/index.html), [Collabora Online](https://help.collaboraoffice.com/) and [Element](https://element.io/user-guide).

## Components

| Component | Version | Role |
|---|---|---|
| [Nextcloud](https://nextcloud.com/) | 34.0.4 | Calendars, files, webmail and Talk |
| [Collabora Online (CODE)](https://www.collaboraonline.com/code/) | 25.04.9.4.1 | Document editor |
| [coturn](https://github.com/coturn/coturn) | 4.18.0 | TURN server for Talk |
| [Tuwunel](https://github.com/matrix-construct/tuwunel) | 2.0.0 | Matrix homeserver |
| [LiveKit](https://livekit.io/) | 1.13.11 | Media server for Matrix calls |
| [lk-jwt-service](https://github.com/element-hq/lk-jwt-service) | 0.10.0 | Gives signed-in Matrix users a LiveKit token |
| [PostgreSQL](https://www.postgresql.org/) | 16 | Nextcloud's database |
| [Redis](https://redis.io/) | 7 | Nextcloud's cache |

Versions are the ones chart 0.6.1 ships.
