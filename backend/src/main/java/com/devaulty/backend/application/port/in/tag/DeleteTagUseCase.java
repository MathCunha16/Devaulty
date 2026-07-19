package com.devaulty.backend.application.port.in.tag;

import java.util.UUID;

public interface DeleteTagUseCase {
    void execute(UUID projectId, UUID id);
}
