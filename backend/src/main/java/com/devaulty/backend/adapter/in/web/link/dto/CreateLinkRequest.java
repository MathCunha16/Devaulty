package com.devaulty.backend.adapter.in.web.link.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record CreateLinkRequest(

        @NotBlank(message = "Title must not be blank")
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,

        @NotBlank(message = "URL must not be blank")
        String url,

        String description
) {
}
