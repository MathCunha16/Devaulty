package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.DeleteProblemUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteProblemImpl implements DeleteProblemUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final ItemTagRepositoryPort itemTagRepositoryPort;

    public DeleteProblemImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort, ItemTagRepositoryPort itemTagRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.itemTagRepositoryPort = itemTagRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if (!problemRepositoryPort.existsByIdAndProjectId(id, projectId)) throw new ResourceNotFoundException("Problem", id);

        itemTagRepositoryPort.removeAllTagsFromItem("problem", id);
        problemRepositoryPort.deleteById(id);

    }
}
