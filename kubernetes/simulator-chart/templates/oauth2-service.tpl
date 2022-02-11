{{- if .Values.oauth2.enabled }}
apiVersion: v1
kind: Service
metadata:
  labels:
    app: {{ .Values.oauth2.serviceName }}
  name: {{ .Values.oauth2.serviceName }}
spec:
  ports:
  - name: http
    port: {{ .Values.oauth2.port }}
    protocol: TCP
    targetPort: {{ .Values.oauth2.port }}
  selector:
    app: {{ .Values.oauth2.serviceName }}
{{- end }}