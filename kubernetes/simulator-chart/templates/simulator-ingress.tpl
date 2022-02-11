# INGRESS
#
# Add to your /etc/hosts: '192.168.49.2 api.flapflap.io' replacing the ip
# address by the one associated with the Ingress
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Values.name }}-ingress
  labels:
    name: {{ .Values.name }}-ingress
  annotations:
    {{- include "simulator-chart.ingressAnnotations" . | nindent 4 }}
    {{- if .Values.oauth2.enabled }}
    nginx.ingress.kubernetes.io/auth-url: "https://$host/oauth2/auth"
    nginx.ingress.kubernetes.io/auth-signin: "https://$host/oauth2/start?rd=$escaped_request_uri"
    {{- end}}
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
            path: /
            backend:
              service:
                name: {{ .Values.name }}-service
                port:
                  number: {{ .Values.port }}
