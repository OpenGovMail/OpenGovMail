# Tuwunel Helm Chart

[Tuwunel](https://matrix-construct.github.io/tuwunel/) Matrix homeserver for OpenGovMail chat. Users sign in through ThunderID with their mail account, so `alice@example.com` chats as `@alice:example.com`.

Tuwunel embeds RocksDB, so there is no database to run. It is a single-replica StatefulSet with one volume for the database and uploaded media, and a second one for backups.

## Why first-party templates

The community [`tuwunel-helm`](https://github.com/AreYouLoco/tuwunel-helm) chart was considered. It is not maintained by the Tuwunel team, and the workload is small: a StatefulSet, Service, ConfigMap, Ingress and two PVCs. Owning those templates keeps the image pinned, the config in one reviewed file, and the ThunderID wiring in this repository.

## Installation

Through the umbrella chart, which also registers the client in ThunderID and creates its Secret:

```bash
helm upgrade --install opengovmail ./charts/opengovmail -n mail \
  --set tuwunel.enabled=true
```

This needs:

- A ThunderID: the bundled one, or an external one in `global.thunder.url`. It must serve a publicly trusted certificate.
- `raven.enabled`: its machine-to-machine client registers the Tuwunel application.
- DNS for `matrix.<domain>` and `<domain>` pointing at the cluster (see [DNS](#dns)).

The registration Job (`<release>-tuwunel-oidc-register`) creates an application called "Matrix chat" in ThunderID with client id `TUWUNEL`, and updates it on every upgrade.

Standalone, without the umbrella, create the ThunderID application yourself and a Secret named by `oidc.existingSecret` with the key `clientSecret`. Register this redirect URI:

```
https://matrix.<domain>/_matrix/client/unstable/login/sso/callback/TUWUNEL
```

## How sign-in works

Password login and registration are off. A Matrix client shows a "ThunderID" button, which sends the user to ThunderID. On first sign-in Tuwunel creates the account, taking the localpart from ThunderID's `username` claim. If that claim is missing, it falls back to the part of `email` before the `@`.

A localpart that is already taken is an error rather than a random ID (`unique_id_fallbacks = false`). The first account to sign in becomes the server admin and joins `#admins:<domain>`.

ThunderID only releases a claim the application lists for userinfo and maps to a requested scope. The registration Job sets both, mapping `username` to the `profile` scope. The user type must define a `username` attribute; ThunderID's built-in `Person` type does.

Only user types in `tuwunelDeps.registerJob.allowedUserTypes` can sign in. Add the user type of each mail domain's organization unit there.

## DNS

| Record | Points at | Why |
|---|---|---|
| `matrix.<domain>` | the cluster | Client API, and federation when enabled |
| `<domain>` | the cluster | Serves `/.well-known/matrix/client` and `/server`, so clients find the server from an address like `@alice:<domain>` |

If `<domain>` is hosted elsewhere, for example a website, set `wellKnown.enabled=false`. Then serve these two files from that host instead:

```
https://<domain>/.well-known/matrix/client  {"m.homeserver":{"base_url":"https://matrix.<domain>/"}}
https://<domain>/.well-known/matrix/server  {"m.server":"matrix.<domain>:443"}
```

Federation goes over 443 through that delegation, so port 8448 is not opened.

## Configuration

| Key | Default | Description |
|-----|---------|-------------|
| `serverName` | `global.domain` | Matrix server name. **Cannot change after the first start without wiping the database.** |
| `host` | `matrix.<serverName>` | Hostname Tuwunel is served on |
| `image.repository` | `ghcr.io/matrix-construct/tuwunel` | Image repository |
| `image.tag` | chart `appVersion` (`v2.0.0`) | Image tag |
| `federation.enabled` | `false` | Federate with other Matrix servers. Off until the federation policy is decided. |
| `maxRequestSize` | `24 MiB` | Largest upload |
| `oidc.enabled` | `true` | Sign in through ThunderID; turns password login off |
| `oidc.clientId` | `TUWUNEL` | Client id registered in ThunderID |
| `oidc.existingSecret` | `tuwunel-idp` | Secret with the client secret under `clientSecret` |
| `oidc.issuerUrl` | `global.thunder.url`, else `https://<domain>:8090` | ThunderID issuer. Must never change. |
| `oidc.useridClaims` | `["username", "email"]` | Claims the localpart is taken from |
| `persistence.size` | `10Gi` | Database and media volume |
| `backup.enabled` | `true` | Online backups to a second volume |
| `backup.keep` | `3` | Backups kept |
| `backup.persistence.size` | `5Gi` | Backup volume |
| `ingress.enabled` | `true` | Ingress for `host` |
| `ingress.className` | `""` | IngressClass; empty uses the cluster default |
| `wellKnown.enabled` | `true` | Ingress for `<serverName>/.well-known/matrix` |
| `extraEnv` | `[]` | Extra `TUWUNEL_*` variables, which override the config file |

Ingress TLS certificates come from cert-manager through `global.tls.issuer`. Traefik passes the `X-Matrix` federation header unchanged and has no request body limit by default. Behind nginx, set `proxy-body-size` to at least `maxRequestSize`.

## Testing

```bash
helm test opengovmail -n mail
```

The test checks the client API, the `.well-known` answer and, with OIDC on, that SSO is the offered login.

To sign in, open [app.element.io](https://app.element.io), choose to sign in, and change the homeserver to `matrix.<domain>`. Then click the ThunderID button.

## Backups

The Matrix database is backed up with Tuwunel's online backup, to the `backups` volume. Backups run on request. As the server admin, send this in the `#admins:<domain>` room:

```
!admin server backup-database
!admin server list-backups
!admin server verify-backup
```

Backups hold the database only. Uploaded media lives under `/var/lib/tuwunel/media` on the data volume. Both volumes sit on the same node, so also take VolumeSnapshots, or copy them off the node, to survive losing it.

### Restoring

Restore replaces the database with the newest backup. Stop Tuwunel, run it once with `--restore-backup`, then start it again:

```bash
NS=mail; STS=opengovmail-tuwunel
kubectl -n $NS scale statefulset $STS --replicas=0
kubectl -n $NS wait --for=delete pod/$STS-0 --timeout=5m

IMAGE=$(kubectl -n $NS get statefulset $STS -o jsonpath='{.spec.template.spec.containers[0].image}')
kubectl -n $NS apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: tuwunel-restore
spec:
  restartPolicy: Never
  securityContext: {runAsUser: 1000, runAsGroup: 1000, runAsNonRoot: true, fsGroup: 1000}
  containers:
    - name: restore
      image: $IMAGE
      args: ["--restore-backup", "--maintenance", "--execute", "server shutdown"]
      env: [{name: TUWUNEL_CONFIG, value: /etc/tuwunel/tuwunel.toml}]
      volumeMounts:
        - {name: config, mountPath: /etc/tuwunel/tuwunel.toml, subPath: tuwunel.toml}
        - {name: oidc, mountPath: /etc/tuwunel/oidc}
        - {name: data, mountPath: /var/lib/tuwunel}
        - {name: backups, mountPath: /var/lib/tuwunel-backups}
  volumes:
    - {name: config, configMap: {name: $STS}}
    - {name: oidc, secret: {secretName: tuwunel-idp, optional: true}}
    - {name: data, persistentVolumeClaim: {claimName: data-$STS-0}}
    - {name: backups, persistentVolumeClaim: {claimName: backups-$STS-0}}
EOF
kubectl -n $NS wait --for=jsonpath='{.status.phase}'=Succeeded pod/tuwunel-restore --timeout=10m
kubectl -n $NS logs tuwunel-restore | grep -i restored   # "Restored database backup"
kubectl -n $NS delete pod tuwunel-restore

kubectl -n $NS scale statefulset $STS --replicas=1
```

Use `--restore-backup=<id>` for a specific backup from `list-backups`. Copy the volume first if you may need to undo the restore.

## Upgrading Tuwunel

Bump `appVersion` in `Chart.yaml`. Database migrations run on start and the listener opens only after they finish, so the startup probe allows 30 minutes. Take a backup first.
