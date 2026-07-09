package com.devaulty.backend.application.port.in.snippet;

import com.devaulty.backend.domain.model.Snippet;
import org.springframework.data.domain.Page;

import java.util.UUID;

public interface GetAllSnippetsByProjectUseCase {
    Page<Snippet> execute(UUID projectId, int page, int size);
}
