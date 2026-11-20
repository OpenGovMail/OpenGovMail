# Connect a calendar app

## Overview

OpenGovMail calendars are in Nextcloud. Nextcloud supplies them through CalDAV. Thus, all CalDAV apps can use them. Nextcloud also supplies contacts through CardDAV. This procedure connects Thunderbird, Outlook, iPhones, iPads and Android phones.

In this procedure, `cloud.yourdomain.com` is the Nextcloud address of your organization. Your administrator can give you this address.

## Before you start: create an app password

You sign in to Nextcloud through ThunderID. Thus, you do not have a Nextcloud password for a calendar app. Create an app password for each device.

1. Sign in at `https://cloud.yourdomain.com`.
2. Open your avatar menu. Select `Personal settings`, then `Security`.
3. Below `Devices & sessions`, type a name, for example "Phone calendar".
4. Click `Create new app password`.
5. Copy the username and the password. Nextcloud shows the password only one time.

The username is a long string of letters and digits. It is not your email address. Calendar apps do not accept your email address here.

To remove the access of a device, delete its entry below `Devices & sessions`.

## Server settings

| Setting | Value |
|---|---|
| Server | `https://cloud.yourdomain.com` |
| CalDAV and CardDAV URL, if the app asks for a full path | `https://cloud.yourdomain.com/remote.php/dav/` |
| Username | The username from the app password screen |
| Password | The app password |

Most apps use only the server address. They find the other settings through `https://cloud.yourdomain.com/.well-known/caldav`.

## Thunderbird

Thunderbird 102 and later versions have CalDAV and CardDAV. You do not have to add add-ons.

1. Open the calendar. Click `New Calendar`.
2. Select `On the Network`. Click `Next`.
3. Enter the username. Enter `https://cloud.yourdomain.com` as the location.
4. Click `Find Calendars`.
5. When Thunderbird asks for the password, enter the app password.
6. Select the calendars to add. Click `Subscribe`.

To add contacts, do these steps:

1. Open the address book.
2. Select `New Address Book`, then `Add CardDAV Address Book`.
3. Enter the same username and location.

## Outlook

Outlook does not have CalDAV. On Windows, use the open source [Outlook CalDav Synchronizer](https://caldavsynchronizer.org/) add-in.

1. Install Outlook CalDav Synchronizer.
2. Create a profile with the type `Nextcloud`.
3. Enter `https://cloud.yourdomain.com/remote.php/dav/` as the DAV URL.
4. Enter the username and the app password.
5. Click `Test or discover settings`.
6. Select the calendar.

## iPhone and iPad

1. Open `Settings`, then `Calendar`, then `Accounts`. On iOS 18 and later versions, `Calendar` is in `Apps`.
2. Tap `Add Account`.
3. Tap `Other`, then `Add CalDAV Account`.
4. Enter `cloud.yourdomain.com` as the server.
5. Enter the username and the app password. Tap `Next`.

To add contacts, return to `Add Account`. Select `Other`, then `Add CardDAV Account`. Enter the same data.

## Android

Android must have a sync app. [DAVx5](https://www.davx5.com/) is an open source sync app.

1. In DAVx5, tap the `+` button.
2. Select `Login with URL and user name`.
3. Enter `https://cloud.yourdomain.com` as the base URL.
4. Enter the username and the app password.
5. Tap `Login`.
6. Select the calendars and address books to sync.

DAVx5 also has a Nextcloud login option. This option opens a browser. You sign in through ThunderID, and DAVx5 then creates the app password for you.

## Share calendars, find free times and send invitations

These functions are the same in all apps. Nextcloud Calendar on the web makes them easiest.

- To share a calendar, open its menu and select `Share`. Then type the name of a person or a group.
- The groups come from ThunderID. A change to a group in ThunderID applies when the person signs in again. At that time, Nextcloud removes group memberships that you added in Nextcloud.
- When you add attendees to an event, `Find a time` shows when they are busy. It shows only the busy periods. It does not show the titles or the details of events.
- Attendees get an email invitation. People outside your organization get `Accept` and `Decline` links. When they click a link, your event changes. They do not have to have an account.

## Known limits

- An outside attendee can reply from their own calendar app, not with `Accept` or `Decline`. Then the reply goes as an email to the mailbox of the organizer. Nextcloud does not read that mailbox. Thus, the event does not change.
- Only the ThunderID user types in `nextcloudDeps.registerJob.allowedUserTypes` can sign in to Nextcloud. Invitations must have an email address. Thus, these users must also have a value in the `email` attribute in ThunderID.

## See also

- [User guide](README.md)
- [Connect a files app](files-client.md)
- [Nextcloud Calendar user manual](https://docs.nextcloud.com/server/latest/user_manual/en/groupware/calendar.html)
