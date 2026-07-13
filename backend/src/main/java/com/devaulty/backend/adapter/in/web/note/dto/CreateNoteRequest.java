package com.devaulty.backend.adapter.in.web.note.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record CreateNoteRequest(
        @NotBlank
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,

        String content
) {
}
