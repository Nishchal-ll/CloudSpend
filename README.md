# ☁️ CloudSpend — Cloud Expense & Budget Intelligence Platform

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat&logo=gin)](https://gin-gonic.com/)
[![Database](https://img.shields.io/badge/Database-SQLite-003B57?style=flat&logo=sqlite)](https://www.sqlite.org/)
[![Docker](https://img.shields.io/badge/Container-Docker-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Cloud](https://img.shields.io/badge/Deploy-Azure_App_Service-0078D4?style=flat&logo=microsoftazure)](https://azure.microsoft.com/)

**CloudSpend** is a full-stack cloud cost intelligence and budget management web application built with **Go (Gin)**, **SQLite**, and **Go HTML Templates**. It allows engineering teams to track, categorize, and monitor multi-cloud infrastructure expenditures across Azure, AWS, GCP, Cloudflare, and SaaS providers with real-time budget threshold warnings.

The application is containerized with a lightweight multi-stage Dockerfile and ready for deployment to **Azure Container Registry (ACR)** and **Azure App Service (Linux Container)**.

---

## 🌟 Key Features

* **Multi-Page Web Application:**
  * **Dashboard (`/`):** Real-time KPI cards (Total Budget, Month-to-Date Spend, Remaining Balance, Utilization rate), category spending distribution, cloud provider share bars, and recent transaction log.
  * **Expense Manager (`/expenses`):** Full CRUD interface with search filters, category pills, provider badges, interactive modal forms, and deletion confirmation dialogs.
* **Server-Side Rendered (SSR) Go Templates:** Fast, SEO-optimized HTML rendering powered by Go's native `html/template` and Gin.
* **Full REST API (`/api/*`):** Standard JSON endpoints for dashboard statistics, budget updates, and expense CRUD operations.
* **Embedded SQLite Database:** Pure-Go zero-CGO SQLite implementation with automatic schema migration and starter dataset seeding.
* **Production Docker Container:** Multi-stage build producing an ultra-lightweight Alpine image (< 25MB).
* **Azure Cloud Ready:** Dynamic port binding (`$PORT` / `$WEBSITES_PORT`), health check endpoint (`/healthz`), and single-command deployment compatibility.

---

## 🏗️ Architecture & Project Structure

```text
CloudSpend/
├── .gitignore              # Git ignored files (binaries, databases, .env)
├── .dockerignore           # Docker build exclusions
├── Dockerfile              # Multi-stage production build (Golang builder -> Alpine)
├── README.md               # Project documentation & Azure deployment guide
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── main.go                 # Application entrypoint & HTTP router configuration
│
├── database/
│   └── db.go               # SQLite initialization, schema migration & seed data
│
├── models/
│   ├── expense.go          # Expense model & database CRUD queries
│   └── budget.go           # Budget model, statistics aggregation & calculations
│
├── handlers/
│   ├── expense_handler.go  # REST API handlers for /api/expenses
│   ├── dashboard_handler.go# REST API handlers for /api/dashboard & /api/budget
│   └── web_handler.go      # Go HTML template renderers for / and /expenses
│
├── templates/
│   ├── index.html          # Dashboard page template
│   └── expenses.html       # Expenses management page template
│
└── static/
    ├── css/
    │   └── style.css       # Modern glassmorphism dark theme styling
    └── js/
        └── app.js          # Interactive modals, AJAX requests, toasts & validation
```

---

## 🚀 Quickstart: Local Development

### Prerequisites
* **Go 1.22+** installed (`go version`)
* **Docker** installed (`docker --version`)

### 1. Run with Go
```bash
# Clone repository
git clone https://github.com/Nishchal-ll/CloudSpend.git
cd CloudSpend

# Download dependencies
go mod download

# Run application
go run main.go
```

Open your browser at: **[http://localhost:8080](http://localhost:8080)**

---

## 📡 REST API Reference

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/healthz` | Container health probe |
| `GET` | `/api/dashboard` | Aggregated spend, budget, and category metrics |
| `GET` | `/api/budget` | Get active monthly budget configuration |
| `POST` | `/api/budget` | Update monthly budget limit |
| `GET` | `/api/expenses` | List expenses (optional: `?category=Compute&search=Azure`) |
| `GET` | `/api/expenses/:id` | Get single expense by ID |
| `POST` | `/api/expenses` | Create new expense |
| `PUT` | `/api/expenses/:id` | Update expense by ID |
| `DELETE` | `/api/expenses/:id` | Delete expense by ID |

### Example: Create Expense Payload
```http
POST /api/expenses HTTP/1.1
Content-Type: application/json

{
  "title": "Azure App Service Linux B1",
  "provider": "Azure",
  "category": "Compute",
  "amount": 54.75,
  "expense_date": "2026-09-17",
  "description": "Production container hosting"
}
```

---

## 🐳 Docker Containerization

### 1. Build Docker Image
```bash
docker build -t cloudspend:v1 .
```

### 2. Run Container Locally
```bash
docker run -d -p 8080:8080 --name cloudspend-app cloudspend:v1
```

Visit: **[http://localhost:8080](http://localhost:8080)**

### 3. Verify Running Container
```bash
docker ps
docker logs cloudspend-app
```

### 4. Stop Container
```bash
docker stop cloudspend-app && docker rm cloudspend-app
```

---

## ☁️ Azure Deployment Step-by-Step

Follow these steps to deploy **CloudSpend** to **Azure Container Registry (ACR)** and **Azure App Service**:

### Step 1: Login to Azure CLI
```bash
az login
```

### Step 2: Create a Resource Group
```bash
az group create --name CloudSpend-RG --location eastus
```

### Step 3: Create Azure Container Registry (ACR)
```bash
# Note: ACR name must be globally unique and alphanumeric
az acr create --resource-group CloudSpend-RG --name cloudspendacr2026 --sku Basic --admin-enabled true
```

### Step 4: Build, Tag, and Push Image to ACR
```bash
# Log in to your ACR
az acr login --name cloudspendacr2026

# Tag local image for ACR
docker tag cloudspend:v1 cloudspendacr2026.azurecr.io/cloudspend:v1

# Push image to ACR
docker push cloudspendacr2026.azurecr.io/cloudspend:v1
```

### Step 5: Create App Service Plan (Linux)
```bash
az appservice plan create --name CloudSpend-Plan --resource-group CloudSpend-RG --sku B1 --is-linux
```

### Step 6: Deploy Web App using Docker Container
```bash
# Retrieve ACR admin password
ACR_PASSWORD=$(az acr credential show --name cloudspendacr2026 --query "passwords[0].value" -o tsv)

# Create Web App
az webapp create \
  --resource-group CloudSpend-RG \
  --plan CloudSpend-Plan \
  --name cloudspend-webapp \
  --deployment-container-image-name cloudspendacr2026.azurecr.io/cloudspend:v1 \
  --docker-registry-server-user cloudspendacr2026 \
  --docker-registry-server-password "$ACR_PASSWORD"
```

### Step 7: Configure Port Setting
```bash
az webapp config appsettings set --resource-group CloudSpend-RG --name cloudspend-webapp --settings WEBSITES_PORT=8080 PORT=8080 GIN_MODE=release
```

### Step 8: Access Your Live Application
Your public URL will be:
```text
https://cloudspend-webapp.azurewebsites.net
```

---

## 📋 University Submission Checklist

- [x] **Source Code:** Clean modular Go structure (`main.go`, `handlers/`, `models/`, `database/`, `templates/`, `static/`).
- [x] **Go HTML Templates:** 2+ full web pages (Dashboard `/` and Expenses Manager `/expenses`).
- [x] **REST API:** Complete CRUD API with structured JSON responses.
- [x] **Dockerfile:** Multi-stage production container build.
- [x] **Docker Image:** Tested locally on `http://localhost:8080`.
- [x] **Azure ACR:** Image pushed to Azure Container Registry.
- [x] **Azure App Service:** Running Linux container deployment with public HTTPS URL.
- [x] **Git Repository:** GitHub connected at [https://github.com/Nishchal-ll/CloudSpend](https://github.com/Nishchal-ll/CloudSpend).

---

## 📜 License
MIT License © 2026 [Nishchal-ll](https://github.com/Nishchal-ll)
