#!/usr/bin/env bash
# Script to fix kubectl configuration for the current user

set -e

echo "Fixing kubectl configuration for current user..."

# Get the correct API server address
API_SERVER=$(sudo kubectl config view --minify | grep server | awk '{print $2}')
echo "Found API server at: $API_SERVER"

# Create necessary directories
mkdir -p $HOME/.kube

# Copy the admin config from root
echo "Copying kubectl config from root..."
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
chmod 600 $HOME/.kube/config

echo "Kubectl configuration copied from root"
echo "Try running: kubectl get nodes"
