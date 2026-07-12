package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.GetAllProblemsByProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.data.domain.Page;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetAllProblemsByProjectImpl implements GetAllProblemsByProjectUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetAllProblemsByProjectImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Problem> execute(UUID projectId, int page, int size) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return problemRepositoryPort.findAllByProject(projectId, page, size);
    }
}
