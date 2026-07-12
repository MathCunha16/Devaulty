package com.devaulty.backend.application.port.in.problem;

import com.devaulty.backend.domain.model.Problem;

import java.util.UUID;

public interface GetProblemByIdUseCase {
    Problem execute(UUID projectId,UUID id);
}
