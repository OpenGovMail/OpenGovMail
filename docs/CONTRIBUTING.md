# Contributing

Contributions are always welcome, no matter how large or small! Before contributing, please read the [code of conduct](https://github.com/OpenGovMail/OpenGovMail/blob/main/docs/CODE_OF_CONDUCT.md).

## Not sure where to start?

- Browse our available Issues on GitHub and select one(s) to work on.
Check out our [Good first issues](https://github.com/OpenGovMail/OpenGovMail/issues?q=is%3Aissue%20state%3Aopen%20label%3A%22good%20first%20issue%22) for a good place to start.

## How to contribute?
- Once you have picked a good issue, if that issue hasn't been assigned to someone else, you're good to go.
- Create a local copy of the OpenGovMail repository by forking the main remote repo and cloning it onto your computer.
- Then, create a separate feature branch in the forked repo. In the feature branch, create your changes.
- Feel free to make a draft Pull Request so maintainers can watch your changes and help you out whenever possible.
- Once you are happy with your changes, make a Pull Request against our main branch.
- Our maintainers will look over your changes, request changes if needed and merge it into the repository. 

## Keeping your fork up to date

Add this repository as `upstream` once:

```bash
git remote add upstream https://github.com/OpenGovMail/OpenGovMail.git
```

Then bring your `main` and your feature branch up to date:

```bash
git fetch upstream
git switch main
git merge --ff-only upstream/main
git push origin main

git switch my-feature
git rebase main
git push --force-with-lease
```

To run your changes on a server, see [Running your own build](../dev/README.md#running-your-own-build-on-kubernetes).
