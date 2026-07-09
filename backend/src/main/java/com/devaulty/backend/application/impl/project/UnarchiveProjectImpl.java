package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.BusinessRuleException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.UnarchiveProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class UnarchiveProjectImpl implements UnarchiveProjectUseCase {

    private final ProjectRepositoryPort projectRepositoryPort;

    public UnarchiveProjectImpl(ProjectRepositoryPort projectRepositoryPort) {
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID id) {
        Project project = projectRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Project", id));

        if(project.isArchived()){
            project.setArchived(false);
            project.setUpdatedAt(LocalDateTime.now());
            projectRepositoryPort.save(project);
        } else {
            throw new BusinessRuleException("Project already unarchived");
        }
    }
}
