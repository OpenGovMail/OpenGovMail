# Set up a files app

## Overview

OpenGovMail files run in Nextcloud. You can use them in the browser, sync them to a computer, or open them on a phone.

In this guide, `cloud.yourdomain.com` stands for your organization's Nextcloud address. Your administrator can tell you what it is.

## In the browser

Sign in at `https://cloud.yourdomain.com` through ThunderID.

- Your group's shared drives appear as folders in `Files`. An administrator creates each one and gives a ThunderID group access.
- Click a document, spreadsheet or presentation to edit it in the browser. Others who open the same file edit it with you at the same time.
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

## See also

- [User guide](README.md)
- [Nextcloud Files user manual](https://docs.nextcloud.com/server/latest/user_manual/en/files/index.html)
