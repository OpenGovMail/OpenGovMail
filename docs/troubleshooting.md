# Troubleshooting OpenGovMail

In the commands below, `example.com` stands for your domain and `opengovmail` for the namespace you installed into.

## A certificate stays not ready

`kubectl get certificate -n opengovmail` shows `READY` as `False` for several minutes.

### Cause

Let's Encrypt couldn't reach the server over port 80 for the hostname, usually because the DNS record doesn't point at the server yet or port 80 is closed. On `le-prod`, you may also have hit the limit of 5 certificates per set of names per week.

### Solution

1. Read the reason cert-manager gives:

   ```bash
   kubectl describe certificate <name> -n opengovmail
   kubectl get challenges -n opengovmail
   kubectl describe challenge <name> -n opengovmail
   ```

2. Check that `dig +short <hostname>` prints the server's IP address, and that port 80 is open to the internet.
3. If the message mentions a rate limit, switch `global.tls.issuer` to `le-staging` until the rest of the setup works, then back to `le-prod`.

### For more information

- [cert-manager ACME troubleshooting](https://cert-manager.io/docs/troubleshooting/acme/)
- [Let's Encrypt rate limits](https://letsencrypt.org/docs/rate-limits/)

## ThunderID setup job never finishes on a first install

On chart 0.6.1, `helm install` doesn't return. `kubectl get pods -n opengovmail` shows the ThunderID setup pod stuck in `ContainerCreating`, and `kubectl describe pod` reports that the `example-com-tls` Secret is missing.

### Cause

The setup job mounts the certificate for your domain, but the chart only creates that certificate after the job finishes.

### Solution

Create the certificate before installing, then let Helm take it over.

1. Stop the stuck install and delete the namespace. A first install that never finished holds no mail or accounts. Deleting the namespace also removes the Secrets the first attempt created, which would otherwise make the chart reject the new secrets you pass when you retry:

   ```bash
   helm uninstall opengovmail -n opengovmail
   kubectl delete namespace opengovmail
   ```

2. Create the namespace and the certificate, using your domain and issuer:

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

3. Run your `helm install` command again, with `--take-ownership` added and `--create-namespace` removed. This needs Helm 3.17 or later.

## Pods fail with "no match for platform in manifest"

A pod stays in `ErrImagePull` or `ImagePullBackOff` with "no match for platform in manifest", or crashes with `exec format error`.

### Cause

The node isn't amd64. The OpenGovMail images (Postfix, OpenDKIM and Raven) are built for amd64 only, so they don't run on ARM servers or on Apple Silicon Macs.

### Solution

Run OpenGovMail on an amd64 (x86_64) node. For local development, use an amd64 machine or virtual machine.

## Pods stay Pending after turning on collaboration

Some pods, often including cert-manager's, stay `Pending`. `kubectl describe pod` reports `Insufficient cpu`.

### Cause

With collaboration turned on, the pods' CPU requests add up to more than one CPU.

### Solution

Give the server at least 2 CPUs, or turn off the collaboration services you don't need.

## Mail clients can't sign in

Clients reject the server's certificate, or a correct password is refused.

### Cause

The install is still on the `le-staging` issuer. Staging certificates aren't trusted, so mail clients reject them, and Raven's calls to ThunderID fail TLS verification.

### Solution

Set `global.tls.issuer` to `le-prod` and run an upgrade. See [Upgrade OpenGovMail](how-to/upgrade.md).

If the issuer is already `le-prod`, check that the user belongs to an organization unit that owns a mail domain. Users of an organization unit with no domain can't sign in. [Add users](how-to/add-users.md) explains organization units.

## Mail goes to spam or fails DKIM

Messages you send land in recipients' spam folders, or [mail-tester](https://www.mail-tester.com/) reports a DKIM, SPF or reverse DNS failure.

### Cause

One of the DNS records is missing or malformed. The DKIM record is the most common: its key is printed in several quoted parts that must be joined.

### Solution

1. Check each record from outside:

   ```bash
   dig +short TXT mail._domainkey.example.com
   dig +short TXT example.com
   dig +short -x 12.34.56.78
   ```

2. The DKIM value must be one string starting `v=DKIM1`, with no quotes or spaces inside the key.
3. The SPF record must include the server's IP address, and the PTR record must return `mail.example.com`.
4. Check the server's IP address on [MXToolbox's blocklist check](https://mxtoolbox.com/blacklists.aspx). A new IP address from a cloud provider is sometimes listed already.

### For more information

- [Installation guide, Step 1 and Step 7](installation-guide.md#step-1-add-dns-records)

## Mail to other servers never arrives

Outgoing mail stays queued, and the Postfix logs show connection timeouts on port 25.

### Cause

Your hosting provider blocks outbound port 25. Many cloud providers do this by default.

### Solution

Ask the provider to open outbound port 25. Check the Postfix logs with:

```bash
kubectl get pods -n opengovmail
kubectl logs -n opengovmail <postfix-pod>
```
