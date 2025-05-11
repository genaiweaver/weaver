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

- **Kubernetes** v1.26+ with **Knative Serving & Eventing** installed.  
- **kubectl**, **kn** CLI, and **docker** (for building examples).

### Installation

1. **Clone the repo**  
   ```bash
   git clone https://github.com/weaver/weaver.git
   cd weaver
   ```
2. **Review the spec** in `SPEC.md` and the JSON-LD contexts under `schemas/`.  
3. **Install the CRD** (optional)  
   ```bash
   kubectl apply -f crds/agent-crd.yaml
   ```
4. **Deploy examples**  
   ```bash
   kubectl apply -f examples/simple-agent-cm.yaml
   ```

## Specification Overview

See **`SPEC.md`** for full details. In brief:

- **CloudEvent core**: `id`, `source`, `type`, `time`, `data` with extension fields like `mcp.action` and `correlationId`.  
- **Service Descriptor** (`schemas/service-schema.json`): Declares `id`, `version`, `endpoint`, `events.consumes`/`produces`, and `capabilities`.  
- **Agent Descriptor** (`schemas/agent-schema.json`): Graph (`nodes`, `edges`, `hooks`) that the Weaver controller reconciles into Knative CRDs.

## Usage

1. **Write a `WeaverAgent`** CRD or ConfigMap in YAML, defining your agent graph.  
2. **Apply** it:  
   ```bash
   kubectl apply -f examples/langgraph-agent.yaml
   ```
3. **Verify** the Knative resources:  
   ```bash
   kn service list -n weaver
   kubectl get triggers -n weaver
   ```
4. **Emit** a CloudEvent to the Broker (e.g. `com.user.chat`) via your Chat-Adapter service.

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
