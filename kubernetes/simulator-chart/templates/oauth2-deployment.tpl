{{- if .Values.oauth2.enabled }}
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: oauth2-proxy
  name: oauth2-proxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: oauth2-proxy
  template:
    metadata:
      labels:
        app: oauth2-proxy
    spec:
      containers:
        - name: oauth2-proxy
          image: quay.io/oauth2-proxy/oauth2-proxy:v7.2.1-amd64
          imagePullPolicy: Always
          ports:
            - containerPort: 4180
              protocol: TCP
          args:
          - --provider=oidc
          - --email-domain=*
          - --provider-display-name=aws_cognito
          - --upstream=file:///dev/null
          - --http-address=0.0.0.0:4180

          # We are setting these in the env vars below
          # - --client-id=MY_USER_POOL_DEVICE_ID
          # - --client-secret=MYUSER_POOL_DEVICE_SECRET

          # - --redirect-url=https://<my-cloud-application-url>/oauth2/callback
          # - --login-url=https://<my-gitlab-url>/oauth/authorize
          # - --redeem-url=https://<my-gitlab-url>/oauth/token
          - --oidc-jwks-url=https://cognito-idp.eu-east-1.amazonaws.com/us-east-1_Nf2xbj7Lq/.well-known/jwks.json
          - --oidc-issuer-url=https://cognito-idp.us-east-1.amazonaws.com/us-east-1_Nf2xbj7Lq

          # - --skip-oidc-discovery
          - --cookie-secure=false
          - --scope=openid
          # - --skip-provider-button=false
          # - --set-xauthrequest=true
          # - --skip-auth-preflight=false
          env:
          - name: OAUTH2_PROXY_CLIENT_ID
            value: {{ .Values.oauth2.clientId }}
          # - name: OAUTH2_PROXY_CLIENT_SECRET
          #   value: <paste-the-secret>
          - name: OAUTH2_PROXY_COOKIE_SECRET
            value: {{ .Values.oauth2.cookieSecret }}
{{- end}}