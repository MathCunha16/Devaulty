package com.devaulty.backend.application.port.in.note;

import com.devaulty.backend.domain.model.Note;
import org.springframework.data.domain.Page;

import java.util.UUID;

public interface GetAllNotesByProjectUseCase {
    Page<Note> execute(UUID projectId, int page, int size);
}
