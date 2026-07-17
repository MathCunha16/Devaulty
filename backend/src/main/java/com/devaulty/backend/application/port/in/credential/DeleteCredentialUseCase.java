package com.devaulty.backend.application.port.in.credential;

import java.util.UUID;

public interface DeleteCredentialUseCase {
    void execute(UUID projectId, UUID id);
}
