
{{/*
Common annotations for all ingresses
*/}}
{{- define "simulator-chart.ingressAnnotations" -}}
kubernetes.io/ingress.class: nginx
{{- if .Values.certManager }}
cert-manager.io/cluster-issuer: {{ .Values.clusterIssuer }}
{{- end }}
{{- end }}
