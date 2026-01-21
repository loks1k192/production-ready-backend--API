# Deployment

## Docker
Build and run:
```bash
docker build -t go-backend:latest .
docker run -p 8080:8080 \
  -e DB_DSN="postgres://postgres:postgres@host.docker.internal:5432/app?sslmode=disable" \
  -e AUTH_JWT_SECRET="change-me" \
  go-backend:latest
```

## Kubernetes
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

Update `DB_DSN` and `AUTH_JWT_SECRET` in the manifests before deploying.
