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
{{- if $idp.clientSecret -}}
{{- $idp.clientSecret -}}
{{- else -}}
{{- $existing := lookup "v1" "Secret" .Release.Namespace ($idp.secretName | default "") -}}
{{- $data := dict -}}
{{- if $existing }}{{- $data = $existing.data | default dict -}}{{- end -}}
{{- if index $data "clientSecret" -}}
{{- index $data "clientSecret" | b64dec -}}
{{- else -}}
{{- fail "global.ravenIdp.clientSecret is required on a first install: it authenticates raven to ThunderID, and it has to be written into both raven's Secret and the bootstrap data that defines its application, which cannot share a generated value. Later installs reuse the value already in the cluster. Generate one and pass it:\n  --set global.ravenIdp.clientSecret=$(openssl rand -hex 32)" -}}
{{- end -}}
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
