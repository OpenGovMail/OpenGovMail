# Identity and users

All OpenGovMail accounts are in [ThunderID](https://github.com/thunder-id/thunderid). ThunderID is an identity server. The account of a user and the email address of the user are the same. Users sign in to mail, calendars, files and chat with one name and one password. An administrator manages all the accounts in one location, the ThunderID console.

## Background

ThunderID uses three items to organize accounts:

- An organization unit is a part of your organization, for example a department. Each user is in one organization unit.
- A user type gives the attributes of a type of user, for example a username, a display name and an email address. Some screens call it a user schema.
- A group is a set of users. The users in a group can be in different organization units.

The chart installs ThunderID with a `default` organization unit, a `Person` user type, an `Administrator` role and the console application. It also registers Raven as an application. Thus, Raven can read the organization units, users and groups. You do not create these items yourself.

## How users connect to mail

Raven gets the data for mail delivery from ThunderID:

- Each organization unit has one or more mail domains. The `default` organization unit has `global.domain`. Only one organization unit can have a given domain.
- A user receives mail at the domains of the organization unit of the user. The user signs in to mail apps with the full email address.
- An organization unit without a domain cannot receive mail. Its users cannot sign in to mail.
- A message to a group address goes to all the members of the group.

To give domains to other organization units, set `raven.organizationUnitDomains`. Use the ID of the organization unit as the key:

```yaml
raven:
  organizationUnitDomains:
    "<ou-id>": ["example.com", "example.org"]
```

## Collaboration sign-in

Nextcloud and Matrix register in ThunderID at each installation and upgrade. They send users to ThunderID to sign in. Only the user types in their lists can sign in. The default list has only `Person`. Nextcloud gets its groups from ThunderID, and updates them at each sign-in. [Add the collaboration services](../how-to/enable-collaboration.md#user-types) tells you how to add other user types.

## Use cases

- An organization with one domain uses the default values. All users go in the `default` organization unit.
- A ministry has some departments, and each department has its own domain. The ministry creates one organization unit for each department. It puts the domain of each department in `raven.organizationUnitDomains`.
- An organization that has a ThunderID already sets `global.thunder.url` to it, and sets `thunderid.enabled` to `false`. It must then do these steps:
  - Set `raven.organizationUnitDomains`.
  - Create the application and the permissions for Raven in its ThunderID.
  - Set `raven.thunder.resourceIdentifier` to the identifier of the System resource server in its ThunderID, for example `https://thunder.example.com/mcp`.

  The default identifier, `https://localhost:8090/mcp`, is correct only for the ThunderID that comes with the chart. If the identifier is not correct, ThunderID refuses the token request from Raven with `invalid_target`. Then users cannot sign in to mail.

## Comparison with a usual mail server

| | Usual mail server | OpenGovMail |
|---|---|---|
| Location of the accounts | A list or a database for mail only | ThunderID, for all the services |
| Mail address | Separate from the login | The same as the account |
| Add a person to a team | Change the mailing lists manually | Add the person to a ThunderID group |

## Related resources

- [Add users](../how-to/add-users.md)
- [Overview of the architecture](architecture.md)
- [ThunderID documentation](https://github.com/thunder-id/thunderid)
