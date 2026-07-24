# Local Development & Process Security Tokens

This document explains how the internal process security token (`X-Devaulty-Internal-Token`) operates, why it exists, and how to develop and test endpoints locally using Swagger UI and cURL.

## The Problem Addressed

Devaulty runs a local Spring Boot server on an ephemeral port inside the user's desktop machine. Since local HTTP ports can theoretically be probed by other local applications running under the same user session, all `/api/v1/**` endpoints require an internal process authentication header:

```http
X-Devaulty-Internal-Token: <token>
```

## Production vs Development Behavior

The token validation logic is enforced by `InternalAppTokenFilter`.

```
Production Mode (prod)         Development Mode (dev)
======================         ======================
- Random UUID generated        - Accepts static token from application-dev.yml
  in-memory at process boot      ("dev-secret-token")
- Injected via JavaFX into     - Accepts random UUID as fallback
  WebView before React mounts  - Enables Swagger UI & cURL testing
```

### 1. Production Mode (`application-prod.yml`)
- `AppTokenContext.PROCESS_TOKEN` generates a cryptographically unique UUID (`UUID.randomUUID()`) in-memory when the JVM boots.
- `DevaultyDesktop.java` injects this token directly into the JavaFX `WebView` JS context (`window.DEVAULTY_INTERNAL_TOKEN`) during the `Worker.State.RUNNING` state before React components mount.
- External requests without this exact in-memory UUID receive `HTTP 403 Forbidden`.

### 2. Local Development Mode (`application-dev.yml`)
- `src/main/resources/application-dev.yml` declares a development token:

```yaml
devaulty:
  dev:
    token: "dev-secret-token"
```

- `InternalAppTokenFilter` reads `@Value("${devaulty.dev.token:#{null}}")`. If present, it permits requests carrying `"dev-secret-token"`.
- In `application-prod.yml`, `devaulty.dev.token` is completely absent. `devToken` evaluates to `null`, ensuring development tokens can **never** be used in production builds.

## How to Test Endpoints Locally

### Option A: Using Swagger UI

1. Open Swagger UI in your browser:
   `http://localhost:8080/swagger-ui.html`
2. Click the green **Authorize 🔓** button in the top right corner.
3. In the `DevaultyInternalToken` value input, enter:
   `dev-secret-token`
4. Click **Authorize** and close the modal.
5. All subsequent "Try it out" API calls from Swagger UI will automatically include the `X-Devaulty-Internal-Token: dev-secret-token` header.

### Option B: Using cURL / HTTP Clients

Pass the `X-Devaulty-Internal-Token` header in your request:

```bash
curl -X GET "http://localhost:8080/api/v1/release/check" \
     -H "X-Devaulty-Internal-Token: dev-secret-token" \
     -H "Accept: application/json"
```

## Security Best Practices Checklist

- [ ] `application-prod.yml` must **never** contain a `devaulty.dev.token` property.
- [ ] `InternalAppTokenFilter` must validate `devToken` only when non-null and non-blank (`devToken != null && !devToken.isBlank()`).
- [ ] Production builds (`jpackage`) must always execute with `-Dspring.profiles.active=prod`.
