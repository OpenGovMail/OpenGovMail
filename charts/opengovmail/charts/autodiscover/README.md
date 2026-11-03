# Autodiscover Helm Chart

Deploys the OpenGovMail autodiscovery service, which answers Thunderbird
autoconfig and Outlook Autodiscover requests for the configured mail domains.

TLS is terminated by the Ingress. The pod listens on plain HTTP port 8080.

## Values

| Value | Default | Description |
|-------|---------|-------------|
| `domains` | `[]` | Domains to serve. Empty derives `global.domain`. The first entry is primary. |
| `ingress.enabled` | `true` | Create the Ingress for the bare domain and `mail.<domain>`. |
| `ingress.hosts` | `[]` | Override the served hosts. Empty derives `[<domain>, mail.<domain>]`. |
| `ingress.tls.secretNames` | `[]` | Override TLS Secret names per host. Empty derives `<host-with-dashes>-tls`. |
| `service.port` | `8080` | Container and Service port. |
| `image.tag` | `main` | Image tag. |

## TLS

By default the Ingress reuses the cert-manager Secrets created by the umbrella
chart:

- `<domain>` -> `<domain-with-dashes>-tls`
- `mail.<domain>` -> `mail-<domain-with-dashes>-tls`

No additional `Certificate` objects are created here. If you override
`ingress.hosts`, either add the hosts to `global.tls.domains` or point
`ingress.tls.secretNames` at Secrets that cover them.

## Install

Installed by the umbrella chart. See
[charts/README.md](../../../README.md).

## Test

```bash
helm test opengovmail -n mail
```
