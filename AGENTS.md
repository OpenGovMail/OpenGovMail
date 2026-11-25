# AGENTS.md

Instructions for AI coding agents working in this repository. People should start with [CONTRIBUTING.md](CONTRIBUTING.md), which this file summarizes and defers to.

## Project overview

OpenGovMail is a self-hosted email and collaboration platform, installed with one Helm chart on Kubernetes. The chart brings together Postfix, OpenDKIM, Rspamd, Unbound, Redis, SeaweedFS, Raven (the mail store), Tuwunel (Matrix chat), Nextcloud (calendars, files and Talk) and ThunderID (the identity provider).

Raven lives in its own repository, [OpenGovMail/raven](https://github.com/OpenGovMail/raven). Changes to the mail store belong there, not here.

[Overview of the architecture](docs/concepts/architecture.md) explains how the parts fit together.

## Repository layout

| Path | What it holds |
|---|---|
| `charts/opengovmail/` | The umbrella Helm chart: `Chart.yaml`, `values.yaml`, `values-example.yaml` and templates |
| `charts/opengovmail/charts/<name>/` | Local subcharts (opendkim, postfix, redis, unbound, rspamd, seaweedfs, raven, tuwunel) |
| `services/smtp/`, `services/dkim/`, `services/certbot/` | Dockerfiles for the images CI publishes |
| `infra/` | Cluster setup applied once, such as the cert-manager issuers |
| `docs/` | Documentation, organized by topic type (see [docs/style-guide.md](docs/style-guide.md)) |
| `test/` | Load and security tests |
| `.github/workflows/` | CI: chart lint and render, image builds, Trivy scans |

Nextcloud and ThunderID are pulled in as remote chart dependencies. The `.tgz` files and `Chart.lock` under `charts/opengovmail/` are build output and are not committed.

## Check a chart change

Run the same checks as CI ([.github/workflows/helm-chart.yml](.github/workflows/helm-chart.yml)) before you call a chart change done. The `--set` flags are placeholders for secrets a first install requires.

```bash
CHART=charts/opengovmail
SECRETS="--set global.ravenIdp.clientSecret=ci --set global.ravenIdp.directAuthSecret=ci --set thunderid.configuration.server.security.directAuthSecret=ci"

helm repo add nextcloud https://nextcloud.github.io/helm/
helm dependency build "$CHART"
helm lint "$CHART" -f "$CHART/values-example.yaml" $SECRETS
helm template ci "$CHART" -f "$CHART/values-example.yaml" $SECRETS > /dev/null
helm template ci "$CHART" --set global.domain=example.com $SECRETS > /dev/null
```

All three must pass. To see what a change renders, drop `> /dev/null` and compare the output before and after.

Running a change end to end needs an amd64 test server. [Run your own build](docs/how-to/run-local-build.md) describes it. Don't assume you have access to one; say what you couldn't test.

## Chart conventions

- Put deployment-specific settings, such as domains and hostnames, in values files. Chart defaults stay generic and safe.
- Every subchart has a `values.yaml` with safe defaults, its templates and a `README.md`. When you add or change a value, document it in that README and, for user-facing values, in [docs/reference/helm-values.md](docs/reference/helm-values.md).
- Use `example.com` for domains and `12.34.56.78` for IP addresses in defaults, examples and docs.
- Never commit secrets. Values that hold secrets are set at install time, never given a real default.

## Documentation

- Follow [docs/style-guide.md](docs/style-guide.md). It builds on the Google developer documentation style guide and lists preferred terms (for example, "ThunderID", not "Thunder").
- Each page follows one topic type: concept, how-to, reference or troubleshooting. Add new pages to [docs/README.md](docs/README.md).
- Helm on Kubernetes is the only supported install method. Don't document other ways to install.
- Write in the second person and present tense, use sentence case for headings, and don't use em dashes.
- Diagrams are Mermaid blocks inside Markdown.
- Update the docs in the same pull request as the change they describe.

## Commits and pull requests

- Use [Conventional Commits](https://www.conventionalcommits.org/): `type(scope): summary`, for example `feat(charts): add Nextcloud Talk with a coturn TURN server` or `docs: document the Helm install`. Keep the body to a line or two.
- Name branches after the issue: `<issue-number>-short-description`.
- Keep each pull request to one change, fill in [the pull request template](.github/pull_request_template.md) and link the issue with `Closes #<number>`.
- Add an entry under `[Unreleased]` in [CHANGELOG.md](CHANGELOG.md) for any change users would notice. The changelog follows Keep a Changelog.

## Releases

Releases follow Semantic Versioning and are prepared by a maintainer. Don't start one unless a maintainer asks. The steps are in the "Releases" section of [CONTRIBUTING.md](CONTRIBUTING.md): the chart `version` and `appVersion` and the Postfix and OpenDKIM image tags move together, in a pull request titled `chore(release): prepare v<version>`.

## Security

- Report vulnerabilities privately, as [SECURITY.md](SECURITY.md) describes. Don't put them in issues, pull requests or commit messages.
- Don't include real domains, IP addresses, hostnames, credentials or personal data in code, docs, tests or commit messages.
- Mail delivery depends on TLS, DKIM, SPF and DMARC. Call out any change that touches them in the pull request.
