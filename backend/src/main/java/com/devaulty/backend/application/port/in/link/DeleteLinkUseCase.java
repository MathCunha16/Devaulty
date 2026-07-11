package com.devaulty.backend.application.port.in.link;

import java.util.UUID;

public interface DeleteLinkUseCase {
    void execute(UUID projectId, UUID id);
}
