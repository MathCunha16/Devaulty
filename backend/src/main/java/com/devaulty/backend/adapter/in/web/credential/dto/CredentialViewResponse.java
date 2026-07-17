package com.devaulty.backend.adapter.in.web.credential.dto;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.time.LocalDateTime;
import java.util.Map;
import java.util.UUID;

public record CredentialViewResponse(
        UUID id,
        UUID projectId,
        String title,
        CredentialSecretType secretType,
        Map<String, String> decryptedPayload,
        String notes,
        String relatedUrl,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
