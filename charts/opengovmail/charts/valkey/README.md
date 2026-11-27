# Valkey Helm Chart

[Valkey](https://valkey.io/) in-memory data store for rspamd backend storage. Valkey is an open-source fork of Redis and speaks the same protocol.

## Installation

```bash
helm install valkey ./valkey \
  --set persistence.enabled=true \
  --set persistence.size=1Gi
```

## Configuration

Key values:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `replicaCount` | int | `1` | Number of replicas. **Must stay 1** — this chart has no clustering/Sentinel; values > 1 create independent, non-replicated Valkey instances and will corrupt rspamd's shared Bayes/learning data |
| `image.repository` | string | `valkey/valkey` | Image repository |
| `image.tag` | string | `9.1-alpine` | Image tag |
| `persistence.enabled` | bool | `true` | Enable persistent volume |
| `persistence.size` | string | `1Gi` | Persistent volume size |
| `service.port` | int | `6379` | Service port |
| `securityContext.fsGroup` | int | `999` | Valkey user UID |

## Usage

Access Valkey within the cluster:

```bash
valkey-cli -h <release-name>-valkey -p 6379
```

For port-forward:

```bash
kubectl port-forward svc/<release-name>-valkey 6379:6379
valkey-cli -h localhost -p 6379
```

## Persistence

Valkey data is stored in `/data` via `volumeClaimTemplates`. The PVC persists across pod restarts.

## Dependencies

The Valkey chart has no external dependencies.
