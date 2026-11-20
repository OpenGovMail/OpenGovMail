# Upgrade OpenGovMail

## Overview

This procedure changes an OpenGovMail installation to a new chart version. Use the same procedure to apply a change in your values file. The two use `helm upgrade`.

## Before you start

- Read the entries in [CHANGELOG.md](../../CHANGELOG.md) between your version and the new version. If a release has breaking changes, its entry gives the necessary steps.
- Get your values file and the password of the ThunderID administrator.
- Before an upgrade to a new major version, back up the Raven volume. This volume has all the mailboxes.

## Upgrade

1. Get the new chart version. Replace `<version>` with the version number:

   ```bash
   cd OpenGovMail
   git fetch --tags
   git checkout v<version>
   helm dependency update ./charts/opengovmail
   ```

2. Get the ThunderID copy of the direct auth secret from the cluster:

   ```bash
   DIRECT_AUTH_SECRET=$(kubectl get secret opengovmail-raven-idp -n opengovmail \
     -o jsonpath='{.data.directAuthSecret}' | base64 -d)
   ```

3. Do the upgrade:

   ```bash
   helm upgrade opengovmail ./charts/opengovmail \
     -f ~/my-values.yaml \
     --set-string global.thunderAdmin.password='<your-admin-password>' \
     --set thunderid.configuration.server.security.directAuthSecret=$DIRECT_AUTH_SECRET \
     --namespace opengovmail
   ```

   Do not give `global.ravenIdp.clientSecret` or `global.ravenIdp.directAuthSecret`. The Raven Secret keeps the values that are in the cluster. If you give a different client secret, the upgrade stops.

4. Make sure that the pods show `Running` or `Completed`:

   ```bash
   kubectl get pods -n opengovmail
   ```

You can put the secrets in a second values file. Then you do not type them again. Do not commit this file. Give it to Helm with a second `-f` option.

## Roll back

If the new version does not operate correctly, return to the previous revision:

```bash
helm history opengovmail -n opengovmail
helm rollback opengovmail <revision> -n opengovmail
```

A rollback restores the Kubernetes resources. It does not restore data.

## See also

- [Installation guide](../installation-guide.md)
- [CHANGELOG.md](../../CHANGELOG.md)
- [Troubleshooting](../troubleshooting.md)
