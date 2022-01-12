# DEPLOYMENT
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.name }}-deployment
spec:
  replicas: {{ .Values.replicas }}
  selector:
    matchLabels:
      app: {{ .Values.name }}
  template:
    metadata:
      labels:
        app: {{ .Values.name }}
    spec:
      containers:
        - name: {{ .Values.name }}
          imagePullPolicy: {{ .Values.imagePullPolicy }}
          image: "{{ .Values.imageRepo }}:{{ .Values.imageTag }}"
          resources:
            limits:
              memory: "128Mi"
              cpu: "500m"
            requests:
              memory: "0"
              cpu: "0"
          ports:
            - containerPort: {{ .Values.port }}
