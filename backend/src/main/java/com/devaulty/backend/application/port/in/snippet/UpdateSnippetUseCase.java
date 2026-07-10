package com.devaulty.backend.application.port.in.snippet;

import com.devaulty.backend.domain.model.Snippet;

public interface UpdateSnippetUseCase {
    Snippet execute(UpdateSnippetCommand command);
}
