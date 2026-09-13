# Backend Intern Assignment: Ticket System in Golang

A lightweight, robust RESTful API backend service for a ticket management system built in Go.

## Features & Implementation Highlights

- **Language & Framework**: Golang with `github.com/go-chi/chi/v5` (idiomatic, lightweight router).
- **Authentication**: JWT authentication (`golang-jwt/jwt/v5`) with Bearer token header verification.
- **Security**: Passwords hashed securely using `bcrypt` (`golang.org/x/crypto/bcrypt`).
- **Database & Persistence**: SQLite using pure-Go driver (`modernc.org/sqlite`). Requires zero CGO setup, making Docker compilation fast and lightweight.
- **Ownership Isolation**: Users can **only** list and access tickets created by themselves. Attempting to access another user's ticket returns `404 Not Found`.
- **Status Flow Enforcement**: Strict status validation rule (`open` $\rightarrow$ `in_progress` $\rightarrow$ `closed`). Reopening closed tickets is prohibited and returns `400 Bad Request`.

---

## Submission Details

- **GitHub Repository**: `[YOUR_GITHUB_REPO_URL_HERE]`
- **Deployed Application URL**: `[YOUR_DEPLOYED_APP_URL_HERE]`
- **Public Health Check URL**: `[YOUR_DEPLOYED_APP_URL_HERE]/health`

---

## API Contract & Required Endpoints

| Method | Endpoint | Auth Required | Description | Status Code |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/health` | No | System health check | `200 OK` |
| `POST` | `/auth/register` | No | Register a new user | `201 Created` |
| `POST` | `/auth/login` | No | Login and receive JWT | `200 OK` |
| `POST` | `/tickets` | Yes (`Bearer <token>`) | Create a ticket (initial status `open`) | `201 Created` |
| `GET` | `/tickets` | Yes (`Bearer <token>`) | List logged-in user's tickets | `200 OK` |
| `GET` | `/tickets/{id}` | Yes (`Bearer <token>`) | Get logged-in user's ticket by ID | `200 OK` / `404 Not Found` |
| `PATCH` | `/tickets/{id}/status` | Yes (`Bearer <token>`) | Update own ticket status | `200 OK` / `400 Bad Request` |

### Status Flow Constraints
- Allowed: `open` $\rightarrow$ `in_progress` $\rightarrow$ `closed`
- Allowed: `open` $\rightarrow$ `closed`
- Forbidden: `closed` $\rightarrow$ `open` or `in_progress` (Returns `400 Bad Request`)

---

## Local Run Instructions

### Prerequisites
- Go 1.22+ installed locally.

### 1. Run Directly with Go
```bash
# Clone the repository
git clone <YOUR_GITHUB_REPO_URL>
cd ticket-system

# Run the server
go run cmd/api/main.go
```
The server will start on port `8080`.

### 2. Verify Health Endpoint
```bash
curl http://localhost:8080/health
```
**Expected Response:**
```json
{
  "status": "ok"
}
```

---

## Local Docker Contract

```bash
# Build the Docker image
docker build -t ticket-system .

# Run the container mapping port 8080
docker run -p 8080:8080 ticket-system

# Verify Health Endpoint
curl http://localhost:8080/health
```

---

## Step-by-Step Free Cloud Deployment Guide (Render.com)

You can easily deploy this container to **Render** for free in under 3 minutes (no credit card required).

### Step 1: Push Code to GitHub
Open your terminal in the project directory and run:
```bash
git init
git add .
git commit -m "Initial commit - Ticket System API"
git branch -M main
git remote add origin https://github.com/<YOUR_USERNAME>/<YOUR_REPO_NAME>.git
git push -u origin main
```

### Step 2: Deploy on Render
1. Visit [render.com](https://render.com) and click **Sign Up** (or Log In) using your GitHub account.
2. On your Render Dashboard, click **New +** $\rightarrow$ Select **Web Service**.
3. Connect your GitHub account and select your repository (`ticket-system`).
4. Fill in the service configuration:
   - **Name**: `golang-ticket-system` (or any custom name)
   - **Language / Environment**: Select **Docker** (Render auto-detects the `Dockerfile`).
   - **Region**: Select your closest region (e.g. Oregon, Frankfurt, Singapore).
   - **Instance Type**: Select **Free**.
5. Click **Create Web Service**.

### Step 3: Copy Your Public Deployed URLs
- Render will automatically compile your Docker image and deploy it.
- Once deployed (takes ~1-2 minutes), Render will display your live public URL at the top of the page, for example:
  `https://golang-ticket-system.onrender.com`
- Your public health check URL is:
  `https://golang-ticket-system.onrender.com/health`

Verify public health check:
```bash
curl https://golang-ticket-system.onrender.com/health
```

