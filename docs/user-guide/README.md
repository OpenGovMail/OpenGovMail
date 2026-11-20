# User guide

This guide is for people who use an OpenGovMail account. It tells you how to read and send mail. If your organization added them, it also tells you about calendars, files, meetings and chat. If you operate the server, read the [installation guide](../installation-guide.md).

## Your account

Your administrator creates your account. The administrator gives you an email address and a password. You use this one account in all the apps: mail apps, Nextcloud and Element.

## Connect your apps

| For | Read |
|---|---|
| Mail in Thunderbird, Apple Mail, Outlook or the Gmail app | [Connect a mail app](mail-client.md) |
| Calendars and contacts | [Connect a calendar app](calendar-client.md) |
| Files and shared drives | [Connect a files app](files-client.md) |
| Video meetings | [Nextcloud Talk user manual](https://docs.nextcloud.com/server/latest/user_manual/en/talk/index.html) |
| Documents, spreadsheets and presentations | [Collabora Online help](https://help.collaboraoffice.com/) |
| Chat | [Element user guide](https://element.io/user-guide) |

Calendars, files, meetings, documents and chat are available only if your organization added them. Ask your administrator.

## FAQ

### Which username do I use to sign in?

In mail apps, use your full email address, for example `alice@example.com`. Nextcloud and Element open a ThunderID sign-in page. Use the same address and password there.

### Why does my calendar app not accept my email address and password?

Calendar and contacts apps cannot use the ThunderID sign-in page. They must have an app password from Nextcloud. [Connect a calendar app](calendar-client.md#before-you-start-create-an-app-password) tells you how to create one.

### What is my Matrix ID in Element?

Your Matrix ID is your email address in a different form. For `alice@example.com`, the Matrix ID is `@alice:example.com`.

### Which ports does my mail app use?

Port 993 to read mail and port 587 to send mail. Some public and office networks close these ports. [Connect a mail app](mail-client.md#server-settings) gives all the settings.

### My app tells me that it does not trust the certificate of the server. Do I accept the certificate?

No. Tell your administrator. A correct installation has a trusted certificate.

### How many people can be in a video meeting?

Nextcloud Talk is satisfactory for meetings of approximately four or five people. Larger meetings possibly do not operate correctly.

### Can I get a deleted file again?

Yes. Nextcloud moves deleted files to `Deleted files`. When you replace a file, Nextcloud keeps its earlier versions in `Versions`.

### How do I change my password?

Ask your administrator. The administrator changes passwords in ThunderID.

### Who can read my mail?

Your mail is on the servers of your organization. The OpenGovMail project does not have access to it.
