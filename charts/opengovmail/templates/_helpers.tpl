{{/*
Expand the name of the chart.
*/}}
{{- define "opengovmail.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "opengovmail.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "opengovmail.labels" -}}
helm.sh/chart: {{ include "opengovmail.name" . }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "opengovmail.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "opengovmail.selectorLabels" -}}
app.kubernetes.io/name: {{ include "opengovmail.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Mail domain (single source of truth).
*/}}
{{- define "opengovmail.domain" -}}
{{- required "global.domain is required (e.g. --set global.domain=example.com)" .Values.global.domain -}}
{{- end }}

{{/*
Mail hostname, defaults to mail.<domain>.
*/}}
{{- define "opengovmail.mailHostname" -}}
{{- default (printf "mail.%s" (include "opengovmail.domain" .)) .Values.global.mailHostname -}}
{{- end }}

{{/*
cert-manager TLS Secret name, defaults to the mail.<domain> cert minted by certificate.yaml
(<domain-dashed>-tls).
*/}}
{{- define "opengovmail.tlsSecretName" -}}
{{- default (printf "%s-tls" (include "opengovmail.mailHostname" . | replace "." "-")) .Values.global.tlsSecretName -}}
{{- end }}

{{/*
Issuer reference helper
*/}}
{{- define "opengovmail.issuerRef" -}}
name: {{ .Values.global.tls.issuer }}
kind: ClusterIssuer
{{- end }}
{{/*
Raven's ThunderID credentials, resolved identically wherever they are needed.

Both are needed in two templates at once — the Secret raven reads, and either the
bootstrap data that defines raven's application or the subchart's own config — and
separate templates cannot share a generated value. So each resolves the same way:
an explicit value wins, otherwise whatever is already in the cluster is reused,
otherwise the render fails with instructions rather than inventing a value that
only one of the two consumers would see.
*/}}
{{- define "opengovmail.ravenClientSecret" -}}
{{- $idp := .Values.global.ravenIdp | default dict -}}
{{- $existing := lookup "v1" "Secret" .Release.Namespace ($idp.secretName | default "") -}}
{{- $data := dict -}}
{{- if $existing }}{{- $data = $existing.data | default dict -}}{{- end -}}
{{- $current := "" -}}
{{- if index $data "clientSecret" }}{{- $current = index $data "clientSecret" | b64dec -}}{{- end -}}
{{- if $idp.clientSecret -}}
{{- /* ThunderID's setup job only runs on install, so a changed value would reach
     raven's Secret but never the registered application, and raven's token
     requests would start failing. Refuse it instead. */ -}}
{{- if and $current (ne $current $idp.clientSecret) -}}
{{- fail (printf "global.ravenIdp.clientSecret differs from the value already in Secret %s. ThunderID's setup job only runs on install, so the new value would never reach raven's registered application and raven could no longer authenticate. Leave global.ravenIdp.clientSecret empty to keep the current value. To rotate it, change the secret on the Raven System application in ThunderID first, then delete Secret %s and upgrade with the new value." $idp.secretName $idp.secretName) -}}
{{- end -}}
{{- $idp.clientSecret -}}
{{- else if $current -}}
{{- $current -}}
{{- else -}}
{{- fail "global.ravenIdp.clientSecret is required on a first install: it authenticates raven to ThunderID, and it has to be written into both raven's Secret and the bootstrap data that defines its application, which cannot share a generated value. Later installs reuse the value already in the cluster. Generate one and pass it:\n  --set global.ravenIdp.clientSecret=$(openssl rand -hex 32)" -}}
{{- end -}}
{{- end }}

{{- define "opengovmail.ravenDirectAuthSecret" -}}
{{- $idp := .Values.global.ravenIdp | default dict -}}
{{- if $idp.directAuthSecret -}}
{{- $idp.directAuthSecret -}}
{{- else -}}
{{- $existing := lookup "v1" "Secret" .Release.Namespace ($idp.secretName | default "") -}}
{{- $data := dict -}}
{{- if $existing }}{{- $data = $existing.data | default dict -}}{{- end -}}
{{- if index $data "directAuthSecret" -}}
{{- index $data "directAuthSecret" | b64dec -}}
{{- else -}}
{{- fail "global.ravenIdp.directAuthSecret is required on a first install: it gates ThunderID's password-check endpoint, and without it every mail login fails with 401. It must be passed twice, because Helm cannot template the vendored subchart's values:\n  --set global.ravenIdp.directAuthSecret=$SECRET \\\n  --set thunderid.configuration.server.security.directAuthSecret=$SECRET" -}}
{{- end -}}
{{- end -}}
{{- end }}

{{/*
A credential generated once and kept: explicit value, else the one already in the
Secret, else random. Every consumer reads the Secret, so one render is enough.
A changed explicit value fails, because the consumer only took the old one on first
start. lookup sees nothing under helm template, so set the value explicitly there.
Takes: ctx, secret, key, value, flag (the values path), owner (who holds the old one).
*/}}
{{- define "opengovmail.keptSecret" -}}
{{- $existing := lookup "v1" "Secret" .ctx.Release.Namespace .secret -}}
{{- $current := "" -}}
{{- if and $existing (index ($existing.data | default dict) .key) }}{{- $current = index $existing.data .key | b64dec -}}{{- end -}}
{{- if and .value $current (ne .value $current) -}}
{{- fail (printf "%s differs from the value already in Secret %s. %s set the old value on first start and will not pick up a new one. Leave %s empty to keep the current value. To rotate it, change it in %s first, then delete Secret %s and upgrade with the new value." .flag .secret .owner .flag .owner .secret) -}}
{{- end -}}
{{- .value | default $current | default (randAlphaNum 40) -}}
{{- end }}

{{- define "opengovmail.nextcloudDbPassword" -}}
{{- include "opengovmail.keptSecret" (dict "ctx" . "secret" .Values.nextcloudDeps.postgres.secretName "key" "password" "value" (.Values.nextcloudDeps.postgres.password | default "") "flag" "nextcloudDeps.postgres.password" "owner" "Postgres") -}}
{{- end }}

{{- define "opengovmail.nextcloudClientSecret" -}}
{{- include "opengovmail.keptSecret" (dict "ctx" . "secret" .Values.nextcloudDeps.idpSecretName "key" "clientSecret" "value" ((((.Values.global).nextcloudIdp).clientSecret) | default "") "flag" "global.nextcloudIdp.clientSecret" "owner" "Nextcloud's user_oidc provider") -}}
{{- end }}

{{/*
Hostname Nextcloud is served on.
*/}}
{{- define "opengovmail.nextcloudHost" -}}
{{- required "nextcloud.nextcloud.host is required when nextcloud.enabled" ((.Values.nextcloud).nextcloud).host -}}
{{- end }}
