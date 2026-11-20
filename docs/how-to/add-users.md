# Add users

## Overview

This procedure adds mail users in the ThunderID console. The ThunderID console is a part of OpenGovMail. Each user that you add gets a mailbox. The user can sign in to all the OpenGovMail services. [Identity and users](../concepts/identity-and-users.md) tells you about organization units, user types and groups.

## Before you start

- Install OpenGovMail. Read the [installation guide](../installation-guide.md).
- Make sure that you can open the console at `https://<your-domain>:8090/console`.
- Get the password of the ThunderID `admin` user. If the chart made the password, show it with this command:

  ```bash
  kubectl get secret opengovmail-thunder-admin -n opengovmail \
    -o jsonpath='{.data.password}' | base64 -d
  ```

## Add a user

1. Open `https://<your-domain>:8090/console`.
2. Sign in as `admin`.
3. Select the organization unit of the user. The `default` organization unit has your main domain.
4. Add the user with the `Person` user type. If you created a different user type, you can use it.
5. Give the user a password. Do one of these:
   - Set a username and a password.
   - Enter a different email address of the user. ThunderID sends an invitation to this address. The user sets a password with the link in the invitation.
6. Send the [user guide](../user-guide/README.md) to the user.

The user signs in to mail apps with the full email address.

## Use a different user type

The `Person` user type is correct for most installations. To keep different attributes, do these steps:

1. In the console, create a user type with the fields that you must have. Include `username` and `password`.
2. Assign the user type to the organization unit of its users.
3. Add users with the user type. Use the procedure in [Add a user](#add-a-user).
4. If you use the collaboration services, do these steps:
   1. Add the user type to `nextcloudDeps.registerJob.allowedUserTypes` and `tuwunelDeps.registerJob.allowedUserTypes`.
   2. Write the attributes of the user type in `nextcloudDeps.claims`.

   Read [Add the collaboration services](enable-collaboration.md#user-types).

## Add an organization unit with a different domain

1. Create the organization unit in the console. Write down its ID.
2. Add the domain to your values file. Raven then delivers mail for the domain, and OpenDKIM signs mail from the domain:

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

   If you add one organization unit to `raven.organizationUnitDomains`, you must also add the `default` organization unit.

3. Do an upgrade. Read [Upgrade OpenGovMail](upgrade.md).
4. Add the DNS records for the new domain. Use the procedure in the [installation guide](../installation-guide.md#step-1-add-dns-records).
5. Add the DKIM record for the new domain. The record is in `/etc/dkimkeys/example.org/mail.txt`.

## See also

- [Identity and users](../concepts/identity-and-users.md)
- [Troubleshooting: Mail apps cannot sign in](../troubleshooting.md#mail-apps-cannot-sign-in)
