```bash
kubectl create namespace observability || true

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

tee values.yaml > /dev/null <<EOF
grafana:
  adminPassword: "supersecret"
  ingress:
    enabled: true
    ingressClassName: traefik
    hosts:
      - grafana.k8s.orb.local
    paths:
      - /
prometheus:
  ingress:
    enabled: false
prometheusOperator:
  enabled: true
EOF

helm install obs prometheus-community/kube-prometheus-stack \
  -n observability \
  -f values.yaml
```
