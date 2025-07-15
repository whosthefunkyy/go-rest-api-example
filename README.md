Simple REST API in Go (Golang) — CRUD for users

This project is a training REST API server written in Go. Basic CRUD for user management is implemented, with PostgreSQL connection and the ability to run via Docker.

REST API with endpoints:

GET /users — get a list of all users
GET /users/{id} — get a user by ID
POST /users — create a new user
PUT /users/{id} — update a user
DELETE /users/{id} — delete a user

- Connect to PostgreSQL

- Docker + Docker Compose

-.env configuration

- Tests (in development)

Language: Golang
Router:  gorilla/mux 
Database: PostgreSQL
Working with DB: GORM
Containerization: Docker, Docker Compose
Testing: testing, httptest (in development)