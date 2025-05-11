#!/bin/bash

# Get the broker URL
BROKER_URL=$(kubectl get broker default -n weaver -o jsonpath='{.status.address.url}')

echo "Sending event to broker at: $BROKER_URL"

# Send the "start" event
curl -v -X POST \
  -H "Ce-Id: test-1" \
  -H "Ce-SpecVersion: 1.0" \
  -H "Ce-Type: com.sample.start" \
  -H "Ce-Source: /tests" \
  -H "Content-Type: application/json" \
  -d '{"foo":"bar"}' \
  "$BROKER_URL"
