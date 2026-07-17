# Sensitive Memory Ownership Rule

This document defines how `byte[]` / `char[]` arrays holding sensitive
cryptographic material (master passwords, decrypted credential payloads) must
be handled across the Devaulty backend.

## The rule

At any point in a call chain, a sensitive array has exactly **one current
owner** — the last piece of code holding a live reference to it.

- Only the **owner** may zero the array (`Arrays.fill(array, (byte) 0)` for
  `byte[]`, `Arrays.fill(array, '\0')` for `char[]`).
- Passing the array into another object (a record field, a return value,
  a collaborator method) **transfers ownership**. The original holder must
  **not** zero it afterward — doing so corrupts the data for whoever
  received it.
- Ownership transfers **once**. The new owner either passes it along again
  (transferring ownership further) or is the **terminal consumer**, in
  which case it must zero the array once it's done extracting what it needs.

In short:

```
Creator → [passes along, does NOT zero] → ... → Terminal consumer → ZEROES
```

## Why this matters

Zeroing an array that another object still references doesn't throw an
exception — it silently overwrites the data everyone else sees. The bug is
invisible until something downstream tries to read the (now-blank) value,
often producing a confusing, unrelated-looking error (e.g. a JSON parser
choking on null bytes) far from the line that actually caused it.

## Real example from this codebase

**Bug (July 2026):** `CredentialWebMapper#toCreateCredentialCommand` built a
`char[] serializedPayload`, passed it into `CreateCredentialCommand`, and
then zeroed `serializedPayload` in its own `finally` block "for safety."
Since `CreateCredentialCommand.payload()` held the *same* array reference,
the command reached `CreateCredentialImpl` already blanked out. The
resulting ciphertext encrypted null bytes, and decryption later failed to
parse as JSON — the error surfaced in a completely different class
(`CredentialWebMapper#jsonToMap`) than the one that caused it.

**Fix:** the creator of `serializedPayload` stopped zeroing it. Ownership
now flows cleanly: `toCreateCredentialCommand` → `CreateCredentialImpl` (uses
it to encrypt, then zeroes it, since nothing downstream needs the plaintext
payload anymore) → the *decrypted* result flows separately to
`DecryptedCredential` → `CredentialWebMapper#jsonToMap` (the terminal
consumer, which zeroes it after parsing).

## Checklist when writing code that touches sensitive arrays

- [ ] Does this method receive a sensitive array as a parameter?
  - If yes: am I the terminal consumer, or am I passing it further
    (directly, or wrapped inside another object/record)?
- [ ] If I'm passing it further: do **not** zero it in my `finally` block.
  Say so explicitly in a comment, so the next person doesn't "helpfully"
  add a zero-fill.
- [ ] If I'm the terminal consumer: zero it in a `finally` block,
  unconditionally, right after I've extracted what I need.
- [ ] Document the ownership contract in the method's Javadoc if the
  method is a public/protected entry point others will call — see
  `CreateCredentialImpl#execute` for the expected format.

## Where this currently applies

- `SetupMasterPasswordImpl` / `UnlockVaultImpl` — the `char[] password` input
- `CreateCredentialImpl` — `command.payload()` (input), `decryptedBytes`
  (output, ownership transferred via `DecryptedCredential`)
- `CredentialWebMapper` — `serializedPayload` (transferred to the command),
  `decryptedBytes` in `jsonToMap` (terminal — zeroed after parsing)
- Any future use case that decrypts a credential (e.g. `GetCredentialImpl`,
  `UpdateCredentialImpl`) must follow the same pattern.
