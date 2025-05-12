# Weaver

Weaver is a Kubernetes-native **protocol specification** for building production-grade, event-driven GenAI applications via Knative. It defines CloudEvent schemas, JSON-LD service & agent descriptors, discovery & versioning rules, and a CNCF-inspired governance model—separating the **“what”** (the protocol) from the **“how”** (implementations like **Genie**).

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
4. [Specification Overview](#specification-overview)
5. [Usage](#usage)
6. [Contributing](#contributing)
7. [Roadmap](#roadmap)
8. [License](#license)
9. [Acknowledgments](#acknowledgments)

## Introduction

Weaver provides an open, implementation-neutral **protocol**—inspired by MCP and A2A—that standardizes how multi-step AI agent components discover, describe, and wire together on Kubernetes via Knative Eventing & Serving. It empowers multiple independent implementations (like **Genie**) to interoperate on a shared foundation.

## Features

- **CloudEvent-First**: Defines a core CloudEvent v1.0 schema with extension points for AI workflows.
- **JSON-LD Descriptors**: Service and Agent descriptor contexts for capabilities, discovery, and versioning.
- **Knative-Native**: Maps high-level graphs to Knative `Service` & `Trigger` CRDs for auto-wiring event flows.
- **Governance Model**: CNCF-style sandbox → incubation → graduation stages, TOC, and Code of Conduct.
- **Extensible Hooks**: Support for error, retry, and A2A extension metadata via descriptor “hooks.”

## Getting Started

### Prerequisites

- **Kubernetes** v1.26+ cluster
- **Knative Serving & Eventing** installed on your cluster with a broker
- **kubectl** CLI tool
- **kn** (Knative CLI) for easier interaction with Knative resources
- **docker** for building and pushing container images
- **make** for building and deploying the controller

<!-- #### Setting up Knative

If you don't have Knative installed, you can follow these steps:

1. **Install Knative Serving**:
   ```bash
   kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-crds.yaml
   kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-core.yaml
   ```

2. **Install Knative Eventing**:
   ```bash
   kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/eventing-crds.yaml
   kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/eventing-core.yaml
   ```

3. **Install a networking layer** (e.g., Kourier):
   ```bash
   kubectl apply -f https://github.com/knative/net-kourier/releases/download/knative-v1.12.0/kourier.yaml
   kubectl patch configmap/config-network \
     --namespace knative-serving \
     --type merge \
     --patch '{"data":{"ingress-class":"kourier.ingress.networking.knative.dev"}}'
   ```

4. **Install the MT Channel Based Broker**:
   ```bash
   kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.12.0/mt-channel-broker.yaml
   ``` -->

### Installation

1. **Clone the repo**
   ```bash
   git clone https://github.com/weaver/weaver.git
   cd weaver
   ```

2. **Review the spec** in `SPEC.md` and the JSON-LD contexts under `schemas/`.

3. **Install Knative Serving and Eventing**
   Ensure you have Knative Serving and Eventing installed in your Kubernetes cluster. If not, follow the [Knative installation guide](https://knative.dev/docs/install/).

4. **Install the Weaver Controller**
   ```bash
   # Build and deploy the controller
   make deploy
   ```

5. **Create a namespace for Weaver resources**
   ```bash
   kubectl create namespace weaver
   ```

6. **Create a Knative Broker**
   ```bash
   kubectl apply -f config/samples/default-broker.yaml
   ```

7. **Deploy the sample services**
   ```bash
   # Deploy the Knative Services for the planner and echo components
   kubectl apply -f config/samples/planner-ksvc.yaml
   kubectl apply -f config/samples/echo-ksvc.yaml
   ```

8. **Deploy the WeaverAgent resource**
   ```bash
   kubectl apply -f config/samples/sample-agent.yaml
   ```

## Architecture

Weaver uses a Kubernetes controller pattern to manage event-driven AI applications:

1. **WeaverAgent CRD**: Defines a high-level graph of services and their event connections
2. **Weaver Controller**: Watches for WeaverAgent resources and reconciles them into Knative resources
3. **Knative Services**: Each node in the graph becomes a Knative Service
4. **Knative Triggers**: Each edge in the graph becomes a Knative Trigger that routes events
5. **Knative Broker**: Central event bus that receives and routes events based on Triggers

![Weaver Architecture](https://raw.githubusercontent.com/weaver/weaver/main/docs/images/architecture.png)

### Controller Behavior

The Weaver controller:
- Watches for changes to WeaverAgent resources
- Creates or updates Knative Services for each node in the graph
- Creates or updates Knative Triggers for each edge in the graph
- Sets owner references so that deleting a WeaverAgent cleans up all related resources
- Updates the status of the WeaverAgent to reflect the current state

## Specification Overview

See **`SPEC.md`** for full details. In brief:

- **CloudEvent core**: `id`, `source`, `type`, `time`, `data` with extension fields like `mcp.action` and `correlationId`.
- **Service Descriptor** (`schemas/service-schema.json`): Declares `id`, `version`, `endpoint`, `events.consumes`/`produces`, and `capabilities`.
- **Agent Descriptor** (`schemas/agent-schema.json`): Graph (`nodes`, `edges`, `hooks`) that the Weaver controller reconciles into Knative CRDs.

## Usage

### Creating and Deploying a WeaverAgent

1. **Write a `WeaverAgent`** CRD in YAML, defining your agent graph.
   ```yaml
   apiVersion: weaver.io/v1alpha1
   kind: WeaverAgent
   metadata:
     name: my-agent
     namespace: weaver
   spec:
     broker: default
     nodes:
       - serviceName: my-service-1
       - serviceName: my-service-2
     edges:
       - eventType: com.example.event1
         to: my-service-1
       - eventType: com.example.event2
         to: my-service-2
   ```

2. **Apply** it:
   ```bash
   kubectl apply -f my-agent.yaml
   ```

### Verifying the Deployment

3. **Verify** the Weaver controller is running:
   ```bash
   kubectl get pods -n weaver-system
   ```

4. **Check** the Knative Services created by the controller:
   ```bash
   kubectl get ksvc -n weaver
   # or using the Knative CLI
   kn service list -n weaver
   ```

5. **Verify** the Knative Triggers created by the controller:
   ```bash
   kubectl get triggers -n weaver
   ```

### Testing the Event Flow

6. **Emit** a CloudEvent to the Broker to test the event flow:
   ```bash
   # Example using curl to send an event to the broker
   curl -v "http://broker-ingress.knative-eventing.svc.cluster.local/weaver/default" \
     -H "Ce-Id: 123456789" \
     -H "Ce-Specversion: 1.0" \
     -H "Ce-Type: com.sample.start" \
     -H "Ce-Source: curl.client" \
     -H "Content-Type: application/json" \
     -d '{"message": "Hello, Weaver!"}'
   ```

7. **Check** the logs of your services to see the event flow:
   ```bash
   # For the planner service
   kubectl logs -n weaver -l serving.knative.dev/service=weaver-planner

   # For the echo service
   kubectl logs -n weaver -l serving.knative.dev/service=weaver-echo
   ```

### Troubleshooting

If you encounter issues with the controller, check the logs:
```bash
kubectl logs -n weaver-system -l app.kubernetes.io/name=weaver
```

Common issues:
- Make sure both Knative Serving and Eventing are properly installed
- Verify that the Knative Broker exists in the correct namespace
- Check that the service images are accessible and can be pulled

## Contributing

Weaver welcomes community contributions! Please see **`CONTRIBUTING.md`**:

- **Fork** the repo and create feature branches.
- **Submit PRs** against `main`; require two maintainer approvals and CI green checks (JSON Schema, linting).
- **Report issues** via GitHub Issues, labeled for triage.

## Roadmap

- **v1.0**: Core event schemas, JSON-LD contexts, JSON Schemas, basic CRD & controller.
- **v1.1**: A2A protocol extensions, prompt registry spec.
- **v2.0**: Multi-broker support, advanced security & compliance profiles.

## License

Apache License 2.0 © 2025 Weaver Contributors. See **`LICENSE`** for details.

## Acknowledgments

- Inspired by **Model Context Protocol (MCP)** and **Knative**.
- Governance patterns adapted from CNCF project practices.
