# RentalFlow

A distributed microservices-based rental management platform designed for scalability and performance. Built with Go, gRPC, and RabbitMQ.

## 🏗 System Architecture

RentalFlow adopts a microservices architecture with 7 independent services, using **gRPC** for synchronous communication and **RabbitMQ** for asynchronous event handling.

### Tech Stack
- **Languages**: Go (Golang) 1.22+
- **Communication**: gRPC (Internal), REST (Gateway), RabbitMQ (Async Events)
- **Database**: PostgreSQL 15 (Per-service isolation)
- **Caching**: Redis 7
- **Deployment**: Docker, Render (Modular Monolith capability)

### Microservices Overview

| Service | Port (HTTP/gRPC) | Description |
|---------|-------------------|-------------|
| **API Gateway** | `8080` | Entry point, REST-to-gRPC transcoding, Auth middleware |
| **Auth Service** | `8081` / `50051` | User registration, JWT issuance, Role-based access |
| **Inventory Service** | `8082` / `50052` | Equipment/Property management, Stock tracking |
| **Booking Service** | `8083` / `50053` | Reservation logic, Availability checks |
| **Payment Service** | `8084` / `50054` | Payment processing integration (Chapa) |
| **Notification Service** | `8085` / `50055` | Email/SMS alerts via RabbitMQ events |
| **Review Service** | `8086` / `50056` | Ratings and feedback management |

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local development)
- Make (optional)

### Quick Start (Local Docker)
Run the entire system with a single command:

```bash
# Start all services, databases, and message broker
docker-compose up --build
```

The API Gateway will be available at `http://localhost:8080`.

### Running Tests
We emphasize rigorous testing with unit and integration suites.

```bash
# Run all unit tests
make test

# Run integration tests
./scripts/test_integration.sh
```

## 📦 Deployment

The project is configured for **Render**. We support a "Modular Monolith" deployment to save costs by running all services in a single container while maintaining logical separation.

- **Unified Dockerfile**: `Dockerfile.render`
- **Render Config**: `render.yaml`

## 📚 Documentation

- **API Documentation**: [View API Specification (OpenAPI)](./docs/openapi.yaml)
- **Project Report**: [Link to PDF]
- **Progress Logs**: See `docs/progress.md`
- **Postman Testing Collections**: ./tests

## 📄 License
MIT License.
