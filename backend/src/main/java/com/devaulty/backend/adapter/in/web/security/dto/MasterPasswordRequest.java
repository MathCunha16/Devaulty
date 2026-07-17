package com.devaulty.backend.adapter.in.web.security.dto;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

public record MasterPasswordRequest(
        @NotNull(message = "Password cannot be null")
        @Size(min = 8, max = 255, message = "Password must be between 8 and 255 characters")
        char[] masterPassword
) {
}
