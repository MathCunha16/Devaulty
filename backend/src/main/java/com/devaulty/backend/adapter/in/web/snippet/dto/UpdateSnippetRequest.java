package com.devaulty.backend.adapter.in.web.snippet.dto;

import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
import jakarta.validation.constraints.Size;

public record UpdateSnippetRequest(
        @Size(min = 2, max = 255, message = "Title must be between 2 and 255 characters")
        String title,
        String description,
        String content,
        SnippetLanguage language,
        SnippetType snippetType
) {
}
