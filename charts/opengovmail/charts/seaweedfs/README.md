# seaweedfs subchart

Single-node [SeaweedFS](https://github.com/seaweedfs/seaweedfs) providing an
S3-compatible API for Raven attachment (blob) storage. Internal only — never
exposed publicly.

## What it creates

- **StatefulSet** running combined `weed server -s3` (master + volume + filer +
  S3 gateway in one process), with a PVC for `/data` via `volumeClaimTemplates`.
- **Service** named `seaweedfs-s3` (stable, from `global.serviceNames.s3`) —
  port `8333` (S3) and `9333` (master, for `weed shell`). Raven reads/writes
  attachments at `http://seaweedfs-s3:8333`.
- **Secret** `<release>-seaweedfs-s3` holding `accessKey`, `secretKey`, and the
  rendered `s3-config.json`. Raven consumes the same Secret so credentials can't
  drift.
- **bucket-init Job** (post-install/upgrade hook) creating the
  `email-attachments` bucket. Idempotent; SeaweedFS also auto-creates on first
  PUT since the identity has `Admin`.

## Credentials

`s3.accessKey` is **pinned** (default `raven`) — it binds to `secretKey` inside
`s3-config.json`, so rotating one without the other breaks existing clients.
`s3.secretKey: ""` auto-generates a key once and preserves it across upgrades
(a `lookup` guard reads the existing Secret; `helm.sh/resource-policy: keep`
protects it from deletion). Set `s3.secretKey` explicitly to pin your own.

## Deployment mode — combined vs 4-role split

Default is **combined** mode (one pod). The retired Docker Compose setup used
the **4-role split** (separate master/volume/filer/s3 containers) — that split
is the *proven* config; combined mode is the lazy single-pod path and has known
startup-ordering quirks (filer/master race).

**Fallback trigger — switch to the 4-role split if:**
- the pod crash-loops on startup, or
- S3 / bucket operations fail against combined mode (the bucket-init Job or
  `helm test` keeps failing to reach the cluster).

The split is not implemented yet (`mode` values other than `combined` `fail`
the render). To add it, run the four roles as separate containers in one pod
(or separate StatefulSets) sharing the data PVC:

```text
weed master -ip.bind=0.0.0.0 -port=9333
weed volume -mserver=<master>:9333 -ip.bind=0.0.0.0 -port=8080 -max=10
weed filer -master=<master>:9333 -ip.bind=0.0.0.0
weed s3 -filer=<filer>:8888 -ip.bind=0.0.0.0 -port=8333 -config=/etc/seaweedfs/s3.json
```

## Storage class

`persistence.storageClass: ""` uses the cluster default. Clusters **without** a
default StorageClass (OpenShift, some k3d setups) must set it, or the PVC stays
`Pending`.
