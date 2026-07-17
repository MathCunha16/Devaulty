package com.devaulty.backend.application.port.in.credential;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.util.UUID;

public record CreateCredentialCommand(
        UUID projectId,
        String title,
        CredentialSecretType secretType,
        char[] payload,
        String notes,
        String relatedUrl
) {
}
