# ClamAV subchart

Runs the ClamAV daemon (clamd) that Rspamd uses to scan incoming mail for viruses. freshclam runs in the same container and keeps the signatures up to date.

## Turn it on

ClamAV is off by default. Turn it on together with Rspamd in your values file:

```yaml
rspamd:
  enabled: true
clamav:
  enabled: true
```

Rspamd reaches clamd at the Service named by `global.serviceNames.clamav`, `clamav-server` by default. That matches `rspamd.dependencies.clamav.host`, so the two connect without more settings. If you change one, change the other.

To run Rspamd without virus scanning, leave `clamav.enabled` off and set `rspamd.modules.antivirus.enabled: false`. Otherwise Rspamd logs a failed scan for every message.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `image.repository` | string | `clamav/clamav` | Image repository |
| `image.tag` | string | `1.5.5` | Image tag. The published tag is amd64 only. |
| `persistence.enabled` | bool | `true` | Keep the signature database on a volume, so a restart doesn't download it again |
| `persistence.size` | string | `2Gi` | Size of the signature database volume |
| `persistence.storageClass` | string | `""` | Storage class. Empty uses the cluster default. |
| `resources.requests.memory` | string | `1Gi` | Memory request |
| `resources.limits.memory` | string | `2867Mi` | Memory limit. clamd holds the signatures in memory, and briefly twice that while it reloads them. |
| `env.noFreshclamd` | bool | `false` | `true` stops freshclam, so signatures are never updated |
| `service.port` | int | `3310` | clamd TCP port |
| `startupProbe.failureThreshold` | int | `60` | Startup checks, 10 seconds apart, before Kubernetes restarts clamd. The first start downloads about 300 MB of signatures. |

## Security

clamd and freshclam run as the image's `clamav` user (UID 100, GID 101), with every Linux capability dropped and privilege escalation off.

- A short init container runs as root with only the `CHOWN` capability. It gives the `clamav` user ownership of the signature volume, because the image's entrypoint changes that ownership itself and fails on a volume it doesn't own.
- `/run/clamav` and `/run/lock` are `emptyDir` volumes, because the entrypoint writes to both and `/run` belongs to root.

## Test it

After an install, `helm test` runs a pod that sends `PING` to clamd and expects `PONG`:

```bash
helm test opengovmail -n opengovmail
```

To check that scanning works, scan the EICAR test file inside the pod:

```bash
kubectl exec -n opengovmail <release>-clamav-0 -- sh -c \
  'printf "X5O!P%%@AP[4\\PZX54(P^)7CC)7}\$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!\$H+H*" > /tmp/eicar && clamdscan --stream --no-summary /tmp/eicar'
```

The result is `/tmp/eicar: Eicar-Test-Signature FOUND`.
