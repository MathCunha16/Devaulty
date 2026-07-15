package com.devaulty.backend.adapter.in.web.credential.dto;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.time.LocalDateTime;
import java.util.UUID;

public record CredentialSummaryResponse(
        UUID id,
        UUID projectId,
        String title,
        CredentialSecretType secretType,
        String relatedUrl,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
