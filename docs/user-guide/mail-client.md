# Connect a mail app

## Overview

This procedure connects a mail app to your OpenGovMail mailbox. It gives the steps for Thunderbird, Apple Mail, Outlook and the Gmail app.

## Before you start

You must have:

- Your full email address and your password. Your administrator gives you these.
- The hostname of your mail server. In this procedure, the hostname is `mail.example.com`. Replace `example.com` with the domain of your organization.

## Server settings

All apps use the same settings. Most apps find the settings automatically. If your app asks for them, use these settings:

| Setting | Value |
|---|---|
| Incoming server (IMAP) | `mail.example.com` |
| Incoming port and security | 993, SSL/TLS |
| Outgoing server (SMTP) | `mail.example.com` |
| Outgoing port and security | 587, STARTTLS |
| Username | Your full email address |
| Authentication | Normal password |

Port 993 encrypts the connection from the start. If an old app cannot use port 993, use port 143 with STARTTLS.

## Thunderbird

1. Install [Thunderbird](https://www.thunderbird.net/).
2. Open `Account Settings`.
3. Click `New Account`, then `Mail Account`.
4. Enter your name and your email address. Click `Continue`.
5. If Thunderbird finds the settings, click `Continue`. Then enter your password.
6. If Thunderbird does not find the settings, click `Manual config`. Enter the [server settings](#server-settings). Click `Re-test`.
7. Click `Continue`.

## Apple Mail (macOS)

1. Open Mail. Select `Mail`, then `Add Account`. On macOS Ventura and later versions, use `Accounts` in System Settings.
2. Select `Other Mail Account`. Click `Continue`.
3. Enter your name, your email address and your password. Click `Sign In`.
4. If Mail does not find the settings, enter the [server settings](#server-settings). Use your full email address as the username for the two servers.
5. Click `Sign In`.

## Outlook (Windows and macOS)

1. Open Outlook. Select `File`, then `Add Account`.
2. Enter your email address. Click `Connect`.
3. When Outlook asks for the account type, select `IMAP`.
4. Enter your password and the [server settings](#server-settings).
5. Click `Next`, then `Finish`.

## Gmail app on iPhone and iPad

1. Open the Gmail app. Tap your profile picture, then `Add another account`.
2. Tap `Other`. Enter your email address. Tap `Next`.
3. Select `IMAP`. Enter your password. Tap `Next`.
4. Enter the [server settings](#server-settings) for the incoming server and the outgoing server.
5. Tap `Next`.

## Gmail app on Android

1. Open the Gmail app. Tap the menu, then `Settings`, then `Add account`.
2. Tap `Other`. Enter your email address. Tap `Manual setup`.
3. Select `Personal (IMAP)`. Enter your password. Tap `Next`.
4. Enter the [server settings](#server-settings) for the incoming server and the outgoing server.
5. Tap `Next`.

## If the app does not connect

- Make sure that your email address and your password are correct.
- Make sure that your network does not close ports 587 and 993. Some public and office networks close these ports.
- If the app shows a warning about the certificate of the server, do not accept the certificate. Tell your administrator. Your administrator can read [Troubleshooting](../troubleshooting.md#mail-apps-cannot-sign-in).

## See also

- [User guide](README.md)
- [Connect a calendar app](calendar-client.md)
