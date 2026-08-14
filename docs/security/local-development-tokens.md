# Local Development & Process Security Tokens (Go Backend)

This document explains how the internal process security token (`DEVAULTY_INTERNAL_TOKEN`) operates in the Go backend (`backend`), why it exists, and how to develop and test endpoints locally using Scalar API Docs and cURL.

## The Problem Addressed

Devaulty runs a local Go HTTP server on a local port inside the user's desktop machine. Since local HTTP ports can theoretically be probed by other local applications running under the same user session, all `/api/v1/*` endpoints require an internal process authentication header:

```http
DEVAULTY_INTERNAL_TOKEN: <token>
```

## Production vs Development Behavior

The token validation logic is enforced by `middleware.AuthMiddleware(apiToken)` using `crypto/subtle.ConstantTimeCompare`.

```text
Production Mode (prod)         Development Mode (dev)
======================         ======================
- Random UUID injected by      - Fallback to static token:
  desktop parent process         "dev-token" (APP_ENV=dev)
  at startup                   - Enables Scalar UI at /docs
- Enforces token isolation     - Enables local cURL testing
```

> **Note on `dev-token`**: In development mode (`APP_ENV=dev`), `"dev-token"` is provided as a convenience default for local developer testing. In production mode, the desktop parent process generates a unique random UUID token and injects it via environment variable.

### 1. Security Middleware (`middleware.AuthMiddleware`)
- The middleware receives `apiToken` initialized at server startup in `cmd/api/main.go`.
- In `main.go`:
  ```go
  apiToken := os.Getenv("DEVAULTY_INTERNAL_TOKEN")
  if apiToken == "" && os.Getenv("APP_ENV") == "dev" {
      apiToken = "dev-token"
  }
  ```
- Performs a constant-time comparison on incoming `DEVAULTY_INTERNAL_TOKEN` headers to prevent timing side-channel attacks.
- Requests without a valid token receive `HTTP 401 Unauthorized`.

### 2. Interactive Documentation (Scalar API Reference at `/docs`)
- When running in `APP_ENV=dev`, `/docs` serves the interactive Scalar API documentation UI.
- Interactive documentation dynamically loads `./docs/openapi.yaml`.

## How to Test Endpoints Locally

### Option A: Using Scalar API Docs (`/docs`)

1. Open Scalar API Docs in your browser:
   `http://localhost:8080/docs`
2. Enter `dev-token` in the security authorization header input.
3. All subsequent "Test Request" calls from Scalar UI will automatically send `DEVAULTY_INTERNAL_TOKEN: dev-token`.

### Option B: Using cURL / HTTP Clients

Pass the `DEVAULTY_INTERNAL_TOKEN` header in your request:

```bash
curl -X GET "http://localhost:8080/api/v1/projects" \
     -H "DEVAULTY_INTERNAL_TOKEN: dev-token" \
     -H "Accept: application/json"
```

## Security Best Practices Checklist

- [ ] Production desktop launcher must always generate a random UUID token and inject it as `DEVAULTY_INTERNAL_TOKEN`.
- [ ] `middleware.AuthMiddleware` must use constant-time byte comparison (`subtle.ConstantTimeCompare`) to prevent timing side-channel attacks.
