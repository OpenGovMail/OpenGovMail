# Turn on collaboration

## Overview

This guide turns on the collaboration services the Helm chart can run next to email. Users sign in to all of them through ThunderID with their mail account. Everything here is off by default.

| Part | What it gives users | Turn it on with |
|---|---|---|
| Nextcloud | Calendars (CalDAV), shared drives, files and webmail | `nextcloud.enabled: true` |
| Collabora | Document editing inside Nextcloud | On with Nextcloud (`nextcloud.collabora.enabled`) |
| Nextcloud Talk | Video meetings, with a TURN server so calls work behind NAT | On with Nextcloud (`nextcloudDeps.turn.enabled`) |
| Matrix (Tuwunel) | Chat in Element or any Matrix client | `tuwunel.enabled: true` |
| Matrix calls | Voice and video calls in Element | `tuwunel.matrixRtc.enabled: true` |

## Before you start

- OpenGovMail is installed and uses the `le-prod` issuer. Nextcloud won't fetch ThunderID's sign-in settings over an untrusted staging certificate.
- The server has at least 2 CPUs. On 1 CPU the pods' CPU requests fill the node, and some pods, including cert-manager's, never get scheduled.
- You have chosen a hostname for Nextcloud and one for Collabora, such as `cloud.example.com` and `office.example.com`.

## Turn on Nextcloud, Collabora and Talk

1. Add A records pointing at the server for the Nextcloud and Collabora hostnames.
2. Open 3478/udp, 3478/tcp and 49160-49200/udp on the node's firewall for Talk.
3. Copy the `nextcloud` block from [values-example.yaml](../../charts/opengovmail/values-example.yaml) into your values file and change every hostname to yours.
4. Run an upgrade. See [Upgrade OpenGovMail](upgrade.md).
5. Open `https://cloud.example.com` and sign in with a mail account.

Collabora's hostname is set in three places: `nextcloud.collabora.ingress`, `collabora.server_name` and `collabora.aliasgroups`. Users never open it directly; Nextcloud loads the editor from it.

## Create a shared drive

Shared drives are folders that belong to a ThunderID group rather than a person. Create one with Nextcloud's `occ` command:

```sh
kubectl exec -n opengovmail deploy/opengovmail-nextcloud -c nextcloud -- \
  su -s /bin/sh www-data -c "php occ groupfolders:create Finance"
# The command prints the new drive's ID. Give a group access to it:
kubectl exec -n opengovmail deploy/opengovmail-nextcloud -c nextcloud -- \
  su -s /bin/sh www-data -c "php occ groupfolders:group <id> <group-id> read write share delete"
```

`php occ group:list` shows each ThunderID group's ID next to its name.

## Turn on Matrix chat and calls

1. Add A records for `matrix.example.com`, and for `matrix-rtc.example.com` if you want calls.
2. For calls, open 7882/udp and 7881/tcp on the node's firewall.
3. Copy the `tuwunel` block from [values-example.yaml](../../charts/opengovmail/values-example.yaml) into your values file. Leave out `matrixRtc` if you don't want calls.
4. Run an upgrade.
5. Sign in to [Element](https://element.io/) with your mail account. `alice@example.com` becomes `@alice:example.com`.

Element finds the server through `https://example.com/.well-known/matrix/client`, which Tuwunel serves on your bare domain. If your bare domain is hosted somewhere else, such as a website, see [DNS in the Tuwunel chart](../../charts/opengovmail/charts/tuwunel/README.md#dns).

## Servers with a private IP address

If the node only has a private IP address behind NAT, set its public address so call media can reach it:

```yaml
nextcloudDeps:
  turn:
    externalIP: 12.34.56.78
tuwunel:
  matrixRtc:
    livekit:
      nodeIP: 12.34.56.78
```

The [Tuwunel chart's calls section](../../charts/opengovmail/charts/tuwunel/README.md#calls) has the details.

## User types

Nextcloud and Matrix register themselves as applications in ThunderID on every install and upgrade. Only the user types they list may sign in, and ThunderID rejects Nextcloud's registration if those types lack the attributes it asks for. The defaults match ThunderID's `Person` type. If your users have another type, set:

- `nextcloudDeps.registerJob.allowedUserTypes` and `tuwunelDeps.registerJob.allowedUserTypes` to its name.
- `nextcloudDeps.claims` to its attributes for the display name and email address, for example `{displayName: full_name, email: external_email}`.

Calendar invitations need an email address, so those users also need the `email` attribute set.

## Limits

Talk calls are peer to peer, which suits meetings of about four or five people. Larger meetings need Nextcloud's High Performance Backend, which the chart doesn't include yet.

## See also

- [Hostnames and ports](../reference/hostnames-and-ports.md)
- [Technology stack](../reference/technology-stack.md) for component versions
- [User guide](../user-guide/README.md) for connecting calendar and file apps
