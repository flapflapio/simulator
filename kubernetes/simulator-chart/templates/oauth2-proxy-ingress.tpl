{{- define "ingress-oauth2" -}}
"{{ .Values.name }}-oauth2-ingress"
{{- end }}

{{- if .Values.oauth2.enabled }}
# INGRESS FOR OAUTH2
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ include "ingress-oauth2" }}
  labels:
    name: {{ include "ingress-oauth2" }}
  annotations:
    {{- include "simulator-chart.ingressAnnotations" . | nindent 4 }}
spec:
  {{- if .Values.certManager }}
  tls:
    - hosts:
      - api.flapflap.io
      secretName: api-flapflap-io-cert
  {{- end }}
  rules:
    - host: api.flapflap.io
      http:
        paths:
          - pathType: Prefix
            path: /oauth2
            backend:
              service:
                name: {{ .Values.oauth2.serviceName }}
                port:
                  number: {{ .Values.oauth2.port }}
{{- end }}