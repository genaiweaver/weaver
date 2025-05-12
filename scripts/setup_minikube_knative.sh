#!/usr/bin/env bash
# Simplified Minikube setup with cri-dockerd
# Usage: sudo ./setup_minikube_knative.sh

set -e

# Ensure running as root
if [[ "$EUID" -ne 0 ]]; then
  echo "Please run this script as root or with sudo."
  exit 1
fi

echo "### 1. Install prerequisites"
apt-get update -y
apt-get install -y curl wget apt-transport-https ca-certificates software-properties-common conntrack

echo "### 2. Install Docker (if not already installed)"
if ! command -v docker &> /dev/null; then
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
  add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
  apt-get update -y
  apt-get install -y docker-ce docker-ce-cli containerd.io
  usermod -aG docker "$SUDO_USER" || true
  echo "Added $SUDO_USER to docker group. Logout/login for group to take effect."
fi

echo "### 3. Install kubectl (if not already installed)"
if ! command -v kubectl &> /dev/null; then
  curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
  install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
  rm kubectl
fi

echo "### 4. Install cri-dockerd"
# Run the dedicated script to install cri-dockerd
./install_cri_dockerd.sh

echo "### 5. Install Minikube (if not already installed)"
if ! command -v minikube &> /dev/null; then
  curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
  install minikube /usr/local/bin/
  rm minikube
fi

echo "### 6. Start Minikube with cri-dockerd"
sysctl -w fs.protected_regular=0

# Create symbolic links to ensure socket is found in the expected location
echo "Checking socket paths..."
mkdir -p /var/run

# Check if the socket exists
if [ ! -S "/run/cri-dockerd.sock" ]; then
  echo "Warning: Socket /run/cri-dockerd.sock does not exist. Is cri-dockerd running?"
  echo "Checking service status..."
  systemctl status cri-docker.socket --no-pager
else
  echo "Socket /run/cri-dockerd.sock exists"
fi

# Create symbolic links if they don't already exist
if [ ! -e "/var/run/cri-docker.sock" ]; then
  echo "Creating symbolic link for cri-docker.sock..."
  ln -sf /run/cri-dockerd.sock /var/run/cri-docker.sock
else
  echo "Symbolic link /var/run/cri-docker.sock already exists"
fi

if [ ! -e "/var/run/cri-dockerd.sock" ]; then
  echo "Creating symbolic link for cri-dockerd.sock..."
  ln -sf /run/cri-dockerd.sock /var/run/cri-dockerd.sock
else
  echo "Symbolic link /var/run/cri-dockerd.sock already exists"
fi

# Check if minikube is already running
MINIKUBE_STATUS=$(minikube status --format={{.Host}} 2>/dev/null || echo "Not Running")
if [ "$MINIKUBE_STATUS" = "Running" ]; then
  echo "Minikube is already running. Stopping it first..."
  minikube stop
fi

# Delete any existing minikube instance
echo "Deleting any existing minikube instance..."
minikube delete || true

# Start minikube with the correct socket path
echo "Starting minikube with cri-dockerd socket..."
minikube start --driver=none --container-runtime=docker --cri-socket=/var/run/cri-dockerd.sock

echo "### 7. Configure kubectl for current user"
# Copy kubectl configuration to current user if running as root
if [ "$EUID" -eq 0 ] && [ -n "$SUDO_USER" ]; then
  echo "Configuring kubectl for user $SUDO_USER..."
  USER_HOME=$(getent passwd "$SUDO_USER" | cut -d: -f6)
  mkdir -p "$USER_HOME/.kube"

  # Copy the admin config
  if [ -f "/etc/kubernetes/admin.conf" ]; then
    cp /etc/kubernetes/admin.conf "$USER_HOME/.kube/config"
    chown "$SUDO_USER:$(id -gn $SUDO_USER)" "$USER_HOME/.kube/config"
    chmod 600 "$USER_HOME/.kube/config"
    echo "Kubectl configured for user $SUDO_USER using admin.conf"
  else
    echo "Warning: admin.conf not found, kubectl may not work for user $SUDO_USER"
  fi

  # Create a helper script for the user
  cat > "$USER_HOME/fix_kubectl_config.sh" << 'EOF'
#!/usr/bin/env bash
# Script to fix kubectl configuration for the current user

set -e

echo "Fixing kubectl configuration for current user..."

# Create necessary directories
mkdir -p $HOME/.kube

# Copy the admin config from root
echo "Copying kubectl config from root..."
sudo cp /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
chmod 600 $HOME/.kube/config

echo "Kubectl configuration copied from root"
echo "Try running: kubectl get nodes"
EOF

  chmod +x "$USER_HOME/fix_kubectl_config.sh"
  chown "$SUDO_USER:$(id -gn $SUDO_USER)" "$USER_HOME/fix_kubectl_config.sh"
  echo "Created fix_kubectl_config.sh script for user $SUDO_USER"
fi

echo "### 8. Verify installation"
kubectl get nodes
minikube status

# Wait for node to be ready
echo "Waiting for node to be ready..."
kubectl wait --for=condition=ready node/$(hostname) --timeout=120s || echo "Node not ready yet, you may need to wait longer"

echo "Setup complete!"
echo "You can now use kubectl with your cluster."
echo "To check the status of your cluster, run: kubectl get nodes"
echo "To check running pods, run: kubectl get pods -A"
echo ""
echo "If you have issues with kubectl permissions, run: ./fix_kubectl_config.sh"
