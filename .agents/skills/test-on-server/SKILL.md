---
name: test-on-server
description: Prepare the commands to try an OpenGovMail chart or image change end to end on a contributor's own test server, following docs/how-to/run-local-build.md. Use when a change has passed helm lint and helm template and needs testing on a real cluster, or when someone asks how to run their branch or a changed Postfix, OpenDKIM or Raven image.
---

# Test a change on a server

Lint and render checks don't prove a change works. Running it on a test server does. [Run your own build](../../../docs/how-to/run-local-build.md) is the source of truth for these steps.

The test server belongs to the contributor, and other people may use it. **Don't log in to it, and don't run commands on it yourself.** Prepare the commands, explain what each one does, and let the contributor run them and share the results.

## Before you start

Check with the contributor that the test server meets these needs. If it doesn't, say so and stop.

- It is a real server on the internet, so Let's Encrypt can issue trusted certificates. Sign-in doesn't work without them.
- It is amd64 (x86_64), because the published images are amd64 only.
- OpenGovMail is already installed on it, as in the [installation guide](../../../docs/installation-guide.md), with the `le-prod` issuer and a spare domain or subdomain.

## Steps

1. **Run the chart checks first,** as in the `helm-chart-change` skill. Don't send a change to a server if it doesn't render.

2. **For a chart change,** the contributor pushes the branch to their fork, then runs on the server:

   ```bash
   git fetch origin
   git checkout <branch>
   helm dependency update ./charts/opengovmail
   ```

   Then they run the `helm upgrade` command from [Upgrade OpenGovMail](../../../docs/how-to/upgrade.md) with their own values file.

3. **For a changed image,** follow the "Use your own image" section of the guide:
   1. Build it on an amd64 machine with a **new tag for every build**, such as the commit hash. If you reuse a tag like `dev`, the pod spec doesn't change, so Kubernetes keeps running the old image. With `IfNotPresent`, a restart can also reuse the cached copy.

      ```bash
      TAG=dev-$(git rev-parse --short HEAD)
      docker build -t ghcr.io/<you>/opengovmail-smtp:$TAG -f services/smtp/Dockerfile services
      ```

      Use `services/dkim/Dockerfile` for OpenDKIM. Raven is built from the [raven repository](https://github.com/OpenGovMail/raven).
   2. Make the image available to the cluster: either push it to a registry, or on a single-node k3s server import it with `docker save <image> | sudo k3s ctr images import -`.
   3. Point the release at it with `helm upgrade ... --reuse-values --set postfix.image.repository=ghcr.io/<you>/opengovmail-smtp --set postfix.image.tag=$TAG`. Add `--set postfix.image.pullPolicy=Never` for an imported image. Use `opendkim.image.*` or `raven.image.*` for those images.

4. **Say what to check,** based on what the change touches. For example:
   - `kubectl get pods -n opengovmail`: all pods are running and none are restarting.
   - `kubectl logs -n opengovmail <pod>`: no new errors in the changed service.
   - For mail changes, send and receive a message from a mail app, and check the DKIM, SPF and DMARC results in the received message's headers.

5. **Say how to go back.** To return to the published images, run the `helm upgrade` command from [Upgrade OpenGovMail](../../../docs/how-to/upgrade.md) as it is: with the contributor's own values file, and without `--reuse-values` or any of the `--set ...image...` overrides from step 3. Helm then takes the image settings from the chart's defaults again.

6. **Report the results in the pull request,** under "Testing Instructions", with log excerpts if they help. Leave out the server's domain, IP address, hostnames and any secrets. Use `example.com` and `12.34.56.78` in their place.

## Never

- Paste secrets, passwords, kubeconfig contents or values files into chat, commits or pull requests.
- Use a production server for testing.
- Commit the contributor's values file or anything specific to their server.
