# Identity and users

Every OpenGovMail account lives in [ThunderID](https://github.com/thunder-id/thunderid), an identity server. A user's account and their email address are the same thing: they sign in to mail, calendars, files and chat with one name and password, and an administrator manages them all in one place, the ThunderID console.

## Background

ThunderID organizes accounts with three ideas:

- An organization unit is a part of your organization, such as a department. Every user belongs to one.
- A user type, sometimes called a user schema, lists the attributes a kind of user has, such as a username, a display name and an email address.
- A group is a set of users. Groups cut across organization units.

The chart installs ThunderID with a `default` organization unit, a `Person` user type, an `Administrator` role and the console application. It also registers Raven as an application, so Raven can read organization units, users and groups. You don't create any of this by hand.

## How users map to mail

Raven decides which mail belongs to whom by asking ThunderID:

- Each organization unit owns one or more mail domains. The `default` organization unit owns `global.domain`. A domain belongs to one organization unit only.
- A user in an organization unit receives mail at that unit's domains, and signs in to mail clients with their full email address.
- An organization unit that owns no domain can't receive mail, and its users can't sign in to mail.
- A message to a group address goes to every member of the group.

To give other organization units their own domains, set `raven.organizationUnitDomains`, keyed by the unit's ID:

```yaml
raven:
  organizationUnitDomains:
    "<ou-id>": ["example.com", "example.org"]
```

## Collaboration sign-in

Nextcloud and Matrix register themselves in ThunderID on every install and upgrade, and send users to ThunderID to sign in. Only the user types they list can sign in: `Person` by default. Groups in Nextcloud come from ThunderID and update at each sign-in. [Turn on collaboration](../how-to/enable-collaboration.md#user-types) explains how to allow other user types.

## Use cases

- A single organization with one domain uses the defaults. Every user goes in the `default` organization unit.
- A ministry with several departments, each with its own domain, creates an organization unit per department and lists each one's domain in `raven.organizationUnitDomains`.
- An organization that already runs ThunderID sets `global.thunder.url` to it and `thunderid.enabled` to `false`. It must then set `raven.organizationUnitDomains`, and create Raven's application and permissions in its ThunderID.

## Comparison with a traditional mail server

| | Traditional mail server | OpenGovMail |
|---|---|---|
| Where accounts live | A mail-only user list or database | ThunderID, shared by every service |
| Mail address | Configured separately from the login | The user's account |
| Adding a person to a team | Edit mailing lists by hand | Add them to a ThunderID group |

## Related resources

- [Add users](../how-to/add-users.md)
- [Overview of the architecture](architecture.md)
- [ThunderID documentation](https://github.com/thunder-id/thunderid)
