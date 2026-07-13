package com.devaulty.backend.application.port.in.note;

import com.devaulty.backend.domain.model.Note;

import java.util.UUID;

public interface GetNoteByIdUseCase {
    Note execute(UUID projectId ,UUID id);
}
