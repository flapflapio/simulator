
{{/*
Common annotations for all ingresses
*/}}
{{- define "simulator-chart.ingressAnnotations" -}}
kubernetes.io/ingress.class: nginx
{{- if .Values.certManager }}
cert-manager.io/cluster-issuer: {{ .Values.clusterIssuer }}
{{- end }}
ingress.kubernetes.io/force-ssl-redirect: "true"
ingress.kubernetes.io/proxy-body-size: "1024m"
nginx.ingress.kubernetes.io/proxy-body-size: "1024m"
{{- end }}
