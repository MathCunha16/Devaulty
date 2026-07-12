package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.CreateProblemCommand;
import com.devaulty.backend.application.port.in.problem.CreateProblemUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateProblemImpl implements CreateProblemUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public CreateProblemImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Problem execute(CreateProblemCommand command) {

        if (!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Problem problem = new Problem();
        problem.setId(UUID.randomUUID());
        problem.setProjectId(command.projectId());
        problem.setTitle(command.title());
        problem.setErrorDescription(command.errorDescription());
        problem.setSolution(command.solution());
        problem.setStatus(command.status());
        problem.setSeverity(command.severity());
        problem.setCreatedAt(LocalDateTime.now());

        return problemRepositoryPort.save(problem);
    }
}
