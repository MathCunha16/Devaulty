package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.BusinessRuleException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.ArchiveProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class ArchiveProjectImpl implements ArchiveProjectUseCase {

    private final ProjectRepositoryPort projectRepositoryPort;

    public ArchiveProjectImpl(ProjectRepositoryPort projectRepositoryPort) {
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID id) {
        Project project = projectRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Project", id));

        if (!project.isArchived()){
            project.setArchived(true);
            project.setUpdatedAt(LocalDateTime.now());
            projectRepositoryPort.save(project);
        } else {
            throw new BusinessRuleException("Project already archived");
        }
    }
}
