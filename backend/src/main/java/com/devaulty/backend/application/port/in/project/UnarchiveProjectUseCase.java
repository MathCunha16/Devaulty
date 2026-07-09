package com.devaulty.backend.application.port.in.project;

import java.util.UUID;

public interface UnarchiveProjectUseCase {
    void execute(UUID id);
}
