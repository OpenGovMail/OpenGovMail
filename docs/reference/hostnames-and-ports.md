# Hostnames and ports

This page lists the hostnames OpenGovMail needs in DNS and the ports it needs open on the server's firewall. `example.com` stands for `global.domain`.

## Hostnames

Point each hostname at the server with an A record. cert-manager gets a certificate for each.

| Hostname | Set with | Serves | Needed for |
|---|---|---|---|
| `mail.example.com` | `global.mailHostname` (default `mail.<domain>`) | SMTP and IMAP | Always |
| `example.com` | `global.domain` | ThunderID on port 8090, and Matrix's `/.well-known/matrix/client` | Always |
| Your choice, such as `cloud.example.com` | `nextcloud.nextcloud.host` and `nextcloud.ingress.tls` | Nextcloud, and Talk's TURN server | Nextcloud |
| Your choice, such as `office.example.com` | `nextcloud.collabora.ingress`, `nextcloud.collabora.collabora.server_name` and `nextcloud.collabora.collabora.aliasgroups` | Collabora. Users never open it directly. | Nextcloud |
| `matrix.example.com` | `tuwunel.host` | Matrix | Matrix |
| `matrix-rtc.example.com` | `tuwunel.matrixRtc.host` | Matrix calls (LiveKit and its token service) | Matrix calls |

## Ports

Open these inbound on the node's firewall.

| Port | Service | Needed for |
|---|---|---|
| 25/tcp | Postfix: mail from other servers | Always |
| 80/tcp | Traefik: Let's Encrypt HTTP-01 challenges | Always |
| 443/tcp | Traefik: Nextcloud, Collabora and Matrix | Collaboration |
| 587/tcp | Postfix: mail clients sending mail | Always |
| 143/tcp | Raven: IMAP with STARTTLS | Optional; 993 is preferred |
| 993/tcp | Raven: IMAP over TLS | Always |
| 8090/tcp | ThunderID: console, sign-in and tokens | Always |
| 3478/udp and 3478/tcp | coturn: Talk's TURN server | Nextcloud Talk |
| 49160-49200/udp | coturn: Talk's relay (`nextcloudDeps.turn.minPort` and `maxPort`) | Nextcloud Talk |
| 7882/udp | LiveKit: Matrix call media | Matrix calls |
| 7881/tcp | LiveKit: Matrix call media on networks that block UDP | Matrix calls |

Postfix, Raven, coturn and LiveKit bind their ports directly on the node. If the node only has a private IP address, set `nextcloudDeps.turn.externalIP` and `tuwunel.matrixRtc.livekit.nodeIP` to its public address.
