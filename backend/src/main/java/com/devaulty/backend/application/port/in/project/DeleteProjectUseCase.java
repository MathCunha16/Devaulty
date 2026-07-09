package com.devaulty.backend.application.port.in.project;

import java.util.UUID;

public interface DeleteProjectUseCase {
    void execute(UUID id);
}
