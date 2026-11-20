# Glossary

| Term | Abbreviation | Definition |
|---|---|---|
| CalDAV and CardDAV | | Protocols calendar and contact apps use to sync with a server. Nextcloud serves both. |
| cert-manager | | A Kubernetes add-on that gets and renews TLS certificates, here from Let's Encrypt. |
| ClusterIssuer | | A cert-manager resource that says where certificates come from. OpenGovMail uses `le-staging` and `le-prod`. |
| DomainKeys Identified Mail | DKIM | A signature added to outgoing mail, checked against a public key in DNS, that proves the mail came from your server. OpenDKIM signs it. |
| Domain-based Message Authentication, Reporting and Conformance | DMARC | A DNS record that tells other servers what to do with mail that fails SPF or DKIM. |
| Helm | | The package manager for Kubernetes. OpenGovMail is installed as one Helm chart. |
| Helm values | | Settings passed to the chart, in a values file or with `--set`. |
| HTTP-01 challenge | | How Let's Encrypt checks you control a hostname: it fetches a file from the server over port 80. |
| Internet Message Access Protocol | IMAP | The protocol mail apps use to read mail. Raven serves it on ports 993 and 143. |
| k3s | | A small Kubernetes distribution, and the tested platform for OpenGovMail. |
| Local Mail Transfer Protocol | LMTP | The protocol Postfix uses to hand incoming mail to Raven. |
| Mail exchanger record | MX | The DNS record that tells other servers where to deliver mail for your domain. |
| Matrix | | An open chat protocol. Tuwunel is OpenGovMail's Matrix server, and Element a Matrix app. |
| Organization unit | OU | A part of your organization in ThunderID, such as a department. Each one owns one or more mail domains. |
| Pointer record | PTR | Reverse DNS: maps the server's IP address back to `mail.<domain>`. Other servers check it before accepting mail. |
| Raven | | OpenGovMail's mail store. It serves IMAP, takes delivery over LMTP and checks passwords for Postfix. |
| Simple Authentication and Security Layer | SASL | How Postfix checks a mail app's password before it accepts mail to send. Raven answers the check. |
| Simple Mail Transfer Protocol | SMTP | The protocol for sending mail between servers (port 25) and from mail apps (port 587). |
| Sender Policy Framework | SPF | A DNS record listing the IP addresses allowed to send mail for your domain. |
| ThunderID | | The identity server that holds OpenGovMail's users, groups and organization units. |
| Traversal Using Relays around NAT | TURN | A relay that lets call media through firewalls and NAT. coturn runs it for Nextcloud Talk. |
| User type | | A ThunderID definition of the attributes a kind of user has. The default is `Person`. Also called a user schema. |
