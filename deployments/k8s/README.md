# 🚀 **Quick Deploy**

> **Note:** Run all commands from the project root directory (`FitByte/`)

## Prerequisites

- Ensure you have `kubectl` installed and configured to interact with your Kubernetes cluster.
- Ensure `helm` is installed and configured if using Helm charts.

## Install Traefik

```bash
TRAEFIK_DASHBOARD_INGRESS_DOMAIN=dashboard.k8s.orb.local

helm repo add traefik https://traefik.github.io/charts
helm repo update

tee deployments/k8s/traefik/values.yaml > /dev/null <<EOF
ports:
  metrics:
    expose:
      default: true

metrics:
  addInternals: true
  prometheus:
    addEntryPointsLabels: true
    addRoutersLabels: true
    addServicesLabels: true
    buckets: "0.1,0.3,1.2,5.0,10.0,30.0,60.0"
    service:
      enabled: true
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9100"
        prometheus.io/path: "/metrics"
    serviceMonitor:
      enabled: true
      namespace: observability
      additionalLabels:
        release: obs # Ganti sesuai release name prometheus stack Anda
      interval: 15s
      scrapeTimeout: 10s
      honorLabels: true

ingressRoute:
  dashboard:
    enabled: true
    matchRule: "Host(\`$TRAEFIK_DASHBOARD_INGRESS_DOMAIN\`)"
    entryPoints:
      - web

global:
  checkNewVersion: false
  sendAnonymousUsage: false
EOF

helm upgrade --install traefik traefik/traefik \
  --version 37.1.0 \
  -f deployments/k8s/traefik/values.yaml \
  -n traefik-system \
  --create-namespace
```

## Install PostgreSQL

```bash
PG_POSTGRES_PASSWORD=supersecret
PG_PASSWORD=supersecret
PG_REPLICATION_PASSWORD=supersecret

# Create Kubernetes Namespace
kubectl create namespace newton-db

# Install Kubernetes Secret
kubectl create secret generic pg-fitbyte-secret \
  --from-literal=postgres-password=$PG_POSTGRES_PASSWORD \
  --from-literal=password=$PG_PASSWORD \
  --from-literal=replication-password=$PG_REPLICATION_PASSWORD \
  -n newton-db

# Install PostgreSQL with custom values
helm upgrade --install newton-pg oci://registry-1.docker.io/bitnamicharts/postgresql \
  --version 16.7.27 \
  --namespace newton-db \
  -f deployments/k8s/postgresql/values.yaml

# Install PgBouncer
kubectl applyf --f deployments/k8s/postgresql/values.yaml
```

## Install Redis

```bash
kubectl apply -f deployments/k8s/redis
```

## Install MinIO

```bash
kubectl apply -f deployments/k8s/minio
```

## Install Observability

```bash
GRAFANA_INGRESS_DOMAIN=grafana.k8s.orb.local
kubectl create namespace observability || true

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

tee deployments/k8s/observability/values.yaml > /dev/null <<EOF
grafana:
  adminPassword: "supersecret"
  ingress:
    enabled: true
    ingressClassName: traefik
    hosts:
      - $GRAFANA_INGRESS_DOMAIN
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
  -f deployments/k8s/observability/values.yaml
```

## Install Application

## Build Docker Image

```bash
# Adjust your image name and version
IMAGE_NAME=apronny/fitbyte
IMAGE_VERSION=1.0.0-beta1
IMAGE_ARCH=linux/amd64
docker build --no-cache --platform=$IMAGE_ARCH -t $IMAGE_NAME:$IMAGE_VERSION .
docker push $IMAGE_NAME:$IMAGE_VERSION
```

## Deploy Docker Image to Kubernetes

Adjust your image name and version at `deployments/k8s/app/deployment.yaml`, adjust your application config at `deployments/k8s/app/configmap.yaml`, and don't forget to adjust your credential app config at `deployments/k8s/app/secret.yaml`.

Adjust your domain name at `deployments/k8s/app/ingress.yaml`.

Deploy the app with the command below:

```bash
kubectl apply -f deployments/k8s/app
```
