---
name: prepare-release
description: Prepare an OpenGovMail release pull request (chore(release) prepare vX.Y.Z) by moving CHANGELOG.md entries under the new version, bumping the chart version and appVersion, and pinning the Postfix and OpenDKIM image tags. Use only when a maintainer explicitly asks to prepare a release; never start a release on your own.
---

# Prepare a release

Releases follow [Semantic Versioning](https://semver.org/). The chart version, the chart `appVersion` and the Postfix and OpenDKIM image tags always move together. The full process is in the "Releases" section of [CONTRIBUTING.md](../../../CONTRIBUTING.md). If it disagrees with this skill, CONTRIBUTING.md wins.

## Steps

1. **Confirm the version with the maintainer.** Read `[Unreleased]` in `CHANGELOG.md` and propose a version:
   - **Patch** (`0.6.1` to `0.6.2`): only fixes.
   - **Minor** (`0.6.1` to `0.7.0`): new features.
   - **Major**: breaking changes, once the project reaches 1.0.

   If `[Unreleased]` is empty, stop and say there is nothing to release.

2. **Branch from an up-to-date `main`:**

   ```bash
   git fetch upstream
   git switch -c release/v<version> upstream/main
   ```

3. **Bump the versions** with the bundled script, run from the repository root:

   ```bash
   .agents/skills/prepare-release/scripts/bump-version.sh <version>
   ```

   It sets `version` and `appVersion` in `charts/opengovmail/Chart.yaml`, and `image.tag` in the `postfix` and `opendkim` subchart `values.yaml` files. Check the result with `git diff`.

4. **Update `CHANGELOG.md`:**
   1. Under `## [Unreleased]`, add a `## [<version>] - <YYYY-MM-DD>` heading with today's date, so the existing entries move under it and `[Unreleased]` is left empty.
   2. Under the new version's `### Changed`, add: ``The postfix (`opengovmail-smtp`) and opendkim (`opengovmail-dkim`) images are pinned to `<version>`.`` Create the `### Changed` group if there isn't one.
   3. If the release has breaking changes, add an `### Upgrading` section with the steps users need to take.
   4. Update the compare links at the bottom of the file:

      ```markdown
      [unreleased]: https://github.com/OpenGovMail/OpenGovMail/compare/v<version>...HEAD
      [<version>]: https://github.com/OpenGovMail/OpenGovMail/compare/v<previous>...v<version>
      ```

5. **Check the chart** with the commands in the `helm-chart-change` skill, or in the "Check a chart change" section of `AGENTS.md`.

6. **Commit and open the pull request:**

   ```bash
   git commit -am "chore(release): prepare v<version>"
   ```

   Title the pull request `chore(release): prepare v<version>`.

7. **Stop there.** After the pull request merges, a maintainer tags the merge commit `v<version>`, pushes the tag (which builds and publishes the images) and publishes the GitHub release. Don't create or push tags yourself.

## Done when

- [ ] `Chart.yaml` `version` and `appVersion` and both image tags show the new version
- [ ] `CHANGELOG.md` has the new version with today's date, an empty `[Unreleased]` and updated compare links
- [ ] The chart checks pass
- [ ] The pull request is titled `chore(release): prepare v<version>` and no tag was pushed
