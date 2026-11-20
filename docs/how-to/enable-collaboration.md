# Add the collaboration services

## Overview

This procedure adds the collaboration services to OpenGovMail. The Helm chart can operate these services with email. Users sign in to all of them through ThunderID with their mail account. By default, all the collaboration services are off.

| Service | What users get | Setting |
|---|---|---|
| Nextcloud | Calendars (CalDAV), shared drives, files and webmail | `nextcloud.enabled: true` |
| Collabora | Document editor in Nextcloud | Starts with Nextcloud (`nextcloud.collabora.enabled`) |
| Nextcloud Talk | Video meetings. A TURN server lets calls go through NAT. | Starts with Nextcloud (`nextcloudDeps.turn.enabled`) |
| Matrix (Tuwunel) | Chat in Element or in a different Matrix app | `tuwunel.enabled: true` |
| Matrix calls | Voice and video calls in Element | `tuwunel.matrixRtc.enabled: true` |

## Before you start

- Install OpenGovMail with the `le-prod` issuer. Nextcloud does not accept the sign-in settings of ThunderID with a staging certificate.
- Make sure that your values file has the `thunderid` certificate settings from [Step 5 of the installation guide](../installation-guide.md#step-5-write-a-values-file). These settings are `thunderid.configuration.tls` and `thunderid.declarativeResources`.

  If these settings are not in the file, ThunderID uses a self-signed certificate. Then the upgrade stops with an error. The error tells you that Nextcloud and Tuwunel must have a trusted certificate. An installation with only `--set` flags and no values file does not have these settings.

- Make sure that the server has a minimum of 2 CPUs. With 1 CPU, Kubernetes cannot start some pods, for example the cert-manager pods.
- Choose a hostname for Nextcloud and a hostname for Collabora, for example `cloud.example.com` and `office.example.com`.

## Add Nextcloud, Collabora and Talk

1. Add an A record for the Nextcloud hostname and for the Collabora hostname. Point the records to the server.
2. For Talk, open 3478/udp, 3478/tcp and 49160-49200/udp on the firewall of the node.
3. Copy the `nextcloud` block from [values-example.yaml](../../charts/opengovmail/values-example.yaml) into your values file.
4. Change each hostname in the block to your hostname.
5. Do an upgrade. Read [Upgrade OpenGovMail](upgrade.md).
6. Open `https://cloud.example.com`. Sign in with a mail account.

The Collabora hostname is in three settings: `nextcloud.collabora.ingress`, `nextcloud.collabora.collabora.server_name` and `nextcloud.collabora.collabora.aliasgroups`. Users do not open Collabora directly. Nextcloud gets the editor from it.

## Create a shared drive

A shared drive is a folder for a ThunderID group, not for one person.

1. Create the drive with the Nextcloud `occ` command:

   ```sh
   kubectl exec -n opengovmail deploy/opengovmail-nextcloud -c nextcloud -- \
     su -s /bin/sh www-data -c "php occ groupfolders:create Finance"
   ```

   The command shows the ID of the new drive.

2. Find the ID of the group. `php occ group:list` shows the ID of each ThunderID group, next to its name.
3. Give the group access to the drive:

   ```sh
   kubectl exec -n opengovmail deploy/opengovmail-nextcloud -c nextcloud -- \
     su -s /bin/sh www-data -c "php occ groupfolders:group <id> <group-id> read write share delete"
   ```

## Add Matrix chat and calls

1. Add an A record for `matrix.example.com`.
2. For calls, do these steps:
   1. Add an A record for `matrix-rtc.example.com`.
   2. Open 7882/udp and 7881/tcp on the firewall of the node.
3. Copy the `tuwunel` block from [values-example.yaml](../../charts/opengovmail/values-example.yaml) into your values file. If you do not want calls, remove `matrixRtc`.
4. Do an upgrade.
5. Sign in to [Element](https://element.io/) with your mail account. The Matrix ID for `alice@example.com` is `@alice:example.com`.

Element finds the server through `https://example.com/.well-known/matrix/client`. Tuwunel supplies this address on your domain. If a different server has your domain, for example a website, read [DNS in the Tuwunel chart](../../charts/opengovmail/charts/tuwunel/README.md#dns).

## Servers with a private IP address

If the node has only a private IP address behind NAT, set its public address. Call media can then get to the node:

```yaml
nextcloudDeps:
  turn:
    externalIP: 12.34.56.78
tuwunel:
  matrixRtc:
    livekit:
      nodeIP: 12.34.56.78
```

For more data, read the [calls section of the Tuwunel chart](../../charts/opengovmail/charts/tuwunel/README.md#calls).

## User types

Nextcloud and Matrix register as applications in ThunderID at each installation and upgrade. Only the user types in their lists can sign in. ThunderID does not accept the Nextcloud registration if these user types do not have the attributes that Nextcloud asks for. The default values are for the ThunderID `Person` user type.

If your users have a different user type, set these values:

- Set `nextcloudDeps.registerJob.allowedUserTypes` and `tuwunelDeps.registerJob.allowedUserTypes` to the name of the user type.
- Set `nextcloudDeps.claims` to the attributes for the display name and the email address. For example: `{displayName: full_name, email: external_email}`.

Calendar invitations must have an email address. Thus, these users must also have a value in the `email` attribute.

## Limits

Talk calls go directly between the people in the call. This is satisfactory for meetings of approximately four or five people. Larger meetings must have the Nextcloud High Performance Backend. The chart does not have it yet.

## See also

- [Hostnames and ports](../reference/hostnames-and-ports.md)
- [Technology stack](../reference/technology-stack.md) for the versions of the components
- [User guide](../user-guide/README.md) for calendar and files apps
