---
name: open-pull-request
description: Open an OpenGovMail pull request the right way, linked to a GitHub issue that tracks the work. Covers finding or opening the issue, checking nobody else is working on it, claiming it, naming the branch, and linking the pull request with Closes or Part of. Use before starting any change that will become a pull request, and when opening one.
---

# Open a pull request

Every pull request in OpenGovMail is linked to an issue. The issue is where people see who is working on what, so two contributors don't build the same thing, and where maintainers agree on the change before code is written. No issue, no pull request.

Run these steps before you start writing code, not just at the end.

## Steps

1. **Find the issue for the work.** Search the open issues for the change:

   ```bash
   gh issue list --repo OpenGovMail/OpenGovMail --state open --search "<keywords>"
   gh issue view <number> --repo OpenGovMail/OpenGovMail --comments
   ```

   The person who asked for the change may already have given you an issue number. If so, read that issue and its comments.

2. **Check that nobody else is working on it.** Look at the assignees and the recent comments.
   - If someone else is assigned, or said in a comment that they're working on it, **stop**. Tell the person you're working for, and suggest they coordinate in the issue or pick a different one.
   - If the issue is closed, or the comments say not to do it, stop and say so.

3. **If there is no issue, propose one.** Draft it with the matching template from `.github/ISSUE_TEMPLATE/`:
   - `bug.yaml` for a defect
   - `feature.yaml` for a new feature or improvement
   - `task.yaml` for other work

   Fill in every required field of the template. Show the draft to the person you're working for, and create the issue only after they agree, because it's public:

   ```bash
   gh issue create --repo OpenGovMail/OpenGovMail --title "<title>" --body-file <file> --label <label>
   ```

   Large changes need agreement in the issue from a maintainer before you write the code, as [CONTRIBUTING.md](../../../CONTRIBUTING.md) asks.

4. **Claim the issue,** so others can see it's taken. Comment on it to say you're working on it. Maintainers assign issues, so ask to be assigned rather than assigning yourself.

5. **Branch from an up-to-date `main`,** named after the issue:

   ```bash
   git fetch upstream
   git switch -c <issue-number>-<short-description> upstream/main
   ```

6. **Make the change** in commits that follow [Conventional Commits](https://www.conventionalcommits.org/), for example `fix(charts): make OpenDKIM sign outgoing mail`.

7. **Open the pull request** against `OpenGovMail/OpenGovMail` `main`, as a draft at first:
   - Fill in every section of [the pull request template](../../../.github/pull_request_template.md), including how you tested the change.
   - Link the issue in the description:
     - `Closes #<number>` when this pull request finishes the issue. GitHub closes the issue when it merges.
     - `Part of #<number>` when the issue needs several pull requests. Use `Closes` in the last one.
   - Keep the pull request to the one change the issue describes. Unrelated fixes get their own issue and pull request.
   - Add an entry under `[Unreleased]` in `CHANGELOG.md` if users will notice the change.

8. **Mark it ready for review** when CI passes.

## Done when

- [ ] An open issue describes the change, and nobody else is working on it
- [ ] The issue has a comment saying you're working on it
- [ ] The branch is named `<issue-number>-<short-description>`
- [ ] The pull request description links the issue with `Closes #<number>` or `Part of #<number>`
- [ ] The pull request contains only the change the issue describes
