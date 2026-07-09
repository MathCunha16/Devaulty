package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetCommand;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateSnippetImpl implements CreateSnippetUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public CreateSnippetImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Snippet execute(CreateSnippetCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())){
            throw new ResourceNotFoundException("Project", command.projectId());
        }

        Snippet snippet = new Snippet();
        snippet.setId(UUID.randomUUID());
        snippet.setProjectId(command.projectId());
        snippet.setTitle(command.title());
        snippet.setDescription(command.description());
        snippet.setContent(command.content());
        snippet.setLanguage(command.language());
        snippet.setSnippetType(command.snippetType());
        snippet.setCreatedAt(LocalDateTime.now());

        return snippetRepositoryPort.save(snippet);
    }
}
