# Add users

## Overview

This guide adds mail users in the ThunderID console that ships with OpenGovMail. Each user you add gets a mailbox and can sign in to every OpenGovMail service. [Identity and users](../concepts/identity-and-users.md) explains organization units, user types and groups.

## Before you start

- OpenGovMail is installed. See the [installation guide](../installation-guide.md).
- You can reach the console at `https://<your-domain>:8090/console`.
- You have the ThunderID `admin` password. If the chart generated it, read it with:

  ```bash
  kubectl get secret opengovmail-thunder-admin -n opengovmail \
    -o jsonpath='{.data.password}' | base64 -d
  ```

## Add a user

1. Open `https://<your-domain>:8090/console` and sign in as `admin`.
2. Choose the organization unit the user belongs to. The `default` unit owns your main domain.
3. Add the user with the `Person` user type, or your own user type if you created one.
4. Choose how the user gets a password:
   - Set a username and password yourself.
   - Enter the user's other email address, and ThunderID emails them an invitation link to set their own password.

The user signs in to mail clients with their full email address. Send them the [user guide](../user-guide/README.md).

## Use your own user type

The `Person` type fits most installs. To store different attributes:

1. In the console, create a user type with the fields you need, including `username` and `password`.
2. Assign it to the organization unit its users belong to.
3. Add users with it, as above.
4. If you use collaboration, list the type in `nextcloudDeps.registerJob.allowedUserTypes` and `tuwunelDeps.registerJob.allowedUserTypes`, and map its attributes in `nextcloudDeps.claims`. See [Turn on collaboration](enable-collaboration.md#user-types).

## Add an organization unit with its own domain

1. Create the organization unit in the console and note its ID.
2. Add the domain to your values file, for Raven to deliver it and OpenDKIM to sign it. Once you list any organization unit, list the `default` one too:

   ```yaml
   raven:
     organizationUnitDomains:
       "<default-ou-id>": ["example.com"]
       "<new-ou-id>": ["example.org"]

   opendkim:
     domains:
       - domain: example.com
         dkimSelector: mail
       - domain: example.org
         dkimSelector: mail
   ```

3. Run an upgrade. See [Upgrade OpenGovMail](upgrade.md).
4. Add DNS records for the new domain, as in the [installation guide](../installation-guide.md#step-1-add-dns-records), including its DKIM record from `/etc/dkimkeys/example.org/mail.txt`.

## See also

- [Identity and users](../concepts/identity-and-users.md)
- [Troubleshooting: mail clients can't sign in](../troubleshooting.md#mail-clients-cant-sign-in)
