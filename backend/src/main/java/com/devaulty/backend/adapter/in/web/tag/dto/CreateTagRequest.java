package com.devaulty.backend.adapter.in.web.tag.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record CreateTagRequest(
        @NotBlank(message = "Name cannot be blank")
        @Size(min = 1, max = 40, message = "Name must be between 1 and 40 characters")
        String name,
        @Pattern(regexp = "^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$", message = "Color must be a valid hex code (e.g. #1A1A2E or #FFF)")
        String color
) {
}
