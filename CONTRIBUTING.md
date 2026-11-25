# Contributing to OpenGovMail

## Welcome

Thank you for wanting to contribute to OpenGovMail. Contributions of any size are welcome: bug reports, documentation fixes, chart changes and new features.

## OpenGovMail overview

OpenGovMail is a self-hosted email and collaboration platform, installed with one Helm chart on Kubernetes. [Overview of the architecture](docs/concepts/architecture.md) explains the parts. This repository holds:

| Path | What it holds |
|---|---|
| `charts/opengovmail/` | The Helm chart and its subcharts |
| `services/smtp/`, `services/dkim/` | Dockerfiles for the Postfix and OpenDKIM images |
| `infra/` | Cluster setup applied once, such as the certificate issuers |
| `docs/` | Documentation |

Raven, the mail store, lives in its own repository, [OpenGovMail/raven](https://github.com/OpenGovMail/raven).

## Ground rules

- Follow the [code of conduct](CODE_OF_CONDUCT.md).
- Report security problems privately, as [SECURITY.md](SECURITY.md) describes, not in a public issue.
- Discuss large changes in an issue before writing the code.

## Community engagement

- Ask questions and talk to the maintainers in the [OpenGovMail Slack](https://join.slack.com/t/lsflk/shared_invite/zt-4ddz9zqnn-8lnG12b9Ty1Pclnx_hdIwA).
- Track work in [GitHub issues](https://github.com/OpenGovMail/OpenGovMail/issues).

## Share ideas

Open an issue with the feature template. Say what problem it solves and who it's for.

## Before you start

- Look for an existing issue. The [good first issues](https://github.com/OpenGovMail/OpenGovMail/issues?q=is%3Aissue%20state%3Aopen%20label%3A%22good%20first%20issue%22) are a good place to begin.
- If nobody is assigned, comment on the issue to say you're working on it.

## Environment setup

You need [Git](https://git-scm.com/), [Helm](https://helm.sh/docs/intro/install/) and [Docker](https://docs.docker.com/engine/install/). To run your changes end to end you also need an amd64 test server with OpenGovMail installed.

[Run your own build](docs/how-to/run-local-build.md) shows how to check chart changes on your own machine and run them on a test server.

### Coding agents

If you use an AI coding agent, such as Claude Code, Cursor, GitHub Copilot or Codex, it reads the project instructions in [AGENTS.md](AGENTS.md) and the skills in `.agents/skills/`. You review and answer for everything it writes, as for your own code.

On Windows, Claude Code finds the skills through the `.claude/skills` symlink. Turn on Developer Mode and run `git config --global core.symlinks true` before you clone, or Git checks the symlink out as a plain text file.

### Troubleshoot

[Troubleshooting](docs/troubleshooting.md) covers the problems people hit most often when installing.

## Best practices

- Keep each pull request to one change.
- Put deployment-specific settings, such as domains and hostnames, in values files, not in the chart's defaults.
- Every subchart has a `values.yaml` with safe defaults, its templates, and a `README.md`.
- Update the documentation in the same pull request as the change it describes.

## Content style guide

Documentation follows the [documentation style guide](docs/style-guide.md) and the templates it lists.

## Contribution workflow

### Fork and clone

1. Fork [OpenGovMail/OpenGovMail](https://github.com/OpenGovMail/OpenGovMail) on GitHub.
2. Clone your fork and add this repository as `upstream`:

   ```bash
   git clone https://github.com/<you>/OpenGovMail.git
   cd OpenGovMail
   git remote add upstream https://github.com/OpenGovMail/OpenGovMail.git
   ```

### Report issues and bugs

Open an issue with the bug template. Include the chart version, what you ran, what you expected and what happened, with logs from `kubectl logs`.

### Issue management

Maintainers label and assign issues. Comment on an issue to ask for it, and link your pull request to it with `Closes #<number>`.

### Branch creation

Create a branch from an up-to-date `main`, named after the issue:

```bash
git switch main
git switch -c 123-short-description
```

### Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/): a type, an optional scope, and a short summary.

```
feat(charts): add Nextcloud Talk with a coturn TURN server
fix(charts): make OpenDKIM sign outgoing mail
docs: document the Helm install
```

Keep the message to a subject line and a line or two of body. Put the details in the pull request.

### Pull requests

1. Open a draft pull request early, so maintainers can follow along and help.
2. Fill in the pull request template, including how you tested the change.
3. Add an entry under `[Unreleased]` in [CHANGELOG.md](CHANGELOG.md) for any change users would notice.
4. Mark it ready for review when CI passes.

A maintainer reviews it, may ask for changes, and merges it into `main`. If a pull request has waited a while without a response, ask in Slack.

### Keeping your fork up to date

Bring your `main` and your branch up to date with `upstream`:

```bash
git fetch upstream
git switch main
git merge --ff-only upstream/main
git push origin main

git switch my-feature
git rebase main
git push --force-with-lease
```

### Running a patched version downstream

If you run OpenGovMail with changes of your own:

- Put settings in your own values file rather than editing the chart. Values files carry over between releases unchanged.
- Keep chart or image changes on a branch in your fork, and rebase it onto each release tag:

  ```bash
  git fetch upstream --tags
  git switch my-patches
  git rebase v<new-version>
  ```

- Build and run patched images as [Run your own build](docs/how-to/run-local-build.md) describes.

### Contributing a patch upstream

A downstream patch that fixes a bug or would help other installs belongs upstream, so you no longer need to carry it. Rebase it onto `upstream/main`, remove anything specific to your install, and open a pull request as above.

### Releases

Releases follow [Semantic Versioning](https://semver.org/). The chart version and the Postfix and OpenDKIM image tags move together. A maintainer prepares each release in a pull request titled `chore(release): prepare v<version>` that:

1. Moves the `[Unreleased]` entries in [CHANGELOG.md](CHANGELOG.md) under the new version, and updates the compare links.
2. Sets `version` and `appVersion` in [charts/opengovmail/Chart.yaml](charts/opengovmail/Chart.yaml).
3. Pins `image.tag` in the Postfix and OpenDKIM subcharts to the new version.

After it merges, the maintainer tags the merge commit `v<version>` and pushes the tag, which builds and publishes the images, then publishes the GitHub release.

A release with breaking changes says so in the changelog, with the steps users need to upgrade.

### Text formats

Documentation is Markdown, and diagrams are [Mermaid](https://mermaid.js.org/) blocks inside it, so GitHub renders them.
