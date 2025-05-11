#!/bin/bash

# Start port-forwarding in the background
kubectl port-forward -n knative-eventing svc/broker-ingress 8080:80 &
PF_PID=$!

# Wait for port-forwarding to be established
sleep 2

# Send the event
curl -v -X POST \
  -H "Ce-Id: test-2" \
  -H "Ce-SpecVersion: 1.0" \
  -H "Ce-Type: com.sample.planned" \
  -H "Ce-Source: /tests" \
  -H "Content-Type: application/json" \
  -d '{"plannedFrom":{"foo":"bar"}}' \
  "http://localhost:8080/weaver/default"

# Kill the port-forwarding process
kill $PF_PID
