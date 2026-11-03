# OpenGovMail Autodiscover Service

Serves the standard email-client autodiscovery documents for every domain in
`conf/opengovmail.yaml`:

- Thunderbird / Mozilla autoconfig (`config-v1.1.xml`)
- Microsoft Outlook Autodiscover (POX XML)

The service is stateless. It reads the same domain list as the rest of the
stack, so adding a mail domain only requires regenerating the configs and
restarting this container.

## Endpoints

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/healthz` | Health check |
| `GET` | `/mail/config-v1.1.xml?emailaddress=user@example.com` | Thunderbird autoconfig |
| `GET` | `/.well-known/autoconfig/mail/config-v1.1.xml?emailaddress=user@example.com` | Thunderbird autoconfig (well-known path) |
| `POST` | `/autodiscover/autodiscover.xml` | Outlook Autodiscover |
| `POST` | `/Autodiscover/Autodiscover.xml` | Outlook Autodiscover (capitalized path) |

Unknown domains return `404`. Malformed or oversized Outlook requests return
`400` / `413` with an Autodiscover error document. Requests with the wrong
method return `405`.

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `CONFIG_FILE` | `/etc/autodiscover/opengovmail.yaml` | Config file containing the `domains:` list |
| `PORT` | `8080` | Listen port |
| `MAIL_HOSTNAME` | `mail.<primary domain>` | IMAP/SMTP hostname advertised to clients |
| `TLS_ENABLED` | `false` | Serve HTTPS directly (Docker Compose mode) |
| `TLS_CERT_DIR` | `/etc/letsencrypt/live` | Parent of the certbot lineage directory |

When `TLS_ENABLED=true`, the service loads
`<TLS_CERT_DIR>/<primary domain>/fullchain.pem` and `privkey.pem`, requires TLS
1.2 or newer, and reloads the pair after certbot renews it. In Kubernetes TLS
is terminated by the Ingress, so `TLS_ENABLED` stays `false`.

## Local run

```bash
go test ./...
go run .
```

Point `CONFIG_FILE` at a config with real domains:

```bash
CONFIG_FILE=../../conf/opengovmail.yaml go run .
```

## Verify

```bash
curl -sS "http://localhost:8080/healthz"

curl -sS \
  "http://localhost:8080/.well-known/autoconfig/mail/config-v1.1.xml?emailaddress=user@example.com"

curl -sS -X POST \
  -H 'Content-Type: text/xml' \
  --data '<?xml version="1.0" encoding="utf-8"?>
<Autodiscover xmlns="http://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006">
  <Request>
    <EMailAddress>user@example.com</EMailAddress>
    <AcceptableResponseSchema>http://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a</AcceptableResponseSchema>
  </Request>
</Autodiscover>' \
  http://localhost:8080/autodiscover/autodiscover.xml
```

## Troubleshooting

**`no domains` at startup** — `CONFIG_FILE` is empty, unreadable, or its
`domains:` list has no entries. The first domain in the list is the primary
domain and determines the default `mail.<domain>` hostname.

**Clients cannot reach the endpoint** — confirm the DNS `A` record for the bare
domain points at the mail server, and that port 443 reaches this container. The
Compose deployment publishes `443:8443`; the container itself listens on 8443
when `PORT=8443`.

**TLS error at startup** — certbot must finish before this service starts. Check
that `<TLS_CERT_DIR>/<primary domain>/fullchain.pem` and `privkey.pem` exist and
are readable by the container user. Certbot writes the private key as
root-only, so the Compose service runs as root to read it; the Helm deployment
keeps the pod non-root because the Ingress terminates TLS.

**Outlook still asks for manual settings** — verify
`https://<domain>/autodiscover/autodiscover.xml` responds with XML (not HTML)
and that the certificate covers `<domain>`. Use Outlook's **Test E-Mail
AutoConfiguration** dialog to see which URL it tried.
