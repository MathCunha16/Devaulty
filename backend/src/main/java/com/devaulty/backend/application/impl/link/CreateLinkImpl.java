package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.CreateLinkCommand;
import com.devaulty.backend.application.port.in.link.CreateLinkUseCase;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateLinkImpl implements CreateLinkUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public CreateLinkImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Link execute(CreateLinkCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Link link = new Link();
        link.setId(UUID.randomUUID());
        link.setProjectId(command.projectId());
        link.setTitle(command.title());
        link.setUrl(command.url());
        link.setDescription(command.description());
        link.setCreatedAt(LocalDateTime.now());

        return linkRepositoryPort.save(link);
    }
}
