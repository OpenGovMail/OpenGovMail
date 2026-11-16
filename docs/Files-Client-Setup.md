# Files client setup

OpenGovMail files run in Nextcloud. You can use them in the browser, sync them to a computer, or open them on a phone.

In this guide, `cloud.yourdomain.com` stands for your Nextcloud host, the value of `nextcloud.nextcloud.host` in the Helm chart.

## In the browser

Sign in at `https://cloud.yourdomain.com` through ThunderID.

- Your group's shared drives appear as folders in `Files`. An administrator creates each one and gives a ThunderID group access.
- To share a file outside the organisation, open its `Share` panel and create a link. Set an expiry date, and a password if the file needs one.
- An overwritten file keeps its earlier versions under `Versions`. A deleted file goes to `Deleted files`, where you can restore it.

## Desktop sync

1. Install the [Nextcloud desktop client](https://nextcloud.com/install/#desktop-files) for Windows, macOS or Linux.
2. Enter `https://cloud.yourdomain.com` as the server address.
3. A browser window opens. Sign in through ThunderID and grant access.
4. Choose a local folder and what to sync.

## Phones and tablets

1. Install the Nextcloud app from the App Store, Google Play or F-Droid.
2. Enter `https://cloud.yourdomain.com` and sign in through ThunderID in the window it opens.

## Other WebDAV apps

Apps that only speak WebDAV need an app password. Create one under `Personal settings`, then `Security`, then `Devices & sessions`, and use the username and password it shows with this address:

```
https://cloud.yourdomain.com/remote.php/dav/files/<username>/
```

## For administrators: creating a shared drive

```sh
kubectl exec -n <namespace> deploy/<release>-nextcloud -c nextcloud -- \
  su -s /bin/sh www-data -c "php occ groupfolders:create Finance"
# The command prints the new drive's ID. Give a group access to it:
kubectl exec -n <namespace> deploy/<release>-nextcloud -c nextcloud -- \
  su -s /bin/sh www-data -c "php occ groupfolders:group <id> <group-id> read write share delete"
```

`php occ group:list` shows each ThunderID group's ID next to its name.
