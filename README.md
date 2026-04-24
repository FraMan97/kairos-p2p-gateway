# Kairos P2P Gateway

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=for-the-badge&logo=kubernetes&logoColor=white)
![Helm](https://img.shields.io/badge/Helm-0F162D?style=for-the-badge&logo=helm&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-GNU_AGPL_v3-blue?style=for-the-badge&logo=gnu&logoColor=white)

The **Kairos P2P Gateway** is the official headless gateway for the **kairos-p2p-engine** ecosystem. It provides a unified, simplified REST API to interact with the decentralized Kairos storage network, handling file orchestration, cryptographic sealing, and network monitoring.

> [!WARNING]
> This project is a Proof of Concept (POC) focused on backend infrastructure and is not intended for production environments without further security audits.

---

## Overview

It acts as a high-performance proxy between client applications (B2C/B2B) and the underlying P2P infrastructure. It abstracts the complexity of shard distribution, Reed-Solomon erasure coding, and Time-Lock encryption provided by the engine.

### Key Features
* **Proxy Streaming**: Efficiently forwards file data to storage nodes using `io.Pipe`, ensuring minimal RAM usage even for large file transfers.
* **Network Metrics Aggregator**: Proxies global network status (active nodes, secured files, storage capacity) from the Kairos Explorer.
* **Agnostic Design**: Serves as a modular interface that can be deployed independently of the core storage nodes.
* **Cloud-Native**: Fully containerized and ready for Kubernetes deployment via Helm.

## API Endpoints

The gateway implements and proxies the following endpoints:

### File Management
* `POST /put`: Uploads and seals a file. Handles multipart form data and forwards it to the Kairos network.
* `GET /get?id=<file_id>`: Retrieves and reconstructs a file from shards (only accessible after the release time).
* `DELETE /delete?id=<file_id>`: Removes the file manifest and all associated shards from the network.
* `GET /upload/status?id=<file_id>`: Checks the asynchronous processing status of an uploaded file.

### Monitoring
* `GET /metrics`: Returns aggregated network health and storage metrics by querying the Kairos Explorer.

## Configuration

The gateway is configured via environment variables on values.yaml helm file:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Local port the gateway listens on | `3000` |
| `KAIROS_NETWORK_URL` | URL of the Kairos Engine LoadBalancer/API | `http://kairos-engine-api:80` |
| `KAIROS_EXPLORER_URL` | URL of the Kairos Explorer service | `http://kairos-explorer:8081` |

## Getting Started

### Prerequisites
* **Go**: 1.24+
* **Kairos P2P Engine**: A running instance of the engine (Bootstrap, Nodes, and Explorer).

### Docker Deploy
```bash
    ./deploy-dev.sh
```