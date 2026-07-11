package com.devaulty.backend.adapter.in.web.link.dto;

import jakarta.validation.constraints.Pattern;

public record UpdateLinkRequest(
        @Pattern(regexp = "(?s).*\\S.*", message = "Title must not be blank")
        String title,
        @Pattern(regexp = "(?s).*\\S.*", message = "url must not be blank")
        String url,
        String description
) {
}
