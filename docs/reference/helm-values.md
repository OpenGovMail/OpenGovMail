# Helm values

This page describes the values you are most likely to set in your values file. Every value, with comments, is in [values.yaml](../../charts/opengovmail/values.yaml), and each subchart's own `values.yaml` under [charts/opengovmail/charts/](../../charts/opengovmail/charts). [values-example.yaml](../../charts/opengovmail/values-example.yaml) is a complete values file to start from.

Secrets go on the command line, not in the values file. See [Secrets](#secrets).

## Domain and certificates

| Field | Description | Example |
|---|---|---|
| `global.domain` | Required. The mail domain. Hostnames, certificate names and service settings all come from it. | `example.com` |
| `global.mailHostname` | Hostname of the mail server. Empty uses `mail.<domain>`. | `mail.example.com` |
| `global.tls.enabled` | Whether the chart requests certificates from cert-manager. | `true` |
| `global.tls.issuer` | ClusterIssuer to use. Defaults to `le-staging`. Production needs `le-prod`. | `le-prod` |
| `global.tls.renewBefore` | How long before expiry to renew. | `720h` |
| `global.tls.domains` | Names on the certificates. Empty uses `mail.<domain>` and `<domain>`. Set only to add names. | `[]` |

## Mail

| Field | Description | Example |
|---|---|---|
| `postfix.enabled` | Run Postfix. | `true` |
| `opendkim.enabled` | Run OpenDKIM. | `true` |
| `opendkim.domains` | Domains to sign for. Empty signs for `global.domain` with selector `mail`. | `[{domain: example.com, dkimSelector: mail}]` |
| `raven.enabled` | Run Raven. | `true` |
| `raven.organizationUnitDomains` | Mail domains each ThunderID organization unit owns, keyed by its ID, main domain first. Empty gives the bundled `default` unit `global.domain`. Required with an external ThunderID. | `{"<ou-id>": ["example.com"]}` |
| `raven.persistence.size` | Size of Raven's mailbox volume. | `5Gi` |

## Spam filtering

| Field | Description | Example |
|---|---|---|
| `rspamd.enabled` | Run Rspamd and check incoming mail with it. | `false` |
| `redis.enabled` | Run Redis for Rspamd. | `false` |
| `unbound.enabled` | Run Unbound, a DNS resolver for Rspamd's blocklist lookups. | `false` |

## Identity

| Field | Description | Example |
|---|---|---|
| `thunderid.enabled` | Run the bundled ThunderID. | `true` |
| `global.thunder.url` | An external ThunderID to use instead, with `thunderid.enabled: false`. Empty uses the bundled one at `https://<domain>:8090`. | `https://thunder.example.com` |
| `raven.thunder.resourceIdentifier` | The resource server Raven asks ThunderID for a token to. Keep the default with the bundled ThunderID. With an external one, use the identifier it registered for its System resource server. | `https://localhost:8090/mcp` |
| `thunderid.configuration.server.publicUrl` | ThunderID's public address. Must match `global.domain`. | `https://example.com:8090` |
| `thunderid.configuration.gateClient.hostname` | ThunderID's sign-in page hostname. Must match `global.domain`. | `example.com` |
| `thunderid.declarativeResources.secret.name` | Certificate Secret ThunderID serves. Your domain with dots as dashes, then `-tls`. | `example-com-tls` |

## Storage

| Field | Description | Example |
|---|---|---|
| `seaweedfs.enabled` | Run SeaweedFS for attachments and Nextcloud files. | `true` |
| `seaweedfs.persistence.size` | Size of the object storage volume. | `10Gi` |
| `global.s3.secretKey` | S3 key shared by SeaweedFS and Raven. Empty derives one from the release name. Set your own in production. | |

## Collaboration

| Field | Description | Example |
|---|---|---|
| `nextcloud.enabled` | Run Nextcloud, with Collabora and Talk. | `false` |
| `nextcloud.nextcloud.host` | Nextcloud's hostname. Must match its ingress. | `cloud.example.com` |
| `nextcloudDeps.turn.externalIP` | The node's public IP address, if it only has a private one. | `12.34.56.78` |
| `nextcloudDeps.registerJob.allowedUserTypes` | ThunderID user types that can sign in to Nextcloud. | `["Person"]` |
| `nextcloudDeps.claims` | ThunderID attributes Nextcloud reads for the display name, email address and other fields. | `{displayName: full_name, email: external_email}` |
| `tuwunel.enabled` | Run Tuwunel, the Matrix server. | `false` |
| `tuwunel.serverName` | The part after the colon in Matrix IDs. Empty uses `global.domain`. It can't change after the first start without wiping the database. | `example.com` |
| `tuwunel.host` | Matrix hostname. Empty uses `matrix.<serverName>`. | `matrix.example.com` |
| `tuwunel.matrixRtc.enabled` | Run LiveKit for Matrix calls. | `false` |
| `tuwunelDeps.registerJob.allowedUserTypes` | ThunderID user types that can sign in to Matrix. | `["Person"]` |

## Secrets

Pass these with `--set` or `--set-string`, or from a values file you never commit.

| Field | Description | When to set |
|---|---|---|
| `global.thunderAdmin.password` | Password for the ThunderID `admin` user. Empty generates one on the first install and keeps it. Set it explicitly when rendering with `helm template`. | Install and upgrade |
| `global.ravenIdp.clientSecret` | Raven's secret for signing in to ThunderID. A different value on upgrade fails the render. | First install only |
| `global.ravenIdp.directAuthSecret` | Secret protecting ThunderID's password check, Raven's copy. | First install only |
| `thunderid.configuration.server.security.directAuthSecret` | The same secret, ThunderID's copy. Must equal Raven's. | Install and upgrade |
