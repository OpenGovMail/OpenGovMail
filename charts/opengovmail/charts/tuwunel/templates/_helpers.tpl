{{- define "tuwunel.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "tuwunel.fullname" -}}
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

{{- define "tuwunel.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "tuwunel.labels" -}}
helm.sh/chart: {{ include "tuwunel.chart" . }}
{{ include "tuwunel.selectorLabels" . }}
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "tuwunel.selectorLabels" -}}
app.kubernetes.io/name: {{ include "tuwunel.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Matrix server name. Defaults to the mail domain so Matrix IDs match addresses.
*/}}
{{- define "tuwunel.serverName" -}}
{{- .Values.serverName | default (((.Values.global).domain)) | required "tuwunel.serverName or global.domain is required" -}}
{{- end }}

{{- define "tuwunel.host" -}}
{{- .Values.host | default (printf "matrix.%s" (include "tuwunel.serverName" .)) -}}
{{- end }}

{{/*
ThunderID issuer. Mirrors opengovmail.thunderUrl, which a subchart cannot call.
*/}}
{{- define "tuwunel.issuerUrl" -}}
{{- .Values.oidc.issuerUrl | default (((.Values.global).thunder).url) | default (printf "https://%s:8090" (include "tuwunel.serverName" .)) | trimSuffix "/" -}}
{{- end }}

{{- define "tuwunel.tlsSecretName" -}}
{{- .Values.ingress.tlsSecretName | default (printf "%s-tls" (include "tuwunel.host" . | replace "." "-")) -}}
{{- end }}

{{- define "tuwunel.wellKnownTlsSecretName" -}}
{{- .Values.wellKnown.tlsSecretName | default (printf "%s-tls" (include "tuwunel.serverName" . | replace "." "-")) -}}
{{- end }}

{{- define "tuwunel.ingressAnnotations" -}}
{{- if and (((.Values.global).tls).enabled) (((.Values.global).tls).issuer) }}
cert-manager.io/cluster-issuer: {{ .Values.global.tls.issuer }}
{{- end }}
{{- with .Values.ingress.annotations }}
{{ toYaml . }}
{{- end }}
{{- end }}
