package com.devaulty.backend.application.port.in.project;

import java.util.UUID;

public interface ArchiveProjectUseCase {
    void execute(UUID id);
}
