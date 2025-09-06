## 🚀 **Quick Deploy**

```bash
helm repo add traefik https://traefik.github.io/charts

helm upgrade --install traefik traefik/traefik \
  --version 37.1.0 \
  -f values.yaml \
  -n traefik-system \
  --create-namespace
```

## 📊 **Verify Metrics**

```bash
# Port forward metrics
kubectl port-forward svc/traefik-metrics 9100:9100 -n traefik-system

# Test metrics endpoint
curl http://localhost:9100/metrics | head -20

# Check specific Golang app metrics
curl http://localhost:9100/metrics | grep traefik_service_request
```

## **IngressRoute**

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: fitbyte
spec:
  entryPoints:
    - web
  routes:
    - match: Host(`your-app.local`)
      kind: Rule
      services:
        - name: fitbyte-service
          port: 8080
```
