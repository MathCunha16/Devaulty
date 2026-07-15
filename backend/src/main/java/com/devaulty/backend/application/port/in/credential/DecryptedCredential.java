package com.devaulty.backend.application.port.in.credential;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.time.LocalDateTime;
import java.util.UUID;

public record DecryptedCredential(
        UUID id,
        UUID projectId,
        String title,
        CredentialSecretType secretType,
        byte[] decryptedPayload,
        String notes,
        String relatedUrl,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
