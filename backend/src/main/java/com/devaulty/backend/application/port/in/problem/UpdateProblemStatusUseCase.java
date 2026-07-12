package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.Problem;

public interface UpdateProblemStatusUseCase {
    Problem execute(UpdateProblemStatusCommand command);
}
