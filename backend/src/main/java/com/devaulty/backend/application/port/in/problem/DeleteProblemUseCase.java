package com.devaulty.backend.application.port.in.problem;

import java.util.UUID;

public interface DeleteProblemUseCase {
    void execute(UUID projectId, UUID id);
}
