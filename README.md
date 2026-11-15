# MiniAuthSystem

A microservices-based authentication system built with **Go, gRPC, Protocol Buffers, PostgreSQL, and JWT**.

MiniAuthSystem demonstrates how independent Go services communicate through gRPC, how authentication workflows are separated into reusable components, and how PostgreSQL is integrated for persistent user storage.

## Architecture

The application consists of four Go modules:

- **Gateway Service:** Receives HTTP REST requests and forwards them to the Auth Service through gRPC.
- **Auth Service:** Handles OTP generation, OTP verification, JWT authentication, and user operations.
- **Protobufs:** Defines the communication contract between services.
- **Service Utils:** Contains shared utilities for configuration, logging, and gRPC communication.

### Architecture Diagram

```mermaid
flowchart TD
    A["Client / Postman"] -->|"HTTP REST"| B["Gateway Service :8080"]

    B --> C["Routes"]
    C --> D["HTTP Handlers"]
    D --> E["gRPC Client"]

    E -->|"gRPC + Protobuf"| F["Auth Service :9090"]

    F --> G["gRPC Server"]
    G --> H["Auth Handler"]

    H --> I["OTP Service"]
    H --> J["JWT Service"]
    H --> K["User Repository"]

    K -->|"SQL / pgx"| L[("PostgreSQL")]

    M["Protobufs"] -.-> E
    M -.-> G
    N["Service Utils"] -.-> B
    N -.-> F
```

### Request Flow

```text
Client
  |
  | HTTP Request
  v
Gateway Service
  |
  | gRPC Request
  v
Auth Service
  |
  | Business Logic
  v
OTP / JWT Services
  |
  | Database Operations
  v
PostgreSQL
```

## Project Structure

```text
MiniAuhSystem/
│
├── gateway-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── clients/
│   │   │   └── auth_client.go
│   │   ├── constants/
│   │   ├── handlers/
│   │   │   └── auth/
│   │   │       └── auth_handler.go
│   │   └── routes/
│   │       └── routes.go
│   ├── .env
│   ├── .air.toml
│   ├── Makefile
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── auth-service/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── client/
│   │   ├── constants/
│   │   ├── handler/
│   │   │   └── auth_handler.go
│   │   ├── model/
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── postgres.go
│   │   │   └── user_repository.go
│   │   ├── server/
│   │   │   └── grpc.go
│   │   └── service/
│   │       ├── otp_service.go
│   │       ├── jwt_service.go
│   │       ├── otp_service_test.go
│   │       └── jwt_service_test.go
│   ├── .env
│   ├── .air.toml
│   ├── Makefile
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
│
├── protobufs/
│   ├── proto/
│   │   └── auth/
│   │       └── v1/
│   │           └── auth.proto
│   ├── generated/
│   │   └── auth/
│   │       └── v1/
│   │           ├── auth.pb.go
│   │           └── auth_grpc.pb.go
│   ├── Makefile
│   └── go.mod
│
├── service-utils/
│   ├── client/
│   │   └── grpc.go
│   ├── config/
│   │   └── config.go
│   ├── grpcerror/
│   │   └── errors.go
│   ├── logger/
│   │   └── logger.go
│   ├── metrics/
│   ├── server/
│   │   └── grpc.go
│   ├── writer/
│   │   └── response.go
│   └── go.mod
│
└── go.work
```

## Technologies Used

| Technology | Purpose |
|---|---|
| Go | Backend development |
| PostgreSQL | Persistent user storage |
| gRPC | Communication between microservices |
| Protocol Buffers | Service contracts and message serialization |
| JWT | Token-based authentication |
| pgx | PostgreSQL driver and connection pooling |
| Air | Automatic application rebuilding during development |
| Makefile | Development commands and automation |
| grpcurl | Testing gRPC services |
| Postman | Testing REST API endpoints |
| Git | Version control |

## Prerequisites

Install the following before running the application:

- Go 1.26.8 or compatible version
- PostgreSQL
- Protocol Buffers compiler (`protoc`)
- `protoc-gen-go`
- `protoc-gen-go-grpc`
- Air
- Make
- Git
- Postman (optional, for API testing)

## Setup Instructions

### 1. Clone the repository

```bash
git clone <YOUR_GITHUB_REPOSITORY_URL>
cd MiniAuhSystem
```

### 2. Configure Go workspace

The root `go.work` file connects all four local Go modules.

Make sure the workspace contains:

```go
go 1.26.8

use (
    ./auth-service
    ./gateway-service
    ./protobufs
    ./service-utils
)
```

### 3. Configure PostgreSQL

Create a database:

```sql
CREATE DATABASE mini_auth;
```

Create an application user:

```sql
CREATE USER mini_auth_app WITH PASSWORD 'YOUR_LOCAL_PASSWORD';

GRANT CONNECT ON DATABASE mini_auth TO mini_auth_app;
```

Connect to the `mini_auth` database and create the users table:

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Grant the application user the necessary schema permissions:

```sql
GRANT USAGE, CREATE ON SCHEMA public TO mini_auth_app;
```

### 4. Configure environment variables

Create `auth-service/.env`:

```env
DATABASE_URL=postgres://mini_auth_app:YOUR_LOCAL_PASSWORD@localhost:5432/mini_auth?sslmode=disable

JWT_SECRET=replace-with-a-long-random-development-secret

GRPC_PORT=9090
```

Create `gateway-service/.env`:

```env
HTTP_PORT=8080

AUTH_SERVICE_ADDRESS=localhost:9090
```

Replace the database address and credentials with values appropriate for your environment.

**Security:** Do not commit `.env` files or real credentials to GitHub. Add them to `.gitignore`.

### 5. Download dependencies

From the project root:

