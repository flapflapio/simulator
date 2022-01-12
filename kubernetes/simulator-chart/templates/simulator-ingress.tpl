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
spec:
  # ingressClassName: nginx
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
