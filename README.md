# 🎟️ Rifa

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)]()
[![TypeScript](https://img.shields.io/badge/TypeScript-5.5+-3178C6?logo=typescript)]()
[![React](https://img.shields.io/badge/React-19.0+-61DAFB?logo=react)]()
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](./LICENSE)

[Rifa](https://suerteconsarah.com) is a raffle platform built with a modern full-stack architecture in **Go** (backend), **TypeScript/React** (frontend) and **Prometheus/Grafana** (monitoring).
It allows users to purchase raffle tickets, manage payments, and view results in a simple and transparent way.

---

## 🚀 Features

- ✨ Full-stack app with **Go + Huma + Chi** backend and **React + Vite + TailwindCss** frontend
- 💳 Ticket purchasing & payment management
- 🐳 Dockerized for easy deployment
- ✅ Unit and Integration tests & GitHub Actions CI/CD
- 📊 OpenTelemetry tracing and Prometheus metrics
- 🪶 Idempotent request middleware
- 🧩 Background job system for async processing
- 🧾 Structured logger and graceful shutdown

---

## 🛠️ Tech Stack

- **Backend:** Go (Huma, Chi, Pgx)
- **Frontend:** React, TypeScript, Vite, TailwindCSS
- **Database:** PostgreSQL
- **CI/CD:** GitHub Actions + Railway deployment
- **Containerization:** Docker

---

## 🧰 Prerequisites

- Go **1.25+**
- Node **18+** / pnpm **9+**
- Docker & Docker Compose
- PostgreSQL **14+** (or use Docker)
- (Optional) Grafana/Prometheus stack (see Monitoring)

---

## 📂 Repository Structure

```
.
├── backend/         # Go backend (Huma, Chi, Pgx)
├── frontend/        # React (Vite) + TypeScript + TailwindCSS
├── monitoring/      # Otel collector + Prometheus + Loki + Tempo + Grafana configuration
├── .github/         # GitHub Actions workflows
├── Dockerfile       # Containerization setup
├── LICENSE
└── README.md
```

---

## ⚡ Getting Started

### 1️⃣ Clone the repo

```bash
git clone https://github.com/cerjey13/rifa.git
cd rifa
```

### 2️⃣ Run with Docker

```bash
docker build -t rifa .
docker run -p 3000:3000 rifa
```

### 3️⃣ Local Development

#### Backend:

```bash
cd backend
go run cmd/app/main.go
```

#### Frontend:

```bash
cd frontend
pnpm install
pnpm dev
```

For embedding into backend binary (backend/cmd/app/dist):

```bash
pnpm build-back
```

---

## 🔑 Environment Variables (example)

These values are examples, adjust to your actual config.

#### Backend

```
DATABASE_URL=postgres://user:pass@localhost:5432/rifa?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=choose-a-strong-secret
COOKIE_SECURE=true
ENV=development
EMAIL_MAILEROO_API_KEY=maileroo-provider-api-key
EMAIL_URL=serviceprovider.smtp.email/api/
EMAIL_SENDER_ACCOUNT=emailsender@provider.com
EMAIL_ACCOUNT=email@example.com
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

#### Frontend

```
VITE_API_URL=http://localhost:8080
```

---

## 🧪 Tests

Run tests for backend and frontend:

### Backend

#### Unit Tests

```bash
cd backend && go test ./...
```

#### Integration Tests

```bash
cd backend && go test -v -tags=integration ./...
```

### Frontend

```bash
cd frontend && pnpm test
```

---

## 🧩 Monitoring

1. **Start Monitoring Stack**
   ```bash
   cd monitoring
   docker compose up -d
   ```
2. **OTel Collector on :4318 (HTTP)**
3. **Prometheus scrapes your backend on port 4318**
4. **Grafana available at localhost:3000 (login: admin/admin)**
5. **Loki/Tempo for logs/traces**
6. **The backend exports traces/metrics when OTEL_EXPORTER_OTLP_ENDPOINT is set**

---

## 🚀 Deployment

Deployed on Railway with GitHub Actions CI/CD.
Each commit to main branch triggers a new deployment.

## 📜 License

This project is licensed under the **MIT License** – see the [LICENSE](./LICENSE) file for details.
