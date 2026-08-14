<p align="center">
  <img src="../docs/assets/icon.png" alt="Devaulty Icon" width="140" />
</p>

<h1 align="center">Devaulty Backend (Go)</h1>

<p align="center">
  <strong>High-performance, zero-allocation memory hygiene, local-first REST API built with Go and Hexagonal Architecture.</strong>
</p>

---

## 🪦 In Memoriam: The Java Era

Devaulty's initial backend prototype was built with Java and Spring Boot. The backend was completely rewritten in Go, reducing memory usage by over 80%, booting in under 10 milliseconds, and producing a single self-contained native binary.

---

## 🏗️ Architecture & Organization

The backend is structured using **Hexagonal Architecture (Ports and Adapters)**, enforcing strict boundary separation between core domain rules, application use cases, and infrastructure adapters:

```text
backend/
├── cmd/api/            # Application entrypoint (main.go)
├── internal/
│   ├── domain/
│   │   ├── model/      # Pure domain entities (Project, Credential, Snippet, etc.)
│   │   └── port/       # Repository & Service interfaces (Inbound / Outbound ports)
│   ├── usecase/        # Application business logic & workflow orchestrators
│   ├── adapter/
│   │   ├── in/web/     # HTTP Handlers (Gin), Routers & Middleware
│   │   ├── in/scheduler/# Background cron/autolock tickers
│   │   └── out/        # Persistence (SQLite/sqlx) & Security (Argon2id, AES-GCM)
│   └── dto/            # Data Transfer Objects & JSON serialization rules
├── migrations/         # SQL migration scripts (embedded via go:embed)
└── docs/               # OpenAPI 3.0 specification (openapi.yaml)
```

---

## 🔐 Security & Cryptography Engine

Security is the core priority of the Devaulty backend. No unencrypted sensitive data is ever written to disk.

### 1. Key Derivation (Argon2id)
- Master passwords are passed directly into the Argon2id key deriver.
- Uses cryptographically secure random salts and memory-hard parameters ($m=64\text{MB}, t=3, p=4$) to mitigate brute-force and GPU hardware attacks.

### 2. Payload Encryption (AES-256-GCM)
- All credential payloads (passwords, tokens, secret keys) are encrypted using **AES-256-GCM** (Galois/Counter Mode).
- Provides both confidentiality and authenticated integrity checks (AEAD), preventing cipher tampering.

### 3. IPC Process Isolation
- To prevent unauthorized local user processes from calling the API, all endpoints require the `DEVAULTY_INTERNAL_TOKEN` header.
- Validated via constant-time comparison (`crypto/subtle.ConstantTimeCompare`) to eliminate timing side-channel attacks.
- Read more: **[Local Security Tokens Documentation](../docs/security/local-development-tokens.md)**

---

## 🧹 Memory Hygiene & Zeroing

Unlike standard garbage-collected strings, immutable memory references can linger in RAM heap dumps. Devaulty implements strict **Memory Ownership & Hygiene**:

- **Mutable Byte Slices**: Sensitive payloads use `[]byte` and `dto.SecretBytes` instead of immutable Go `string`s.
- **Single Ownership Transfer**: Functions explicitly transfer slice ownership down the execution stack.
- **Prompt Zeroing**: Terminal consumers call `defer clear(slice)` immediately upon processing to zero memory bytes in RAM.
- Read more: **[Memory Hygiene Rules Documentation](../docs/security/memory-hygiene.md)**

---

## 💾 Database & Embedded Migrations

- **Database**: SQLite3 managed via `jmoiron/sqlx` for fast, lightweight local persistence.
- **Embedded Migrations**: All `.sql` migration files in `migrations/` are compiled directly into the binary using Go 1.16+ `//go:embed`. No external migration directories need to be distributed with the application.

---

## 📖 Interactive API Documentation (Scalar)

When running in development mode (`APP_ENV=dev`), Devaulty serves interactive API documentation at:

```http
http://localhost:8080/docs
```

Powered by **Scalar API Reference**, it dynamically parses `docs/openapi.yaml` to provide interactive endpoint testing directly from your browser using the default development token (`dev-token`).

---

## 📦 Key Dependencies

- **[Gin Gonic](https://github.com/gin-gonic/gin)** — High-performance HTTP web framework.
- **[sqlx](https://github.com/jmoiron/sqlx)** + **[go-sqlite3](https://github.com/mattn/go-sqlite3)** — SQLite database driver and extensions.
- **[golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto)** — Argon2id key derivation & cryptographic primitives.
- **[go-scalar-api-reference](https://github.com/MarceloPetrucio/go-scalar-api-reference)** — Embedded OpenAPI documentation renderer.

---

## 🛠️ Local Development & Testing

### Running in Development Mode
```bash
APP_ENV=dev go run ./cmd/api/
```
*App runs on `http://127.0.0.1:8080` with API docs available at `/docs`.*

### Running Tests
```bash
go test ./... -v
```
