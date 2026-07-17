package com.devaulty.backend.adapter.in.web.credential.dto;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

public record CreateCredentialRequest(

        @NotBlank(message = "Title is required")
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,

        @NotNull(message = "Secret type is required")
        CredentialSecretType secretType,

        char[] username,
        char[] password,
        char[] apiKey,
        char[] rawTextContent,

        String notes,
        String relatedUrl

) {
    @AssertTrue(message = "Required secret fields are missing for the selected secret type")
    public boolean isValidSecretPayload() {
        if (secretType == null) return false;

        return switch (secretType) {
            case LOGIN -> password != null && password.length > 0;
            case API_KEY -> apiKey != null && apiKey.length > 0;
            case RAW_TEXT -> rawTextContent != null && rawTextContent.length > 0;
        };
    }
}
