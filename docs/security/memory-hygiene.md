# Sensitive Memory Ownership & Hygiene Rule (Go Backend)

This document defines how `[]byte` slices holding sensitive cryptographic material (master passwords, derived AES keys, decrypted credential payloads) must be handled across Devaulty's Go backend (`backend`).

## The Rule

At any point in a execution path, a sensitive byte slice has exactly **one current owner** — the function or struct currently holding a live reference to it.

- Use **mutable `[]byte` slices** instead of Go `string`s for sensitive cryptographic data wherever possible (Go `string`s are immutable and cannot be zeroed in memory).
- **Single Ownership Transfer**: Passing a sensitive `[]byte` slice into a downstream method (e.g. passing `password` from `SecurityHandler` to `VaultUseCase`) **transfers ownership** to that component.
- Only the **terminal consumer owner** zeroes the slice memory using `defer clear(slice)`. Callers must **not** perform duplicate zeroing after transferring ownership, unless they explicitly pass a defensive copy.
- When storing a byte slice in a long-lived component (e.g. `MasterKeySession.SetKey`), the receiver must make a defensive copy (`copy(dst, src)`) or explicitly assume ownership.

In short:

```text
HTTP Handler (converts JSON to []byte, clears DTO string, transfers ownership)
  → UseCase (terminal consumer owner: runs `defer clear(password)`)
    → Key Derivation (computes derived key []byte, runs `defer clear(secretKey)`)
      → Session Holder (copies bytes into RAM, zeroes old key via `clear(oldKey)`)
```

## Why This Matters

Leaving sensitive passwords or cryptographic keys in memory RAM indefinitely creates vulnerability windows for heap-dumping attacks or cold-boot attacks.

In Go, `clear(slice)` (available in Go 1.21+) overwrites all byte elements with zero values (`0x00`) immediately, significantly reducing the retention window of sensitive data in RAM (though it zeroes the provided slice specifically rather than immutable strings or intermediate CPU/system buffers).

## Code Patterns in `backend`

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
    req.MasterPassword = "" // clear reference to string DTO (ownership of password slice transfers to VaultUseCase)

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
    clear(m.masterKey) // zeroes provided slice bytes to reduce retention window
    m.masterKey = nil
    m.lastActivityAt = nil
}
```

### 4. Secret Payload DTOs (`dto.SecretBytes`)

In `internal/dto/credential_dto.go`:

```go
type SecretBytes []byte

func (sb *SecretBytes) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*sb = SecretBytes(s)
	return nil
}
```

`SecretBytes` allows JSON unmarshaling of plain strings directly into byte slices, enabling immediate zeroing with `defer clear(cmd.Password)`, `defer clear(cmd.APIKey)`, etc. in `CredentialUseCase`.

## Checklist for Sensitive Data Handling

- [ ] Does this function receive or construct a sensitive `[]byte` / `SecretBytes` slice (password, salt, raw key, secret payload)?
- [ ] Is `defer clear(slice)` called immediately after slice creation/reception?
- [ ] If storing key bytes in a long-lived struct, does `Clear()` / `SetKey()` overwrite existing bytes with `clear(m.key)` before nil-assigning?
- [ ] Are input DTO string references cleared or unmarshaled directly to `SecretBytes` for zeroing in memory?
