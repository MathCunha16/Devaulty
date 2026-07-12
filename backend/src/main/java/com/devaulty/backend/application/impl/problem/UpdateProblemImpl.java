package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.UpdateProblemCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateProblemImpl implements UpdateProblemUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateProblemImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Problem execute(UpdateProblemCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Problem problem = problemRepositoryPort.findById(command.id())
                .orElseThrow(() -> new ResourceNotFoundException("Problem", command.id()));

        if(!problem.getProjectId().equals(command.projectId())) throw new ResourceNotFoundException("Problem", command.id());

        if (command.title() != null ) problem.setTitle(command.title());
        if (command.errorDescription() != null ) problem.setErrorDescription(command.errorDescription());
        if (command.solution() != null ) problem.setSolution(command.solution());
        if (command.severity() != null ) problem.setSeverity(command.severity());
        problem.setUpdatedAt(LocalDateTime.now());

        return problemRepositoryPort.save(problem);
    }
}
