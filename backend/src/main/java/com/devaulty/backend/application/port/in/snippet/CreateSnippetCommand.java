package com.devaulty.backend.application.port.in.snippet;

import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;

import java.util.UUID;

public record CreateSnippetCommand(
        UUID projectId,
        String title,
        String description,
        String content,
        SnippetLanguage language,
        SnippetType snippetType
) {
}
