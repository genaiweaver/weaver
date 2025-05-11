#!/bin/bash

# Install Knative Eventing CRDs
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/eventing-crds.yaml

# Install Knative Eventing core
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/eventing-core.yaml

# Install in-memory channel provisioner
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/in-memory-channel.yaml

# Install MT Channel Based Broker
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/mt-channel-broker.yaml

# Create a default broker in the weaver namespace
cat <<EOF | kubectl apply -f -
apiVersion: eventing.knative.dev/v1
kind: Broker
metadata:
  name: default
  namespace: weaver
spec:
  # Using the MT Channel Based Broker
  config:
    apiVersion: v1
    kind: ConfigMap
    name: config-br-default-channel
    namespace: knative-eventing
EOF

echo "Knative Eventing has been installed."
echo "A default broker has been created in the weaver namespace."
