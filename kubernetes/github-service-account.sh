#!/bin/bash
#
# A script for creating the ServiceAccount used by the GitHub Actions CI for
# deploying on AWS
#

main() {
    create_service_account
    create_cluster_role_binding
	create_kubeconfig
}

create_kubeconfig() (
server='https://8EA6D58744BF5F2E3904C549516E9118.gr7.us-east-1.eks.amazonaws.com'
name=github-token-89rp4
ca="$(kubectl get secret/$name -o jsonpath='{.data.ca\.crt}')"
token="$(kubectl get secret/$name -o jsonpath='{.data.token}' | base64 --decode)"
namespace="$(kubectl get secret/$name -o jsonpath='{.data.namespace}' | base64 --decode)"

echo "
apiVersion: v1
kind: Config
clusters:
- name: default-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: default-context
  context:
    cluster: default-cluster
    namespace: default
    user: default-user
current-context: default-context
users:
- name: default-user
  user:
    token: ${token}
" > kubeconfig-github.yml
)

# Create the ServiceAccount
create_service_account() {
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: github
EOF
}

# Create the CluserRoleBinding
create_cluster_role_binding() {
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: github-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: github
    namespace: flapflap
EOF
}

[[ $0 == ${BASH_SOURCE[0]} ]] && main "$@"
