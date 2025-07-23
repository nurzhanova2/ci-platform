# CI Platform (DevOps Pet Project)

This is a ***DevOps pet project*** — a minimal CI/CD platform written in Go and running inside Docker. It’s designed for learning how real-world CI systems work internally: job orchestration, container builds, webhook triggers, and configuration management.


## Tech Stack

- **Go** 1.21+
- **Docker CLI** inside Alpine
- **Viper** for configuration
- **Gorilla Mux** for HTTP routing
- RESTful API

## Features

- Trigger pipelines from HTTP (e.g., webhook)
- Execute simple build/test jobs inside containers
- Containerized with Docker Compose
- Modular project structure (Go best practices)
- Easily extendable

## How to Run


```bash
// clone the Repository
git clone https://github.com/nurzhanova2/ci-platform.git
cd ci-platform

// Build & Run with Docker Compose
docker compose up --build
