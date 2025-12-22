## Go REST API with CI/CD (GitHub Actions + AWS Elastic Beanstalk) by Artem Melnychuk

This project is a training REST API server written in Go (Golang) that implements basic CRUD operations for user management.

### Features

- REST API with the following endpoints:
  - `GET /users` — retrieve all users
  - `GET /users/{id}` — retrieve a user by ID
  - `POST /users` — create a new user
  - `PUT /users/{id}` — update an existing user
  - `DELETE /users/{id}` — delete a user

- PostgreSQL integration using GORM
- Environment-based configuration via `.env`
- Containerized with Docker and Docker Compose
- Tests (in progress) using `testing` and `httptest`

### CI/CD

A CI/CD pipeline is implemented using **GitHub Actions** to automatically build and deploy the application to **AWS Elastic Beanstalk**.

On every push to the `main` branch, the pipeline performs the following steps:
- Checks out the source code
- Builds a Linux binary of the Go application
- Packages the binary and configuration files into a deployment artifact
- Uploads the artifact to Amazon S3
- Creates a new Elastic Beanstalk application version
- Deploys the new version to the target environment

### Tech Stack

- Language: Go (Golang)
- Router: gorilla/mux
- Database: PostgreSQL
- ORM: GORM
- Containerization: Docker, Docker Compose
- CI/CD: GitHub Actions
- Cloud: AWS (Elastic Beanstalk, S3)
