 Backend Intern Assignment: Ticket System in Golang

A Restful API backend service for a ticket management system built in Go.

 Features & Implementation Highlights

-  Language & Framework: Golang with `github.com/go-chi/chi/v5` (idiomatic, lightweight router).
-  Authentication: JWT authentication (`golang-jwt/jwt/v5`) with Bearer token header verification.
-  Security: Passwords hashed securely using `bcrypt` (`golang.org/x/crypto/bcrypt`).
-  Database & Persistence: SQLite using pure-Go driver (`modernc.org/sqlite`).
-  Ownership Isolation: Users can **only** list and access tickets created by themselves. Attempting to access another user's ticket returns `404 Not Found`.
-  Status Flow Enforcement: Strict status validation rule (`open` $\rightarrow$ `in_progress` $\rightarrow$ `closed`). Reopening closed tickets is prohibited and returns `400 Bad Request`.

---

 Submission Details

- Deployed Application URL: ` `
- Public Health Check URL: ` '

---

API Contract & Required Endpoints

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


