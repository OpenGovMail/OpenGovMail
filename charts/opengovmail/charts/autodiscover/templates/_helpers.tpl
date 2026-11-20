{{- define "autodiscover.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "autodiscover.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{- define "autodiscover.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "autodiscover.labels" -}}
helm.sh/chart: {{ include "autodiscover.chart" . }}
app.kubernetes.io/name: {{ include "autodiscover.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "autodiscover.selectorLabels" -}}
app.kubernetes.io/name: {{ include "autodiscover.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "autodiscover.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "autodiscover.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
The domain list as a JSON array. Empty values.domains derives global.domain, so
the umbrella's one-command install works without repeating the domain here.
*/}}
{{- define "autodiscover.domains" -}}
{{- if .Values.domains -}}
{{- toJson .Values.domains -}}
{{- else -}}
{{- toJson (list (required "autodiscover: set global.domain, or list autodiscover.domains explicitly" ((.Values.global).domain))) -}}
{{- end -}}
{{- end -}}

{{- define "autodiscover.domain" -}}
{{- first (include "autodiscover.domains" . | fromJsonArray) -}}
{{- end -}}

{{- define "autodiscover.mailHost" -}}
{{- default (printf "mail.%s" (include "autodiscover.domain" .)) (((.Values.global).mailHostname)) -}}
{{- end -}}

{{- define "autodiscover.ingressHosts" -}}
{{- if .Values.ingress.hosts -}}
{{- toJson .Values.ingress.hosts -}}
{{- else -}}
{{- toJson (list (include "autodiscover.domain" .) (include "autodiscover.mailHost" .)) -}}
{{- end -}}
{{- end -}}

{{/* Takes a dict with a "host" key. Matches certificate.yaml's naming. */}}
{{- define "autodiscover.tlsSecretName" -}}
{{- printf "%s-tls" (.host | replace "." "-") -}}
{{- end -}}
