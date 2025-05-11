#!/bin/bash

# Get the WeaverAgent
AGENT_NAME="sample-agent"
AGENT_NAMESPACE="weaver"

# Get the broker name
BROKER=$(kubectl get weaveragent -n $AGENT_NAMESPACE $AGENT_NAME -o jsonpath='{.spec.broker}')

# Get the edges
EDGES=$(kubectl get weaveragent -n $AGENT_NAMESPACE $AGENT_NAME -o jsonpath='{.spec.edges}')

# Parse the edges and create triggers
echo "Creating triggers for WeaverAgent $AGENT_NAME in namespace $AGENT_NAMESPACE"
echo "Using broker: $BROKER"

# Get the number of edges
NUM_EDGES=$(kubectl get weaveragent -n $AGENT_NAMESPACE $AGENT_NAME -o jsonpath='{.spec.edges[*].eventType}' | wc -w)

# Loop through the edges
for i in $(seq 0 $((NUM_EDGES-1))); do
  EVENT_TYPE=$(kubectl get weaveragent -n $AGENT_NAMESPACE $AGENT_NAME -o jsonpath="{.spec.edges[$i].eventType}")
  TO=$(kubectl get weaveragent -n $AGENT_NAMESPACE $AGENT_NAME -o jsonpath="{.spec.edges[$i].to}")
  
  # Create the trigger
  echo "Creating trigger for event type $EVENT_TYPE to service $TO"
  
  cat <<EOF | kubectl apply -f -
apiVersion: eventing.knative.dev/v1
kind: Trigger
metadata:
  name: $AGENT_NAME-edge-$i
  namespace: $AGENT_NAMESPACE
spec:
  broker: $BROKER
  filter:
    attributes:
      type: $EVENT_TYPE
  subscriber:
    ref:
      apiVersion: serving.knative.dev/v1
      kind: Service
      name: $TO
EOF
done

echo "Done creating triggers"