```bash
cd auth-service
go mod download

cd ../gateway-service
go mod download

cd ../protobufs
go mod download

cd ../service-utils
go mod download
```

### 6. Run the Auth Service

Open a terminal:

```bash
cd auth-service
make run
```

The Auth Service starts its gRPC server on port `9090`.

### 7. Run the Gateway Service

Open another terminal:

```bash
cd gateway-service
make run
```

The Gateway Service starts its HTTP server on port `8080`.

The application endpoints are now available at:

```text
http://localhost:8080
```

## API Endpoints

Base URL:

```text
http://localhost:8080
```

### 1. Login

Generates an OTP for the supplied phone number.

**Endpoint**

```http
POST /api/v1/auth/login
```

**Request Body**

```json
{
  "phone_number": "9876543210"
}
```

**Example Response**

```json
{
  "message": "OTP generated successfully"
}
```

For local development, the generated OTP is printed in the Auth Service terminal.

### 2. Verify OTP

Verifies the OTP and generates a JWT access token.

**Endpoint**

```http
POST /api/v1/auth/verify-otp
```

**Request Body**

```json
{
  "phone_number": "9876543210",
  "otp": "1234"
}
```

**Example Response**

```json
{
  "access_token": "YOUR_JWT_ACCESS_TOKEN"
}
```

### 3. Get User Profile

Retrieves the authenticated user's profile.

**Endpoint**

```http
GET /api/v1/auth/profile
```

**Headers**

```http
Authorization: Bearer YOUR_JWT_ACCESS_TOKEN
```

**Example Response**

```json
{
  "user_id": "USER_UUID",
  "phone_number": "9876543210"
}
```

### API Summary

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/auth/login` | Generate OTP |
| POST | `/api/v1/auth/verify-otp` | Verify OTP and generate JWT |
| GET | `/api/v1/auth/profile` | Retrieve authenticated user profile |

## Authentication Flow

### Step 1: Login

The client sends a phone number to the Gateway Service.

The Gateway forwards the request to the Auth Service using gRPC.

The Auth Service generates a four-digit OTP and temporarily stores it with an expiration time.

### Step 2: OTP Verification

The client submits the phone number and OTP.

The Auth Service validates the OTP.

After successful verification:

- The OTP is invalidated.
- The user is retrieved from PostgreSQL.
- A new user record is created if the user does not already exist.
- A JWT access token is generated.

### Step 3: Access Protected Profile

The client sends the JWT in the Authorization header.

The Auth Service validates the token and retrieves the associated user information from PostgreSQL.

The profile is returned through the Gateway.

### Authentication Sequence Diagram

```mermaid
sequenceDiagram
    participant C as Client
    participant G as Gateway
    participant A as Auth Service
    participant O as OTP Service
    participant J as JWT Service
    participant R as Repository
    participant DB as PostgreSQL

    C->>G: POST /auth/login
    G->>A: Login RPC
    A->>O: Generate OTP
    O-->>A: OTP generated
    A-->>G: Login response
    G-->>C: Success response

    C->>G: POST /auth/verify-otp
    G->>A: VerifyOTP RPC
    A->>O: Verify OTP
    O-->>A: Verification success
    A->>R: Get user by phone
    R->>DB: SELECT user
    DB-->>R: User record or no rows
    R-->>A: User information
    A->>R: Create user if needed
    A->>J: Generate JWT
    J-->>A: Access token
    A-->>G: VerifyOTP response
    G-->>C: JWT access token

    C->>G: GET /auth/profile + JWT
    G->>A: GetProfile RPC
    A->>J: Validate JWT
    J-->>A: Valid claims
    A->>R: Get user by phone
    R->>DB: SELECT user
    DB-->>R: User record
    R-->>A: Profile information
    A-->>G: Profile response
    G-->>C: HTTP 200 + user profile
```

## Development Commands

Run these commands from the relevant service directory.

### Auth Service

| Command | Description |
|---|---|
| `make run` | Start with Air hot reloading |
| `make build` | Build the application |
| `make test` | Run tests |
| `make clean` | Remove build artifacts |
| `make fmt` | Format Go source files |
| `make lint` | Run golangci-lint |

### Gateway Service

| Command | Description |
|---|---|
| `make run` | Start with Air hot reloading |
| `make build` | Build the application |
| `make test` | Run tests |
| `make clean` | Remove build artifacts |
| `make fmt` | Format Go source files |
| `make lint` | Run golangci-lint |

### Inspect gRPC Services

With Auth Service running:

```bash
grpcurl -plaintext localhost:9090 list
```

Inspect the authentication service:

```bash
grpcurl -plaintext localhost:9090 list auth.v1.AuthService
```

## Testing

The project includes:

- OTP generation and verification unit tests.
- JWT generation and validation unit tests.
- User repository tests.
- Authentication handler tests.
- Integration tests involving PostgreSQL.

Run the tests:

```bash
cd auth-service
make test
```

```bash
cd gateway-service
make test
```

## Current Development Scope

- OTP delivery uses terminal output for local development; no SMS provider is integrated.
- OTP data is stored in application memory.
- gRPC uses insecure transport credentials for local development.
- The project has been tested locally with PostgreSQL, gRPC, and Postman.
- Production deployment, TLS, external OTP delivery, and production-grade security hardening are not yet implemented.

## Learning Objectives

This project demonstrates practical implementation of:

- Go backend development.
- Microservices architecture.
- HTTP REST APIs.
- gRPC communication.
- Protocol Buffers.
- PostgreSQL integration.
- Repository and service-layer separation.
- OTP-based authentication.
- JWT token generation and validation.
- Unit and integration testing.
- Go workspace management.

## Author

**Roshan**

GitHub: [Add your GitHub profile]

LinkedIn: [Add your LinkedIn profile]
