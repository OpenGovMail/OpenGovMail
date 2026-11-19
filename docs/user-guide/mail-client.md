# Set up a mail app

## Overview

This guide connects a mail app, such as Thunderbird, Apple Mail, Outlook or the Gmail app, to your OpenGovMail mailbox.

## Before you start

You need:

- Your full email address and password. Your administrator gives you these.
- Your mail server's hostname. In this guide it's `mail.example.com`; replace `example.com` with your organization's domain.

## Server settings

Every app needs the same settings. Most find them by themselves; if yours asks, use these:

| Setting | Value |
|---|---|
| Incoming server (IMAP) | `mail.example.com` |
| Incoming port and security | 993, SSL/TLS |
| Outgoing server (SMTP) | `mail.example.com` |
| Outgoing port and security | 587, STARTTLS |
| Username | Your full email address |
| Authentication | Normal password |

Port 993 encrypts the connection from the start. If an older app can't use it, port 143 with STARTTLS also works.

## Thunderbird

1. Install [Thunderbird](https://www.thunderbird.net/) and open `Account Settings`.
2. Click `New Account`, then `Mail Account`.
3. Enter your name and email address, and click `Continue`.
4. If Thunderbird finds the settings, click `Continue` and enter your password.
5. If it doesn't, click `Manual config`, enter the [server settings](#server-settings), and click `Re-test`.
6. Click `Continue` to finish.

## Apple Mail (macOS)

1. Open Mail and choose `Mail`, then `Add Account`. On macOS Ventura and later, use `Accounts` in System Settings.
2. Choose `Other Mail Account` and click `Continue`.
3. Enter your name, email address and password, and click `Sign In`.
4. If Mail can't find the settings, enter the [server settings](#server-settings), with your full email address as the username for both servers.
5. Click `Sign In`.

## Outlook (Windows and macOS)

1. Open Outlook and choose `File`, then `Add Account`.
2. Enter your email address and click `Connect`.
3. Choose `IMAP` when asked.
4. Enter your password and the [server settings](#server-settings).
5. Click `Next`, then `Finish`.

## Gmail app on iPhone and iPad

1. Open the Gmail app, tap your profile picture, then `Add another account`.
2. Tap `Other`, enter your email address, and tap `Next`.
3. Choose `IMAP`, enter your password, and tap `Next`.
4. Enter the [server settings](#server-settings) for the incoming and outgoing servers.
5. Tap `Next` to finish.

## Gmail app on Android

1. Open the Gmail app, tap the menu, then `Settings` and `Add account`.
2. Tap `Other`, enter your email address, and tap `Manual setup`.
3. Choose `Personal (IMAP)`, enter your password, and tap `Next`.
4. Enter the [server settings](#server-settings) for the incoming and outgoing servers.
5. Tap `Next` to finish.

## If it doesn't connect

- Check your email address and password.
- Make sure your network doesn't block ports 587 and 993. Some public and office networks do.
- If the app warns about the server's certificate, don't accept it. Tell your administrator, who can check [Troubleshooting](../troubleshooting.md#mail-clients-cant-sign-in).

## See also

- [User guide](README.md)
- [Set up a calendar app](calendar-client.md)
