# User guide

This guide is for people who use an OpenGovMail account: reading and sending mail, and, where your organization has turned them on, calendars, files, meetings and chat. If you run the server, start with the [installation guide](../installation-guide.md) instead.

## Your account

Your administrator creates your account and gives you an email address and a password. You use that one account everywhere: in mail apps, in Nextcloud and in Element.

## Set up your apps

| To use | Read |
|---|---|
| Mail in Thunderbird, Apple Mail, Outlook or the Gmail app | [Set up a mail app](mail-client.md) |
| Calendars and contacts | [Set up a calendar app](calendar-client.md) |
| Files and shared drives | [Set up a files app](files-client.md) |
| Video meetings | [Nextcloud Talk user manual](https://docs.nextcloud.com/server/latest/user_manual/en/talk/index.html) |
| Documents, spreadsheets and presentations | [Collabora Online help](https://help.collaboraoffice.com/) |
| Chat | [Element user guide](https://element.io/user-guide) |

Calendars, files, meetings, documents and chat are only there if your organization has turned them on. Ask your administrator.

## FAQ

### What username do I sign in with?

Your full email address, such as `alice@example.com`, in mail apps. Nextcloud and Element send you to a ThunderID sign-in page where you use the same address and password.

### Why does my calendar app reject my email address and password?

Calendar and contact apps can't use the ThunderID sign-in page, so they need an app password from Nextcloud instead. [Set up a calendar app](calendar-client.md#before-you-start-create-an-app-password) shows how to create one.

### What is my Matrix ID in Element?

Your email address with the `@` moved: `alice@example.com` becomes `@alice:example.com`.

### Which ports does my mail app need?

993 for reading mail and 587 for sending. Some public and office networks block them. [Set up a mail app](mail-client.md#server-settings) lists every setting.

### My app warns that the server's certificate isn't trusted. Should I accept it?

No. Tell your administrator. A correctly installed server has a trusted certificate.

### How many people can join a video meeting?

Nextcloud Talk suits meetings of about four or five people. Larger meetings may be unreliable.

### Can I get a deleted file back?

Yes. Deleted files go to `Deleted files` in Nextcloud, and overwritten files keep their earlier versions under `Versions`.

### How do I change my password?

Ask your administrator. Password changes are made in ThunderID.

### Who can read my mail?

Your mail is stored on your organization's own servers. The OpenGovMail project has no access to it.
