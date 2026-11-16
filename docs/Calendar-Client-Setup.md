# Calendar client setup

OpenGovMail calendars run in Nextcloud and are served over CalDAV, so any CalDAV client can use them. Contacts are served over CardDAV from the same server.

In this guide, `cloud.yourdomain.com` stands for your Nextcloud host, the value of `nextcloud.nextcloud.host` in the Helm chart.

## Before you start: create an app password

You sign in to Nextcloud through ThunderID, so you have no Nextcloud password to give a calendar app. Create an app password for each device instead.

1. Sign in at `https://cloud.yourdomain.com`.
2. Open your avatar menu and choose `Personal settings`, then `Security`.
3. Under `Devices & sessions`, type a name such as "Phone calendar" and click `Create new app password`.
4. Copy both the username and the password it shows. The password is shown only once.

The username is a long string of letters and digits, not your email address. Calendar apps reject your email address here. To revoke a device later, delete its entry under `Devices & sessions`.

## Server settings

| Setting | Value |
|---|---|
| Server | `https://cloud.yourdomain.com` |
| CalDAV and CardDAV URL, if the app asks for a full path | `https://cloud.yourdomain.com/remote.php/dav/` |
| Username | the username from the app password screen |
| Password | the app password |

Most apps only need the server address. They find the rest through `https://cloud.yourdomain.com/.well-known/caldav`.

## Thunderbird

Thunderbird 102 and later support CalDAV and CardDAV without add-ons.

1. Open the calendar, then click `New Calendar`.
2. Choose `On the Network` and click `Next`.
3. Enter the username, and `https://cloud.yourdomain.com` as the location. Click `Find Calendars`.
4. Enter the app password when asked, select the calendars to add, and click `Subscribe`.

For contacts, open the address book, choose `New Address Book` and then `Add CardDAV Address Book`, and use the same username and location.

## Outlook

Outlook has no built-in CalDAV support. On Windows, install the open source [Outlook CalDav Synchronizer](https://caldavsynchronizer.org/) add-in. Create a profile of type `Nextcloud`, enter `https://cloud.yourdomain.com/remote.php/dav/` as the DAV URL, then the username and app password. Click `Test or discover settings` and pick the calendar.

## iPhone and iPad

1. Open `Settings`, then `Calendar`, then `Accounts`, and tap `Add Account`. On iOS 18 and later, `Calendar` is under `Apps`.
2. Tap `Other`, then `Add CalDAV Account`.
3. Enter `cloud.yourdomain.com` as the server, then the username and app password. Tap `Next`.

For contacts, go back to `Add Account`, choose `Other`, then `Add CardDAV Account`, and enter the same details.

## Android

Android needs a sync app. [DAVx5](https://www.davx5.com/) is open source.

1. In DAVx5, tap the `+` button.
2. Choose `Login with URL and user name`.
3. Enter `https://cloud.yourdomain.com` as the base URL, then the username and app password.
4. Tap `Login`, then choose which calendars and address books to sync.

DAVx5 also offers a Nextcloud login option. It opens a browser where you sign in through ThunderID, then creates the app password for you.

## Sharing, availability and invitations

These work the same in every client, but are easiest in Nextcloud Calendar on the web.

- To share a calendar, open its menu, choose `Share`, and type a person's or group's name. Groups come from ThunderID, and membership changes there apply at each person's next sign-in. Group memberships added in Nextcloud itself are removed at that point.
- When you add attendees to an event, `Find a time` shows when they are busy. It shows busy periods only, not event titles or details.
- Attendees get an email invitation. People outside your organisation get `Accept` and `Decline` links that update your event when clicked, without an account.

## Known limits

- An outside attendee who replies from their own calendar app, instead of clicking `Accept` or `Decline`, sends the reply as email to the organizer's mailbox. Nextcloud does not read that mailbox, so the event does not update.
- Only ThunderID user types listed in `nextcloudDeps.registerJob.allowedUserTypes` can sign in to Nextcloud. Invitations need an email address, so those users also need the `email` attribute set in ThunderID.
