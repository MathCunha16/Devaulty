package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.UpdateLinkCommand;
import com.devaulty.backend.application.port.in.link.UpdateLinkUseCase;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateLinkImpl implements UpdateLinkUseCase {

    private final LinkRepositoryPort linkRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateLinkImpl(LinkRepositoryPort linkRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.linkRepositoryPort = linkRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Link execute(UpdateLinkCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Link link = linkRepositoryPort.findById(command.id())
                .orElseThrow(() -> new ResourceNotFoundException("Link", command.id()));

        if(!link.getProjectId().equals(command.projectId())) throw new ResourceNotFoundException("Link", command.id());

        if(command.title() != null) link.setTitle(command.title());
        if(command.url() != null) link.setUrl(command.url());
        if(command.description() != null) link.setDescription(command.description());
        link.setUpdatedAt(LocalDateTime.now());

        return linkRepositoryPort.save(link);
    }
}
