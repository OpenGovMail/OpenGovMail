# Upgrade OpenGovMail

## Overview

This guide moves an OpenGovMail install to a new chart version, or applies a change to your values file. Both use `helm upgrade`.

## Before you start

- Read the entries in [CHANGELOG.md](../../CHANGELOG.md) between your version and the new one. A release with breaking changes lists the steps it needs.
- Have your values file and the ThunderID admin password at hand.
- Back up Raven's volume, which holds every mailbox, before a major version upgrade.

## Upgrade

1. Get the new chart version:

   ```bash
   cd OpenGovMail
   git fetch --tags
   git checkout v<version>
   helm dependency update ./charts/opengovmail
   ```

2. Read back ThunderID's copy of the direct auth secret from the cluster:

   ```bash
   DIRECT_AUTH_SECRET=$(kubectl get secret opengovmail-raven-idp -n opengovmail \
     -o jsonpath='{.data.directAuthSecret}' | base64 -d)
   ```

3. Run the upgrade:

   ```bash
   helm upgrade opengovmail ./charts/opengovmail \
     -f ~/my-values.yaml \
     --set-string global.thunderAdmin.password='<your-admin-password>' \
     --set thunderid.configuration.server.security.directAuthSecret=$DIRECT_AUTH_SECRET \
     --namespace opengovmail
   ```

   Leave out `global.ravenIdp.clientSecret` and `global.ravenIdp.directAuthSecret`. Raven's Secret keeps the values already in the cluster, and a different client secret fails the upgrade.

4. Check that the pods come back:

   ```bash
   kubectl get pods -n opengovmail
   ```

To keep from retyping the secrets, put them in a second values file that you never commit, and pass it with another `-f`.

## Roll back

If the new version misbehaves, go back to the previous revision:

```bash
helm history opengovmail -n opengovmail
helm rollback opengovmail <revision> -n opengovmail
```

A rollback restores Kubernetes resources, not data.

## See also

- [Installation guide](../installation-guide.md)
- [CHANGELOG.md](../../CHANGELOG.md)
- [Troubleshooting](../troubleshooting.md)
