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

- **GitHub Repository**: `https://github.com/Ashutosh437/Ticket-Booking`
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

 Status Flow Constraints
- Allowed: `open` $\rightarrow$ `in_progress` $\rightarrow$ `closed`
- Allowed: `open` $\rightarrow$ `closed`
- Forbidden: `closed` $\rightarrow$ `open` or `in_progress` (Returns `400 Bad Request`)


