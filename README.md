# 🎟️ Rifa

[Rifa](https://suerteconsarah.com) is a raffle platform built with a modern full-stack architecture in **Go** (backend) and **TypeScript/React** (frontend).  
It allows users to purchase raffle tickets, manage payments, and view results in a simple and transparent way.

---

## 🚀 Features

- ✨ Full-stack app with **Go + Huma + Chi** backend and **React + Vite** frontend
- 💳 Ticket purchasing & payment management
- 🐳 Dockerized for easy deployment
- ✅ Unit tests & GitHub Actions CI/CD

---

## 🛠️ Tech Stack

- **Backend:** Go (Chi, PostgreSQL)
- **Frontend:** React, TypeScript, Vite
- **Database:** PostgreSQL
- **CI/CD:** GitHub Actions + Railway deployment
- **Containerization:** Docker

---

## 📂 Repository Structure

```
.
├── backend/         # Go backend (Huma, Chi, PostgreSQL)
├── frontend/        # React + TypeScript frontend (Vite)
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
EMAIL_ACCOUNT=email@example.com
```

#### Frontend

```
VITE_API_URL=http://localhost:8080
```

---

## 🧪 Tests

Run unit tests for backend and frontend:

#### Backend

```bash
cd backend && go test ./...
```

#### Frontend

```bash
cd frontend && pnpm test
```

## 🚀 Deployment

Deployed on Railway with GitHub Actions CI/CD.
Each commit to main branch triggers a new deployment.

## 📜 License

This project is licensed under the **MIT License** – see the [LICENSE](./LICENSE) file for details.
