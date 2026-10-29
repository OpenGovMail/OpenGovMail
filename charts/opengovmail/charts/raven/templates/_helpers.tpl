{{- define "raven.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "raven.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name (include "raven.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/* Stable Service name so Postfix / delivery.yaml resolve raven's socketmap. */}}
{{- define "raven.serviceName" -}}
{{- (((.Values.global).serviceNames).raven) | default "raven" }}
{{- end }}

{{- define "raven.selectorLabels" -}}
app.kubernetes.io/name: {{ include "raven.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "raven.labels" -}}
{{ include "raven.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/* Primary domain + derived mail hostname (single source: global.domain). */}}
{{- define "raven.domain" -}}
{{- required "global.domain is required (e.g. --set global.domain=example.com)" ((.Values.global).domain) }}
{{- end }}

{{- define "raven.mailHostname" -}}
{{- (.Values.global).mailHostname | default (printf "mail.%s" (include "raven.domain" .)) }}
{{- end }}

{{/* S3 endpoint (defaults to the in-cluster seaweedfs Service). */}}
{{- define "raven.s3Endpoint" -}}
{{- if .Values.blobStorage.endpoint }}
{{- .Values.blobStorage.endpoint }}
{{- else }}
{{- printf "http://%s:8333" ((((.Values.global).serviceNames).s3) | default "seaweedfs-s3") }}
{{- end }}
{{- end }}

{{/*
S3 accessKey / secretKey. Priority: global.s3.* > local blobStorage.* > derived.
The secretKey derivation MUST match the seaweedfs subchart's
(`printf "%s-opengovmail-seaweedfs-s3" .Release.Name | sha256sum | trunc 40`) so both
agree on a fresh one-command install without any cross-lookup. Override
global.s3.secretKey in prod.
*/}}
{{- define "raven.s3AccessKey" -}}
{{- (((.Values.global).s3).accessKey) | default .Values.blobStorage.accessKey }}
{{- end }}

{{- define "raven.s3SecretKey" -}}
{{- (((.Values.global).s3).secretKey) | default .Values.blobStorage.secretKey | default (printf "%s-opengovmail-seaweedfs-s3" .Release.Name | sha256sum | trunc 40) }}
{{- end }}

{{/*
Thunder hosts. Both default to the BARE domain, not mail.<domain>: Thunder serves
the <domain> certificate on :8090 and publishes its OAuth issuer as
https://<domain>:8090, so mail.<domain> is not a SAN and TLS verification fails.

Both are the public name rather than the in-cluster Service name because raven
speaks HTTPS with no skip-verify or custom-CA option, and opengovmail-thunder-service
is not on the cert either. In-cluster calls therefore hairpin via the node's
external IP (see finding #7).

global.thunder.url (an external ThunderID) supplies host and port when the
matching thunder.host / thunder.publicHost is empty.
*/}}
{{- define "raven.externalThunder" -}}
{{- with ((.Values.global).thunder).url -}}
{{- $u := urlParse . -}}
{{- $port := trimPrefix (printf "%s:" $u.hostname) $u.host -}}
{{- dict "host" $u.hostname "port" (ternary "443" $port (eq $port $u.host)) | toJson -}}
{{- end -}}
{{- end }}

{{- define "raven.thunderHost" -}}
{{- $ext := include "raven.externalThunder" . | default "{}" | fromJson -}}
{{- .Values.thunder.host | default $ext.host | default (include "raven.domain" .) }}
{{- end }}

{{- define "raven.thunderPort" -}}
{{- $ext := include "raven.externalThunder" . | default "{}" | fromJson -}}
{{- if and $ext.port (not .Values.thunder.host) }}{{ $ext.port }}{{ else }}{{ .Values.thunder.port }}{{ end }}
{{- end }}

{{- define "raven.thunderPublicHost" -}}
{{- $ext := include "raven.externalThunder" . | default "{}" | fromJson -}}
{{- .Values.thunder.publicHost | default $ext.host | default (include "raven.domain" .) }}
{{- end }}

{{- define "raven.thunderPublicPort" -}}
{{- $ext := include "raven.externalThunder" . | default "{}" | fromJson -}}
{{- if and $ext.port (not .Values.thunder.publicHost) }}{{ $ext.port }}{{ else }}{{ .Values.thunder.publicPort }}{{ end }}
{{- end }}

{{/*
cert-manager TLS Secret for IMAPS. Mirrors the umbrella's opengovmail.tlsSecretName
(charts/opengovmail/templates/_helpers.tpl) rather than calling it, so this subchart
still renders standalone — the same reason raven.domain/raven.mailHostname are
duplicated above.

The dots->dashes naming rule exists in three places and they MUST agree:
  charts/opengovmail/templates/_helpers.tpl   (opengovmail.tlsSecretName)
  charts/opengovmail/templates/certificate.yaml (the Certificate's secretName)
  here
*/}}
{{- define "raven.tlsSecretName" -}}
{{- ((.Values.global).tlsSecretName) | default (printf "%s-tls" (include "raven.mailHostname" . | replace "." "-")) }}
{{- end }}

{{/* Secret holding the rendered raven.yaml + delivery.yaml (contain S3 key). */}}
{{- define "raven.configSecretName" -}}
{{- printf "%s-config" (include "raven.fullname" .) }}
{{- end }}

{{/*
Secret mounted at /certs. Precedence:
  tls.existingSecret -> global.tlsSecretName -> mail-<domain-dashed>-tls -> self-signed
The cert-manager name is preferred over the generated self-signed Secret so a
one-command install with global.domain set picks up the real cert automatically.
Falls back to self-signed only when there is no global.domain at all (standalone).
*/}}
{{- define "raven.certSecretName" -}}
{{- if .Values.tls.existingSecret }}
{{- .Values.tls.existingSecret }}
{{- else if and (((.Values.global).tls).enabled) ((.Values.global).domain) }}
{{- include "raven.tlsSecretName" . }}
{{- else }}
{{- printf "%s-certs" (include "raven.fullname" .) }}
{{- end }}
{{- end }}
