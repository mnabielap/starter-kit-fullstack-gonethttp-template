# 🚀 Starter Kit Fullstack - Go (net/http)

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/doc/devel/release.html#go1.22)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker&logoColor=white)](Dockerfile)
[![Swagger](https://img.shields.io/badge/Swagger-OpenAPI-85EA2D?logo=swagger&logoColor=black)](http://localhost:8080/swagger/index.html)

A production-ready **Fullstack Starter Kit** built with **Native Go** (`net/http`) and **GORM**.

This project implements a **Hybrid Architecture**:
1.  **RESTful JSON API**: Secured with JWT (Access & Refresh Tokens).
2.  **Server-Side Rendering**: HTML Templates served by Go, communicating with the API via a persistent JavaScript client.

It is designed as a robust alternative to PHP Native or Node.js boilerplates, offering type safety, high performance, and dual database support (SQLite & PostgreSQL).

---

## ✨ Features

- **🏗 Standard Go Layout**: Clean separation of concerns (`cmd`, `internal`, `pkg`).
- **💾 Dual Database Support**: Zero-config switch between **SQLite** (Pure Go) and **PostgreSQL**.
- **🔐 Secure Authentication**:
  - JWT Implementation (Access & Refresh Tokens).
  - CSRF Protection Middleware.
  - BCrypt Password Hashing.
- **🎨 Fullstack UI**:
  - **HTML/Templates**: Server-side rendered views (`web/templates`).
  - **JS Client**: Built-in `api-client.js` handles JWT storage and API fetching.
  - **Bootstrap 5**: Responsive dashboard UI.
- **🛡 Security**: Helmet-equivalent headers, Rate Limiting, and Input Validation.
- **🐳 Docker Ready**: Multi-stage build (Alpine Linux) with manual orchestration support.
- **📝 Swagger Docs**: Auto-generated API documentation.
- **🧪 Automated Testing**: Python-based script suite for endpoint verification (No Postman needed!).

---

## 📂 Project Structure

```text
starter-kit-fullstack-gonethttp-template/
├── cmd/server/            # Application entry point (main.go)
├── config/                # Configuration & Database connection
├── internal/
│   ├── handlers/          # HTTP Handlers (API & Web separated)
│   ├── middleware/        # Auth, CSRF, Logger, Rate Limiter
│   ├── models/            # GORM Database Models
│   ├── repository/        # Data Access Layer
│   ├── routes/            # Router & Middleware wiring
│   └── services/          # Business Logic
├── pkg/                   # Public Utilities (Response, View Engine)
├── web/
│   ├── static/            # CSS, JS (api-client.js), Images
│   └── templates/         # HTML Templates (Layouts, Partials, Pages)
├── api_tests/             # Python API Testing Scripts
├── migrations/            # SQL Migrations (Auto-migrates on startup)
├── .env.example           # Environment variables template
├── Dockerfile             # Docker build configuration
└── README.md              # Documentation
```

---

## 🛠️ Getting Started (Local Development)

**Recommended:** We suggest running the project locally first to understand the structure before containerizing it.

### Prerequisites
- **Go** (version 1.22 or higher)
- **Git**

### 1. Initialize Project
```bash
go mod tidy
```

### 2. Configuration
Copy the example environment file:
```bash
cp .env.example .env
```

**Configure `.env`:**
By default, it uses **SQLite**, so no extra setup is required.
```properties
APP_ENV=development
PORT=8080

# Switch to 'postgres' if needed
DB_DRIVER=sqlite
DB_NAME=starter_kit_db.sqlite
```

### 3. Run the Application
```bash
go run cmd/server/main.go
```

- **Web App**: Visit `http://localhost:8080`
- **Swagger UI**: Visit `http://localhost:8080/swagger/index.html`
- **Update Swagger Docs**: Run `swag init -g cmd/server/main.go -o docs`

---

## 🐳 Docker Deployment

We use manual Docker orchestration (without Compose) to simulate a real-world environment with persistent data.

### 1. Create Network 🌐
Create a dedicated network so the App and Database can communicate.
```bash
docker network create fullstack_gonethttp_network
```

### 2. Create Volumes 📦
Create persistent storage for the Database and Media/Static files.
```bash
# Volume for Postgres Data
docker volume create fullstack_gonethttp_postgres_data

# Volume for Application Media/Uploads
docker volume create fullstack_gonethttp_media_volume

# Volume for SQLite (Optional, if using SQLite in Docker)
docker volume create fullstack_gonethttp_db_volume
```

### 3. Start Database (PostgreSQL) 🐘
Run the Postgres container attached to our network and volume.
```bash
docker run -d \
  --name fullstack-gonethttp-postgres \
  --network fullstack_gonethttp_network \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -e POSTGRES_DB=starter_kit_db \
  -v fullstack_gonethttp_postgres_data:/var/lib/postgresql/data \
  postgres:15-alpine
```

### 4. Setup Environment for Docker 📝
Create a `.env.docker` file. **Crucial:** `DB_HOST` must match the Postgres container name.

```properties
APP_ENV=production
PORT=5005
APP_URL=http://localhost:5005

DB_DRIVER=postgres
DB_HOST=fullstack-gonethttp-postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=mysecretpassword
DB_NAME=starter_kit_db

JWT_SECRET=production_secure_secret
```

### 5. Build the App Image 🏗
```bash
docker build -t fullstack-gonethttp-app .
```

### 6. Run the App Container 🚀
Run the application container. We inject the environment file and mount volumes.

```bash
docker run -d -p 5005:5005 \
  --env-file .env.docker \
  --network fullstack_gonethttp_network \
  -v fullstack_gonethttp_db_volume:/app/db \
  -v fullstack_gonethttp_media_volume:/app/web/static \
  --name fullstack-gonethttp-container \
  fullstack-gonethttp-app
```

The application is now accessible at `http://localhost:5005`.

---

## 📦 Docker Management Cheat Sheet

Essential commands to manage your containers and volumes.

#### 📜 View Logs
See what's happening inside the container in real-time.
```bash
docker logs -f fullstack-gonethttp-container
```

#### 🛑 Stop Container
Safely stop the running application.
```bash
docker stop fullstack-gonethttp-container
```

#### ▶️ Start Container
Resume a stopped container (data persists).
```bash
docker start fullstack-gonethttp-container
```

#### 🗑 Remove Container
Remove the container instance (requires stopping first). Your data remains safe in the volumes.
```bash
docker rm fullstack-gonethttp-container
```

#### 📂 List Volumes
View all persistent storage volumes.
```bash
docker volume ls
```

#### ⚠️ Remove Volume
**WARNING:** This deletes your database data **permanently**!
```bash
docker volume rm fullstack_gonethttp_postgres_data
```

---

## 🧪 API Testing (Automated)

Forget Postman! This project comes with a suite of **Python scripts** to test every endpoint. These scripts automatically handle token storage (`secrets.json`), CSRF extraction, and request chaining.

### Setup
1. Ensure **Python 3.x** is installed.
2. Navigate to the `api_tests` directory.
3. **Important:** Edit `api_tests/utils.py` and ensure `BASE_URL` matches your running server:
   - Local: `"http://localhost:8080/v1"`
   - Docker: `"http://localhost:5005/v1"`

### How to Run
Run the scripts sequentially. No arguments needed.

**1. Authentication Flow:**
```bash
# Register a new user (Saves tokens automatically)
python api_tests/A1.auth_register.py

# Login (Refreshes tokens)
python api_tests/A2.auth_login.py

# Refresh Token Exchange
python api_tests/A3.auth_refresh.py
```

**2. User Management (Admin Role):**
*Note: Ensure you are logged in as an Admin (via `A2` using admin creds) to perform these actions.*

```bash
# Create a User manually
python api_tests/B1.user_create.py

# Get List of Users (Pagination & Filtering)
python api_tests/B2.user_get_list.py

# Update a User
python api_tests/B4.user_update.py

# Delete a User
python api_tests/B5.user_delete.py
```

---

## 📝 License

This project is licensed under the MIT License.