# Connect a files app

## Overview

OpenGovMail files are in Nextcloud. You can use them in a browser, sync them to a computer, or open them on a phone.

In this procedure, `cloud.yourdomain.com` is the Nextcloud address of your organization. Your administrator can give you this address.

## In the browser

Sign in at `https://cloud.yourdomain.com` through ThunderID.

- The shared drives of your group are folders in `Files`. An administrator creates each shared drive and gives a ThunderID group access to it.
- To edit a document, a spreadsheet or a presentation in the browser, click it. Other people can open the same file and edit it at the same time.
- To share a file with people outside your organization, open its `Share` panel and create a link. Set an expiry date. If the file is sensitive, also set a password.
- When you replace a file, Nextcloud keeps its earlier versions in `Versions`. When you delete a file, Nextcloud moves it to `Deleted files`. You can restore it from there.

## Desktop sync

1. Install the [Nextcloud desktop client](https://nextcloud.com/install/#desktop-files) for Windows, macOS or Linux.
2. Enter `https://cloud.yourdomain.com` as the server address.
3. A browser window opens. Sign in through ThunderID and give access.
4. Select a local folder. Select the folders to sync.

## Phones and tablets

1. Install the Nextcloud app from the App Store, Google Play or F-Droid.
2. Enter `https://cloud.yourdomain.com`.
3. Sign in through ThunderID in the window that opens.

## Other WebDAV apps

Apps that use only WebDAV must have an app password.

1. In Nextcloud, open `Personal settings`, then `Security`, then `Devices & sessions`.
2. Create an app password.
3. In the app, use the username and the password that Nextcloud shows, with this address:

   ```
   https://cloud.yourdomain.com/remote.php/dav/files/<username>/
   ```

## See also

- [User guide](README.md)
- [Nextcloud Files user manual](https://docs.nextcloud.com/server/latest/user_manual/en/files/index.html)
