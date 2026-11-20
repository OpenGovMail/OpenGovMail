# Documentation style guide

## Introduction

This guide is for anyone writing or changing OpenGovMail documentation. It records the choices this project has made. For everything else, follow the base guide below.

## Intended audience and scope

It covers the Markdown files in this repository: the README, everything under `docs/`, CONTRIBUTING.md and the chart READMEs. It doesn't cover code comments or commit messages; see [CONTRIBUTING.md](../CONTRIBUTING.md) for those.

## Our preferred style guide

The [Google developer documentation style guide](https://developers.google.com/style). Where this page disagrees with it, this page wins.

## Glossary of preferred terms

| Use | Not |
|---|---|
| ThunderID | Thunder, Thunder ID, Asgardeo Thunder |
| organization unit | OU in running text, organisation unit |
| user type | user schema, except where the ThunderID console uses that name |
| mail app | MUA, mail user agent, email client |
| Helm chart, the chart | Helm charts, when meaning OpenGovMail's one chart |
| `example.com` | `yourdomain.com`, a real domain |
| `12.34.56.78` | a real IP address |

Define new terms in the [glossary](glossary.md).

## Topic types and templates

Each page follows one template from [The Good Docs Project](https://gitlab.com/tgdp/templates). Pick the type before you write.

| Type | Template | Where | Example |
|---|---|---|---|
| Project overview | README | `README.md` | [README.md](../README.md) |
| Quickstart | Quickstart | `docs/quickstart.md` | [Quickstart](quickstart.md) |
| Installation guide | Installation guide | `docs/installation-guide.md` | [Installation guide](installation-guide.md) |
| Explanation of an idea | Concept | `docs/concepts/` | [Identity and users](concepts/identity-and-users.md) |
| Steps for one task | How-to | `docs/how-to/`, or `docs/user-guide/` for end users | [Add users](how-to/add-users.md) |
| Facts to look up | Reference | `docs/reference/` | [Helm values](reference/helm-values.md) |
| Problems and fixes | Troubleshooting | `docs/troubleshooting.md` | [Troubleshooting](troubleshooting.md) |
| Terms | Glossary | `docs/glossary.md` | [Glossary](glossary.md) |
| Changes per release | Changelog | `CHANGELOG.md` | [CHANGELOG.md](../CHANGELOG.md) |
| How to contribute | Contributing guide | `CONTRIBUTING.md` | [CONTRIBUTING.md](../CONTRIBUTING.md) |

Add the new page to the [documentation index](README.md).

## General writing tips

- Helm on Kubernetes is the only supported way to install OpenGovMail. Don't document other install methods.
- Write in the second person and the present tense.
- Use sentence case for headings.
- Lead with what the reader does. Put the reason after the step, in one sentence.
- One task per how-to, with no more than about 10 steps.
- Put commands in fenced code blocks that can be copied and run as they are. Mark the parts to replace with `<angle brackets>`.
- Use bold sparingly, and don't use em dashes.
- When a page states a version, say which chart release it's for.

## Accessible writing

- Give every image alt text that says what it shows.
- Make link text say where it goes. Not "click here".
- Don't rely on colour alone in diagrams.

## Inclusive and bias-free writing

Follow Google's [inclusive documentation](https://developers.google.com/style/inclusive-documentation) guidance. Use "they" for a person whose pronouns you don't know.

## How the style guide is updated

Propose a change in a pull request that edits this page, and record it in the decision log below.

## Decision log

| Date | Decision |
|---|---|
| 2026-11-19 | Adopted The Good Docs Project templates and this guide. Helm is the only documented install method. |
