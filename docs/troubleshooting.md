# Troubleshooting OpenGovMail

In the commands on this page, `example.com` is your domain and `opengovmail` is the namespace of the installation.

## A certificate does not become ready

`kubectl get certificate -n opengovmail` shows `False` in the `READY` column for more than some minutes.

### Cause

Let's Encrypt cannot connect to the server on port 80 for the hostname. Usually, the DNS record does not point to the server yet, or port 80 is closed. With `le-prod`, a different cause is possible: the limit of 5 certificates for each set of names in one week.

### Solution

1. Find the cause that cert-manager gives:

   ```bash
   kubectl describe certificate <name> -n opengovmail
   kubectl get challenges -n opengovmail
   kubectl describe challenge <name> -n opengovmail
   ```

2. Make sure that `dig +short <hostname>` shows the IP address of the server.
3. Make sure that port 80 is open to the internet.
4. If the message is about a rate limit, set `global.tls.issuer` to `le-staging`. When the other parts of the installation operate correctly, set it to `le-prod` again.

### For more information

- [cert-manager ACME troubleshooting](https://cert-manager.io/docs/troubleshooting/acme/)
- [Let's Encrypt rate limits](https://letsencrypt.org/docs/rate-limits/)

## The ThunderID setup job does not complete on a first installation

The `helm install` command does not complete. `kubectl get pods -n opengovmail` shows the ThunderID setup pod in `ContainerCreating`. `kubectl describe pod` shows that the certificate Secret of ThunderID does not exist.

### Cause

The setup job uses the certificate in `thunderid.declarativeResources.secret.name`. The chart gets this certificate before the job only if the Secret name is `thunderid-tls`. Chart 0.6.1 does not do this for any name. With a different name or with chart 0.6.1, the certificate does not exist when the job starts.

### Solution for chart versions after 0.6.1

1. Stop the installation and delete the namespace:

   ```bash
   helm uninstall opengovmail -n opengovmail
   kubectl delete namespace opengovmail
   ```

   An installation that did not complete has no mail and no accounts. Thus, you lose no data. You must delete the namespace because the first installation made Secrets. If these Secrets stay, the chart does not accept the new secrets in the next installation.

2. In your values file, set `thunderid.declarativeResources.secret.name` to `thunderid-tls`.
3. Do your `helm install` command again.

### Solution for chart 0.6.1

Create the certificate before you install the chart. Then let Helm take control of the certificate.

1. Stop the installation and delete the namespace:

   ```bash
   helm uninstall opengovmail -n opengovmail
   kubectl delete namespace opengovmail
   ```

   An installation that did not complete has no mail and no accounts. Thus, you lose no data. You must delete the namespace because the first installation made Secrets. If these Secrets stay, the chart does not accept the new secrets in the next installation.

2. Create the namespace and the certificate. Use your domain and your issuer:

   ```bash
   kubectl create namespace opengovmail
   kubectl apply -n opengovmail -f - <<EOF
   apiVersion: cert-manager.io/v1
   kind: Certificate
   metadata:
     name: example-com-tls
   spec:
     secretName: example-com-tls
     issuerRef:
       name: le-prod
       kind: ClusterIssuer
     dnsNames:
       - example.com
   EOF
   kubectl wait --for=condition=Ready certificate/example-com-tls -n opengovmail --timeout=10m
   ```

3. Do your `helm install` command again. Add `--take-ownership` and remove `--create-namespace`. You must have Helm 3.17 or later.

## Pods show "no match for platform in manifest"

A pod stays in `ErrImagePull` or `ImagePullBackOff` and shows "no match for platform in manifest". Or the pod stops with `exec format error`.

### Cause

The node is not amd64. The OpenGovMail images (Postfix, OpenDKIM and Raven) are for amd64 only. They do not operate on ARM servers or on Apple Silicon Macs.

### Solution

Install OpenGovMail on an amd64 (x86_64) node. For development, use an amd64 test server.

## Pods stay Pending after you add the collaboration services

Some pods stay `Pending`. Frequently, the cert-manager pods are among them. `kubectl describe pod` shows `Insufficient cpu`.

### Cause

The collaboration services request more CPU than one CPU can supply.

### Solution

Do one of these:

- Give the server a minimum of 2 CPUs.
- Remove the collaboration services that you do not use.

## Mail apps cannot sign in

The mail app does not accept the certificate of the server. Or the server does not accept a correct password.

### Cause

The installation uses the `le-staging` issuer. Mail apps do not trust staging certificates. Also, the connections from Raven to ThunderID fail the TLS check.

A different cause is possible: the user is in an organization unit that has no mail domain. Users in such an organization unit cannot sign in.

### Solution

1. If the issuer is `le-staging`, set `global.tls.issuer` to `le-prod`. Then do an upgrade. Read [Upgrade OpenGovMail](how-to/upgrade.md).
2. If the issuer is `le-prod`, make sure that the organization unit of the user has a mail domain. [Add users](how-to/add-users.md) tells you about organization units.

## Mail goes to spam or fails DKIM

Recipients find your mail in their spam folders. Or [mail-tester](https://www.mail-tester.com/) shows a DKIM, SPF or reverse DNS failure.

### Cause

A DNS record does not exist or is not correct. Usually, the cause is the DKIM record. The command shows the DKIM key in some parts, and you must join these parts.

### Solution

1. Examine each record from outside your network:

   ```bash
   dig +short TXT mail._domainkey.example.com
   dig +short TXT example.com
   dig +short -x 12.34.56.78
   ```

2. Make sure that the DKIM value is one text that starts with `v=DKIM1`. The key must have no quotation marks or spaces in it.
3. Make sure that the SPF record has the IP address of the server.
4. Make sure that the PTR record shows `mail.example.com`.
5. Examine the IP address of the server with the [MXToolbox blocklist check](https://mxtoolbox.com/blacklists.aspx). Sometimes, a new IP address from a cloud provider is on a blocklist already.

### For more information

- [Installation guide, Step 1 and Step 7](installation-guide.md#step-1-add-dns-records)

## Mail to other servers does not arrive

Outgoing mail stays in the queue. The Postfix logs show timeouts on port 25.

### Cause

Your hosting provider closes outbound port 25. Many cloud providers do this by default.

### Solution

1. Ask your provider to open outbound port 25.
2. To examine the Postfix logs, find the Postfix pod and show its logs:

   ```bash
   kubectl get pods -n opengovmail
   kubectl logs -n opengovmail <postfix-pod>
   ```
