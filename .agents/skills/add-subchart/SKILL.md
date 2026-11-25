---
name: add-subchart
description: Add a new service to OpenGovMail as a local Helm subchart under charts/opengovmail/charts/, wire it into the umbrella Chart.yaml and values.yaml, and document it. Use when adding a new component or service to the chart, not for changing an existing one (use helm-chart-change for that).
---

# Add a subchart

Each OpenGovMail service is a local subchart in `charts/opengovmail/charts/<name>/`, listed as a dependency of the umbrella chart and switched on with `<name>.enabled`. Discuss a new service in an issue before you add it, as [CONTRIBUTING.md](../../../CONTRIBUTING.md) asks for large changes.

## Steps

1. **Copy the closest existing subchart** as a starting point, rather than running `helm create`. `redis` and `unbound` are small examples with a StatefulSet, a Service and a ServiceAccount:

   ```bash
   cp -R charts/opengovmail/charts/unbound charts/opengovmail/charts/<name>
   ```

   Then, in every file of the copy, rename the chart and the template helpers (for example `unbound.fullname` becomes `<name>.fullname`), and remove what the new service doesn't need.

2. **Write `Chart.yaml`.** Use `apiVersion: v2`, `name: <name>`, `type: application`, `version: 0.1.0`, an `appVersion` that matches the image version, and `maintainers: [{ name: OpenGovMail Team }]`.

3. **Write `values.yaml` with safe defaults.**
   - Pin `image.tag` to a specific version, never `latest`.
   - Read shared settings, such as the domain and the TLS Secret, from `.Values.global`. Don't add your own copies of them.
   - Never give a secret a real default. Use `example.com` and `12.34.56.78` in examples.
   - Set resource requests, and run as a non-root user where the image allows it.

4. **If other services must reach it by name,** add the name under `global.serviceNames` in `charts/opengovmail/values.yaml`, and make the subchart create a Service with exactly that name. Postfix and other vendored configs depend on these stable names.

5. **Wire it into the umbrella chart.**
   - In `charts/opengovmail/Chart.yaml`, add the dependency next to the other local ones:

     ```yaml
       - name: <name>
         version: 0.1.0
         repository: file://charts/<name>
         condition: <name>.enabled
     ```

   - In `charts/opengovmail/values.yaml`, add a section:

     ```yaml
     <name>:
       enabled: true
     ```

   - If an install must set something for the service to work, add it to `charts/opengovmail/values-example.yaml`.

6. **Rebuild the dependencies** and run the checks in the `helm-chart-change` skill. `Chart.lock` and the `.tgz` packages are build output and are not committed.

   ```bash
   helm dependency update charts/opengovmail
   ```

   Then render with `<name>.enabled` set to `true` and to `false`, and check that disabling it removes all its resources.

7. **Document it:**
   - `charts/opengovmail/charts/<name>/README.md`: what the service does and a table of its values (key, type, default, description). See `charts/opengovmail/charts/redis/README.md` for the format.
   - `docs/reference/helm-values.md`: the values users set, under the matching section.
   - `docs/concepts/architecture.md`: where the service fits.
   - `docs/reference/hostnames-and-ports.md`: any hostname or port it exposes.
   - `CHANGELOG.md`: an `Added` entry under `[Unreleased]`.

   Use the `write-docs` skill for the pages under `docs/`.

## Done when

- [ ] The subchart has `Chart.yaml`, `values.yaml`, templates and a `README.md`
- [ ] The umbrella `Chart.yaml` lists it with `condition: <name>.enabled`, and `values.yaml` has its section
- [ ] The image tag is pinned and no secret has a real default
- [ ] The chart checks pass, and the chart renders with the service both enabled and disabled
- [ ] The docs and `CHANGELOG.md` are updated
