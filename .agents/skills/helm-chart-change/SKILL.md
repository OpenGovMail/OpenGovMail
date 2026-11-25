---
name: helm-chart-change
description: Change the OpenGovMail Helm chart or one of its subcharts (templates, values.yaml, values-example.yaml, Chart.yaml) and check it the way CI does with helm lint and helm template. Use for any edit under charts/opengovmail/, such as adding a value, fixing a template or changing a default.
---

# Change the Helm chart

The chart lives in `charts/opengovmail/`. Local subcharts are in `charts/opengovmail/charts/<name>/`. Nextcloud and ThunderID are remote dependencies, so change their settings through the umbrella `values.yaml`, never by editing their packaged `.tgz` files.

## Steps

1. **Find where the setting belongs.**
   - Shared settings, such as the domain, service names and the TLS Secret, go under `global` in `charts/opengovmail/values.yaml`. Subcharts read them from there.
   - Settings for one service go in that subchart's `values.yaml`, or under the subchart's key in the umbrella `values.yaml`.

2. **Make the change.**
   - Keep defaults generic and safe. Settings that differ per install, such as domains and hostnames, belong in the user's own values file, not in chart defaults.
   - Use `example.com` and `12.34.56.78` in defaults and examples, never a real domain or IP address.
   - Never give a secret a real default. Secrets are set at install time.
   - If a new value is needed for an install to work, add it to `charts/opengovmail/values-example.yaml` too.

3. **Document the value** in the same change:
   - In the subchart's `README.md`, if it has one.
   - In `docs/reference/helm-values.md`, under the matching section, if users set it.

4. **Run the CI checks** from the repository root. The `--set` flags are placeholders for the secrets a first install requires.

   ```bash
   CHART=charts/opengovmail
   # An array, so the flags split the same way in bash and zsh.
   SECRETS=(--set global.ravenIdp.clientSecret=ci --set global.ravenIdp.directAuthSecret=ci --set thunderid.configuration.server.security.directAuthSecret=ci)

   helm repo add nextcloud https://nextcloud.github.io/helm/ 2>/dev/null || true
   helm dependency build "$CHART"
   helm lint "$CHART" -f "$CHART/values-example.yaml" "${SECRETS[@]}"
   helm template ci "$CHART" -f "$CHART/values-example.yaml" "${SECRETS[@]}" > /dev/null
   helm template ci "$CHART" --set global.domain=example.com "${SECRETS[@]}" > /dev/null
   ```

   If CI changes, [.github/workflows/helm-chart.yml](../../../.github/workflows/helm-chart.yml) is the source of truth for these commands.

5. **Check that the change renders as intended.** Compare the rendered chart before and after the change, and check that only the intended parts differ:

   ```bash
   .agents/skills/helm-chart-change/scripts/render-diff.sh upstream/main
   ```

   - The first argument is the base to compare against. Use the branch the change started from.
   - The script renders the base in a temporary git worktree, so it never touches your working tree or stash, and it includes committed, uncommitted and new files.
   - The chart generates passwords and keys on every offline render. The script masks Secret values and `checksum/*` annotations, so they don't show as changes.
   - Extra arguments go to `helm template`. For example, add `--set <name>.enabled=true` to render an optional service.

6. **Add a changelog entry.** If users will notice the change, add it under `[Unreleased]` in `CHANGELOG.md`, in the `Added`, `Changed` or `Fixed` group. Link the pull request.

7. **Say what you couldn't test.** Lint and render don't prove the service works. Running it end to end needs a test server, as described in [docs/how-to/run-local-build.md](../../../docs/how-to/run-local-build.md). If you didn't run it there, say so in the pull request.

## Done when

- [ ] `helm lint` passes with `values-example.yaml`
- [ ] Both `helm template` commands render without errors
- [ ] The rendered diff contains only the intended changes
- [ ] New or changed values are documented in the subchart README and `docs/reference/helm-values.md`
- [ ] `CHANGELOG.md` has an `[Unreleased]` entry if users will notice the change
- [ ] No real domains, IP addresses or secrets were added
