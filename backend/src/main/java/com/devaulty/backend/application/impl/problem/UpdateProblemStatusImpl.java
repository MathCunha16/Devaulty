package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.UpdateProblemStatusCommand;
import com.devaulty.backend.application.port.in.problem.UpdateProblemStatusUseCase;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateProblemStatusImpl implements UpdateProblemStatusUseCase {

    private final ProblemRepositoryPort problemRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateProblemStatusImpl(ProblemRepositoryPort problemRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.problemRepositoryPort = problemRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Problem execute(UpdateProblemStatusCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Problem problem = problemRepositoryPort.findById(command.id())
                .orElseThrow(() -> new ResourceNotFoundException("Problem", command.id()));

        if(!problem.getProjectId().equals(command.projectId())) throw new ResourceNotFoundException("Problem", command.id());

        problem.setStatus(command.status());
        problem.setUpdatedAt(LocalDateTime.now());

        return problemRepositoryPort.save(problem);
    }
}
