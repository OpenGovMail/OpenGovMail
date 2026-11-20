# Helm values

This page gives the values that you will possibly set in your values file. All the values, with comments, are in [values.yaml](../../charts/opengovmail/values.yaml) and in the `values.yaml` of each subchart in [charts/opengovmail/charts/](../../charts/opengovmail/charts). [values-example.yaml](../../charts/opengovmail/values-example.yaml) is a complete values file. You can start from it.

Give secrets on the command line. Do not put them in the values file. Read [Secrets](#secrets).

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
| `thunderid.declarativeResources.secret.name` | The certificate Secret that ThunderID uses. With `thunderid-tls`, the chart gets the certificate before the ThunderID setup job. With a different name, you must supply the Secret yourself. | `thunderid-tls` |

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
| `tuwunel.serverName` | The part after the colon in Matrix IDs. Empty uses `global.domain`. You cannot change it after the first start. A change makes it necessary to delete the database. | `example.com` |
| `tuwunel.host` | Matrix hostname. Empty uses `matrix.<serverName>`. | `matrix.example.com` |
| `tuwunel.matrixRtc.enabled` | Run LiveKit for Matrix calls. | `false` |
| `tuwunelDeps.registerJob.allowedUserTypes` | ThunderID user types that can sign in to Matrix. | `["Person"]` |

## Secrets

Give these values with `--set` or `--set-string`. Or put them in a values file that you do not commit.

| Field | Description | When to set |
|---|---|---|
| `global.thunderAdmin.password` | Password for the ThunderID `admin` user. Empty generates one on the first install and keeps it. Set it explicitly when rendering with `helm template`. | Install and upgrade |
| `global.ravenIdp.clientSecret` | The secret that Raven uses to sign in to ThunderID. If you give a different value in an upgrade, the upgrade stops. | First install only |
| `global.ravenIdp.directAuthSecret` | The secret that protects the password check of ThunderID. This is the copy for Raven. | First install only |
| `thunderid.configuration.server.security.directAuthSecret` | The same secret. This is the copy for ThunderID. It must be the same as the copy for Raven. | Install and upgrade |
