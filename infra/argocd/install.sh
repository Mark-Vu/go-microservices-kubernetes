#!/bin/bash
set -e

echo "=========================================="
echo "Installing ArgoCD on GKE"
echo "=========================================="

# Create namespace
echo "1. Creating argocd namespace..."
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -

# Install ArgoCD
echo "2. Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "3. Waiting for ArgoCD pods to be ready (this may take 2-3 minutes)..."
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=argocd-server -n argocd --timeout=300s

# Change the service type to LoadBalancer to access ArgoCD UI from outside the cluster, default is ClusterIP
kubectl patch svc argocd-server -n argocd -p '{"spec":{"type":"LoadBalancer"}}'

echo ""
echo "=========================================="
echo "ArgoCD Installation Complete!"
echo "=========================================="
echo ""
echo "Getting ArgoCD admin password..."
ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

echo ""
echo "ArgoCD Admin Credentials:"
echo "  Username: admin"
echo "  Password: $ARGOCD_PASSWORD"
echo ""
kubectl get svc argocd-server -n argocd --watch

echo ""
echo "Access ArgoCD UI at: https://<EXTERNAL-IP>"
echo ""
echo "Next steps:"
echo "1. Login to ArgoCD UI"
echo "2. Deploy application: kubectl apply -f infra/argocd/application.yaml"
