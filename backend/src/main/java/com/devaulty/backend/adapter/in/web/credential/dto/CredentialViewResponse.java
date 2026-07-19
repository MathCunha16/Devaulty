package com.devaulty.backend.adapter.in.web.credential.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import com.devaulty.backend.adapter.in.web.tag.dto.TagSummaryResponse;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.time.LocalDateTime;
import java.util.List;
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
        @Schema(description = "List of tags associated with this credential")
        List<TagSummaryResponse> tags,
        LocalDateTime createdAt,
        LocalDateTime updatedAt
) {
}
