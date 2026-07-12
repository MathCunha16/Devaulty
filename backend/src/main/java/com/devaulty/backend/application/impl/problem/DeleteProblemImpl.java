package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.DeleteProblemUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteProblemImpl implements DeleteProblemUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public DeleteProblemImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        if (problemRepositoryPort.findById(id)
                .filter(problem -> projectId.equals(problem.getProjectId()))
                .isEmpty()) {
            throw new ResourceNotFoundException("Problem", id);
        }

        problemRepositoryPort.deleteById(id);

    }
}
