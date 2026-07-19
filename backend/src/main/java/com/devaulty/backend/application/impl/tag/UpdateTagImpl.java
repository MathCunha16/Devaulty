package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceAlreadyExistsException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.UpdateTagCommand;
import com.devaulty.backend.application.port.in.tag.UpdateTagUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateTagImpl implements UpdateTagUseCase {

    private final TagRepositoryPort tagRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateTagImpl(TagRepositoryPort tagRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.tagRepositoryPort = tagRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Tag execute(UpdateTagCommand command) {

        if (!projectRepositoryPort.existsById(command.projectId()))
            throw new ResourceNotFoundException("Project", command.projectId());

        Tag tag = tagRepositoryPort.findById(command.id())
                .filter(t -> command.projectId().equals(t.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Tag", command.id()));


        if (command.name() != null) {

            String sanitizedName = command.name().trim();
            if (tagRepositoryPort.existsByNameAndProjectId(command.projectId(), sanitizedName) && !tag.getName().equals(sanitizedName)) {
                throw new ResourceAlreadyExistsException("Tag", sanitizedName);
            }
            tag.setName(sanitizedName);

        }
        if (command.color() != null) tag.setColor(command.color());
        tag.setUpdatedAt(LocalDateTime.now());

        return tagRepositoryPort.save(tag);

    }
}
