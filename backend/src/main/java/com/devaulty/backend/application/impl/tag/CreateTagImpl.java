package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceAlreadyExistsException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.CreateTagCommand;
import com.devaulty.backend.application.port.in.tag.CreateTagUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateTagImpl implements CreateTagUseCase {

    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public CreateTagImpl(TagRepositoryPort tagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Tag execute(CreateTagCommand command) {

        if (!projectRepositoryPort.existsById(command.projectId())) {
            throw new ResourceNotFoundException("Project", command.projectId());
        }

        String sanitizedName = command.name().trim();

        if (tagRepositoryPort.existsByNameAndProjectId(command.projectId(), sanitizedName)) {
            throw new ResourceAlreadyExistsException("Tag", sanitizedName);
        }

        Tag tag = new Tag();
        tag.setId(UUID.randomUUID());
        tag.setProjectId(command.projectId());
        tag.setName(command.name());
        tag.setColor(command.color());
        tag.setCreatedAt(LocalDateTime.now());

        return tagRepositoryPort.save(tag);
    }
}
