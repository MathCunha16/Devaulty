package com.devaulty.backend.application.port.in.note;

import java.util.UUID;

public interface DeleteNoteUseCase {
    void execute(UUID projectId, UUID id);
}
