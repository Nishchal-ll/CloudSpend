# ☁️ CloudSpend — Cloud Expense & Budget Intelligence Platform

[![Live Demo](https://img.shields.io/badge/Live%20Demo-Azure%20App%20Service-0078D4?style=for-the-badge&logo=microsoftazure&logoColor=white)](https://cloudspend-app-a9a3bxcce7hye8hv.eastasia-01.azurewebsites.net/)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![Database](https://img.shields.io/badge/Database-SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![Docker](https://img.shields.io/badge/Container-Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

> **Live Production URL:** [https://cloudspend-app-a9a3bxcce7hye8hv.eastasia-01.azurewebsites.net/](https://cloudspend-app-a9a3bxcce7hye8hv.eastasia-01.azurewebsites.net/)  
> **Demo Credentials:** `admin@example.com` / `password123`

---

## 📸 Application Screenshots

### 1. Landing Page
![CloudSpend Landing Page](https://raw.githubusercontent.com/Nishchal-ll/CloudSpend/main/image1.png)

### 2. Executive Dashboard (100% Full-Width)
![CloudSpend Dashboard](https://raw.githubusercontent.com/Nishchal-ll/CloudSpend/main/image2.png)

---

## 📌 Project Overview

**CloudSpend** is an enterprise-grade cloud cost intelligence and budget management platform built with **Go (Gin)**, **SQLite**, and **Go HTML Templates**, containerized with **Docker**, and deployed on **Microsoft Azure App Service** via **Azure Container Registry (ACR)**.

It empowers engineering teams and cloud architects to monitor multi-cloud infrastructure expenditures across **Microsoft Azure, AWS, Google Cloud, Cloudflare, and SaaS services**, calculate real-time budget utilization, analyze spend by categories, and prevent cost overruns.

---

## 🌟 Key Features

* **Real-time Cost Intelligence Dashboard:**
  * **3 Summary KPI Cards:** Total Monthly Spend, Budget Consumption Percentage, and Active Cloud Providers.
  * **Category Breakdown:** Real-time distribution across *Compute, Database, Storage, Networking, Containers, Monitoring, and DevOps*.
  * **Provider Distribution:** SVG Donut chart displaying cloud expenditure allocation.
  * **Recent Transactions:** Quick-access audit log of recent cloud expenses.
* **Full CRUD Expenses Management (`/expenses`):**
  * Search, filter by category, add new resources, edit entries, and delete obsolete infrastructure records.
* **Authentication & Role Security:**
  * Bcrypt-hashed password authentication with protected session middleware.
* **Target Budget Management:**
  * Configurable monthly budget limits with automatic threshold health status badges (*Healthy*, *Warning*, *Over Limit*).
* **REST API Surface (`/api/*`):**
  * JSON endpoints for automation, third-party integrations, and health monitoring (`/healthz`).
* **Ultra-Lightweight Production Container:**
  * Multi-stage Docker build producing a secure Alpine Linux image (< 25MB).

---

## 🏗️ Architecture & Technology Stack

```text
┌─────────────────────────────────────────────────────────────┐
│                 Microsoft Azure Cloud                       │
│                                                             │
│   ┌────────────────────────┐      ┌─────────────────────┐   │
│   │ Azure Container        │      │   Azure App Service │   │
│   │ Registry (ACR)         ├─────►│   (Linux Container) │   │
│   │ cloudspend.azurecr.io  │      │   :8080 (HTTPS)     │   │
│   └────────────────────────┘      └──────────┬──────────┘   │
└──────────────────────────────────────────────┼──────────────┘
                                               │
                                               ▼
┌─────────────────────────────────────────────────────────────┐
│                    CloudSpend Container                      │
│                                                             │
│   ┌────────────────────┐      ┌─────────────────────────┐   │
│   │   HTML5 Templates  │      │     REST API (Gin)      │   │
│   │   (Go server-side) │◄────►│  /api/expenses          │   │
│   │   CSS + JS Client  │      │  /api/budget            │   │
│   └────────────────────┘      │  /api/dashboard         │   │
│                               └───────────┬─────────────┘   │
│                                           │                 │
│                               ┌───────────▼─────────────┐   │
│                               │   SQLite Database (DB)  │   │
│                               │   (Users, Expenses,     │   │
│                               │    Budget Targets)      │   │
│                               └─────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Backend** | Go (Golang) + Gin | High-performance, concurrent web framework & routing |
| **Frontend** | HTML5 + CSS3 + Vanilla JS | Server-side rendered Go templates with responsive Crimson Red theme |
| **Database** | SQLite (`modernc.org/sqlite`) | Embedded pure-Go SQL database with auto migrations and seed data |
| **Security** | `golang.org/x/crypto/bcrypt` | Password hashing & session-based authentication |
| **Container** | Docker (Multi-stage) | Minimal Alpine Linux base image (< 25MB) |
| **Dev Reload**| Air (`github.com/air-verse/air`)| Instant hot-reloading for code, templates, and styles |
| **Registry** | Azure Container Registry | Private Docker container registry (`cloudspend.azurecr.io`) |
| **Hosting** | Azure App Service (Linux) | Managed container hosting with automated HTTPS |

---

## 🚀 Local Development

### Option 1: Live-Reload with Docker Compose & Air (Recommended)
```bash
# Clone the repository
git clone https://github.com/Nishchal-ll/CloudSpend.git
cd CloudSpend

# Start container with live hot-reloading
docker compose up
```
Open **[http://localhost:8080](http://localhost:8080)** in your browser. Any edits to `.go`, `.html`, or `.css` files will instantly rebuild and reload!

### Option 2: Run with Go
```bash
go mod download
go run main.go
```

### Option 3: Standard Docker Build
```bash
docker build -t cloudspend:v1 .
docker run -d -p 8080:8080 --name cloudspend-app cloudspend:v1
```

---

## 📡 REST API Reference

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/healthz` | Container health probe | No |
| `POST`| `/api/auth/login` | Authenticate user & start session | No |
| `POST`| `/api/auth/logout`| Terminate user session | Yes |
| `GET` | `/api/dashboard` | Aggregated spend, budget, & category breakdown | Yes |
| `GET` | `/api/budget` | Get active monthly budget limit | Yes |
| `POST`| `/api/budget` | Update monthly budget limit | Yes |
| `GET` | `/api/expenses` | List all expenses (with search & category filters) | Yes |
| `POST`| `/api/expenses` | Create new expense record | Yes |
| `PUT` | `/api/expenses/:id` | Update existing expense by ID | Yes |
| `DELETE`| `/api/expenses/:id`| Delete expense record by ID | Yes |

---

## ☁️ Azure Deployment Workflow

```bash
# 1. Login to Azure Container Registry
docker login cloudspend.azurecr.io

# 2. Tag production image
docker tag cloudspend:v1 cloudspend.azurecr.io/cloudspend:v1

# 3. Push image to Azure Container Registry
docker push cloudspend.azurecr.io/cloudspend:v1

# 4. Azure App Service automatically pulls and deploys the container on:
# https://cloudspend-app-a9a3bxcce7hye8hv.eastasia-01.azurewebsites.net/
```

---

## 📜 License
MIT License © 2026 [Nishchal Acharya](https://github.com/Nishchal-ll)
