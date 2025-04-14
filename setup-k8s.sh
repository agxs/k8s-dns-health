kubectl create serviceaccount dns-health-sa --namespace kube-system

cat <<EOF | kubectl apply -f -
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: dns-health-role
  namespace: kube-system
rules:
  - apiGroups: [""]
    resources: ["pods"]
    verbs: ["get", "list", "delete"]
EOF

kubectl create rolebinding dns-health-rolebinding \
  --namespace kube-system \
  --role=dns-health-role \
  --serviceaccount=kube-system:dns-health-sa


