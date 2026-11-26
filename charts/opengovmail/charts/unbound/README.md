# Unbound subchart

Runs Unbound, a recursive DNS resolver, for Rspamd's DNS and blocklist (RBL) lookups. It resolves from the root servers itself, because Spamhaus, SURBL and other blocklists refuse queries that arrive through public resolvers.

## Turn it on

Unbound is off by default. Turn it on together with Rspamd and Redis in your values file. Rspamd waits for Redis at startup and doesn't start without it.

```yaml
rspamd:
  enabled: true
unbound:
  enabled: true
redis:
  enabled: true
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `image.repository` | string | `klutchell/unbound` | Image repository. A distroless image that runs as the `unbound` user. |
| `image.tag` | string | `1.26.2` | Image tag |
| `configuration` | string | Recursive `unbound.conf` (see `values.yaml`) | Required. Mounted over the image's `unbound.conf`. The default resolves recursively, validates DNSSEC, runs without chroot or a user switch, and sizes its caches for a 256Mi limit. |
| `resources.limits.memory` | string | `256Mi` | Memory limit. Lower it only together with the cache sizes in `configuration`. |
| `service.port` | int | `53` | Service port |

## Security

Unbound runs as the image's `unbound` user (UID 101, GID 102) from the start, with every Linux capability dropped, privilege escalation off and the `RuntimeDefault` seccomp profile.

- To let a non-root user listen on port 53, the pod sets the sysctl `net.ipv4.ip_unprivileged_port_start=53`. Kubernetes treats it as a safe sysctl since version 1.22, so the node needs no extra setting.
- `/var/unbound` belongs to the `unbound` user. It holds the root hints and the DNSSEC trust anchor, which Unbound rewrites when the root key rolls over.

## Test it

From any pod in the namespace:

```bash
# A normal lookup
dig @<release-name>-unbound example.com

# Spamhaus's test entry: answers 127.0.0.2 when blocklist lookups work
dig @<release-name>-unbound 2.0.0.127.zen.spamhaus.org +short
```

## Persistence

The StatefulSet mounts a volume at `/var/lib/unbound`. Unbound keeps its cache in memory and doesn't use that volume. It stays because Kubernetes doesn't allow a StatefulSet's volume claims to change in place, so removing it would break upgrades.
