## Prerequisites

- Ensure you have `kubectl` installed and configured to interact with your Kubernetes cluster.
- Ensure `helm` is installed and configured if using Helm charts.

## Run

```bash
# Create Kubernetes Namespace
kubectl create ns newton-db

# Install Kubernetes Secret
kubectl create secret generic pg-fitbyte-secret \
  --from-literal=postgres-password=supersecret \
  --from-literal=password=supersecret \
  --from-literal=replication-password=supersecret \
  -n newton-db

# Install PostgreSQL with custom values
helm upgrade --install newton-pg oci://registry-1.docker.io/bitnamicharts/postgresql \
  --version 16.7.27 \
  --namespace newton-db \
  --set auth.existingSecret=pg-fitbyte-secret \
  --set auth.secretKeys.adminPasswordKey=postgres-password \
  --set auth.secretKeys.userPasswordKey=password \
  --set auth.secretKeys.replicationPasswordKey=replication-password \
  --set architecture=standalone \
  --set primary.resourcesPreset=none \
  --set primary.resources.requests.cpu=1000m \
  --set primary.resources.requests.memory=1024Mi \
  --set primary.resources.limit.memory=1024Mi
```
