#!/usr/bin/env bash
# Simple script to install cri-dockerd using fixed version
# Based on https://www.mirantis.com/blog/how-to-install-cri-dockerd-and-migrate-nodes-from-dockershim/

set -e

echo "### Installing cri-dockerd..."

# Check if cri-dockerd is already installed
if command -v cri-dockerd &> /dev/null; then
  echo "cri-dockerd is already installed, skipping binary installation"
  INSTALLED_VERSION=$(cri-dockerd --version 2>&1 | grep -o "v[0-9]*\.[0-9]*\.[0-9]*" || echo "unknown")
  echo "Installed version: ${INSTALLED_VERSION}"
else
  # Use a fixed version that is known to work
  VERSION="0.3.17"
  echo "Using version: v${VERSION}"

  # Download the binary
  echo "Downloading cri-dockerd..."
  wget https://github.com/Mirantis/cri-dockerd/releases/download/v${VERSION}/cri-dockerd-${VERSION}.amd64.tgz

  # Extract the binary
  echo "Extracting cri-dockerd..."
  tar xvf cri-dockerd-${VERSION}.amd64.tgz

  # Move the binary to /usr/local/bin
  echo "Installing cri-dockerd binary..."
  install -o root -g root -m 0755 cri-dockerd/cri-dockerd /usr/local/bin/cri-dockerd

  # Clean up
  rm -rf cri-dockerd cri-dockerd-${VERSION}.amd64.tgz
fi

# Check if service files already exist
if [ -f "/etc/systemd/system/cri-docker.service" ] && [ -f "/etc/systemd/system/cri-docker.socket" ]; then
  echo "Service files already exist, skipping download"
else
  echo "Setting up systemd service files..."
  # Download the systemd service and socket files
  wget https://raw.githubusercontent.com/Mirantis/cri-dockerd/master/packaging/systemd/cri-docker.service
  wget https://raw.githubusercontent.com/Mirantis/cri-dockerd/master/packaging/systemd/cri-docker.socket

  # Install the service files
  install -m 0644 cri-docker.service /etc/systemd/system/
  install -m 0644 cri-docker.socket /etc/systemd/system/

  # Update the service file to point to the correct binary location
  sed -i -e 's,/usr/bin/cri-dockerd,/usr/local/bin/cri-dockerd,' /etc/systemd/system/cri-docker.service

  # Clean up downloaded files
  rm -f cri-docker.service cri-docker.socket
fi

# Check if service is already running
if systemctl is-active --quiet cri-docker.socket; then
  echo "cri-docker.socket is already running"
else
  # Reload systemd and enable the service
  echo "Enabling and starting cri-docker service..."
  systemctl daemon-reload
  systemctl enable --now cri-docker.socket
fi

echo "cri-dockerd installation complete!"
echo "To use with minikube, run: minikube start --container-runtime=docker --cri-socket=/var/run/cri-dockerd.sock"
