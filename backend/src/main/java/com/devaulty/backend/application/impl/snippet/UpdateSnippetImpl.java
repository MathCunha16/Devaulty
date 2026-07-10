package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.UpdateSnippetCommand;
import com.devaulty.backend.application.port.in.snippet.UpdateSnippetUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateSnippetImpl implements UpdateSnippetUseCase {

    private final SnippetRepositoryPort snippetRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateSnippetImpl(SnippetRepositoryPort snippetRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.snippetRepositoryPort = snippetRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Snippet execute(UpdateSnippetCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Snippet snippet = snippetRepositoryPort.findById(command.id())
                .orElseThrow(() -> new ResourceNotFoundException("Snippet", command.id()));

        if(!snippet.getProjectId().equals(command.projectId())) throw new ResourceNotFoundException("Snippet", command.id());

        if (command.title() != null ) snippet.setTitle(command.title());
        if (command.description() != null ) snippet.setDescription(command.description());
        if (command.content() != null ) snippet.setContent(command.content());
        if (command.language() != null ) snippet.setLanguage(command.language());
        if (command.snippetType() != null ) snippet.setSnippetType(command.snippetType());
        snippet.setUpdatedAt(LocalDateTime.now());

        return snippetRepositoryPort.save(snippet);
    }
}
