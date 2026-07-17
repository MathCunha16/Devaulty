package com.devaulty.backend.adapter.in.web.credential.dto;

import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import jakarta.validation.constraints.AssertTrue;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record UpdateCredentialRequest(
        @Pattern(regexp = "(?s).*\\S.*", message = "Title must not be blank")
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,

        CredentialSecretType secretType,

        char[] username,
        char[] password,
        char[] apiKey,
        char[] rawTextContent,

        String notes,
        String relatedUrl
){
    @AssertTrue(message = "Required secret fields are missing for the selected secret type")
    public boolean isValidSecretPayload() {
        boolean isAnySecretFieldPresent = secretType != null
                || (username != null && username.length > 0)
                || (password != null && password.length > 0)
                || (apiKey != null && apiKey.length > 0)
                || (rawTextContent != null && rawTextContent.length > 0);

        if (!isAnySecretFieldPresent) {
            return true;
        }

        if (secretType == null) {
            return false;
        }

        return switch (secretType) {
            case LOGIN -> password != null && password.length > 0;
            case API_KEY -> apiKey != null && apiKey.length > 0;
            case RAW_TEXT -> rawTextContent != null && rawTextContent.length > 0;
        };
    }
}
