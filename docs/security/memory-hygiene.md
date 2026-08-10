# Sensitive Memory Ownership & Hygiene Rule (Go Backend)

This document defines how `[]byte` slices holding sensitive cryptographic material (master passwords, derived AES keys, decrypted credential payloads) must be handled across Devaulty's Go backend (`backend-go`).

## The Rule

At any point in a execution path, a sensitive byte slice has exactly **one current owner** — the function or struct currently holding a live reference to it.

- Use **mutable `[]byte` slices** instead of Go `string`s for sensitive cryptographic data wherever possible. (Go `string`s are immutable and cannot be zeroed in memory).
- Only the **owner** or **terminal consumer** may zero the slice memory using `clear(slice)` or `for i := range slice { slice[i] = 0 }`.
- In Go functions that process sensitive `[]byte` parameters, use `defer clear(slice)` to guarantee memory wipe upon function completion, even if errors occur.
- When passing a byte slice into another component (e.g., copying a key into `MasterKeySession.SetKey`), the receiver must make its own copy (`copy(dst, src)`) or explicitly assume ownership.

In short:

```text
HTTP Handler (reads JSON, converts to []byte, defer clear)
  → UseCase (receives []byte, defer clear)
    → Key Derivation (computes derived key []byte, defer clear temp buffers)
      → Session Holder (copies bytes into RAM, zeroes old key via clear(oldKey))
```

## Why This Matters

Leaving sensitive passwords or cryptographic keys in memory RAM indefinitely creates vulnerability windows for heap-dumping attacks or cold-boot attacks.

In Go, `clear(slice)` (available in Go 1.21+) overwrites all byte elements with zero values (`0x00`) immediately, avoiding sensitive data persistence in memory before garbage collection cycles.

## Code Patterns in `backend-go`

### 1. HTTP Handlers (`SecurityHandler`)

In `internal/adapter/in/web/handler/security_handler.go`:

```go
func (h *SecurityHandler) SetupMasterPassword(c *gin.Context) {
    var req dto.MasterPassword
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    password := []byte(req.MasterPassword)
    req.MasterPassword = "" // clear reference to string DTO
    defer clear(password)   // guarantee RAM wipe on return

    err := h.vaultUseCase.SetupMasterPassword(c.Request.Context(), password)
    // ...
}
```

### 2. Use Cases (`VaultUseCase`)

In `internal/usecase/vault_usecase.go`:

```go
func (uc *VaultUseCase) SetupMasterPassword(ctx context.Context, password []byte) error {
    defer clear(password) // terminal consumer of input password

    saltBytes, err := uc.keyDeriver.GenerateSalt(SaltLength)
    if err != nil {
        return err
    }
    defer clear(saltBytes) // terminal consumer of salt bytes

    // ...
}
```

### 3. Session Holder (`MasterKeySessionHolderAdapter`)

In `internal/adapter/out/security/master_key_session_holder.go`:

```go
func (m *MasterKeySessionHolderAdapter) clearLocked() {
    clear(m.masterKey) // zeroes underlying RAM bytes
    m.masterKey = nil
    m.lastActivityAt = nil
}
```

## Checklist for Sensitive Data Handling

- [ ] Does this function receive or construct a sensitive `[]byte` slice (password, salt, raw key)?
- [ ] Is `defer clear(slice)` called immediately after slice creation/reception?
- [ ] If storing key bytes in a long-lived struct, does `Clear()` / `SetKey()` overwrite existing bytes with `clear(m.key)` before nil-assigning?
- [ ] Are input DTO string references cleared (`req.Password = ""`) immediately after converting to `[]byte`?
