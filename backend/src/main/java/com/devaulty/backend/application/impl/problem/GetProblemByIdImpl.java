package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.GetProblemByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetProblemByIdImpl implements GetProblemByIdUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetProblemByIdImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Problem execute(UUID projectId, UUID id) {

        if (!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", id);

        Problem problem = problemRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Problem", id));

        if(!problem.getProjectId().equals(projectId)) throw new ResourceNotFoundException("Problem", id);

        return problem;
    }
}
