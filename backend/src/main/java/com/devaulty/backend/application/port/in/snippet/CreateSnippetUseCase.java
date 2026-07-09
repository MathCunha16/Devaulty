package com.devaulty.backend.application.port.in.snippet;

import com.devaulty.backend.domain.model.Snippet;

public interface CreateSnippetUseCase {
    Snippet execute(CreateSnippetCommand command);
}
