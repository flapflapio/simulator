# SERVICE
---
apiVersion: v1
kind: Service
metadata:
  name: "{{ .Values.name }}-service"
spec:
  selector:
    app: {{ .Values.name }}
  ports:
    - port: {{ .Values.port }}
      targetPort: {{ .Values.port }}
