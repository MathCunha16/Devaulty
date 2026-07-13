package com.devaulty.backend.application.port.in.note;

import java.util.UUID;

public interface UnarchiveNoteUseCase {
    void execute(UUID projectId ,UUID id);
}
