package com.devaulty.backend.application.port.in.note;

import com.devaulty.backend.domain.model.Note;

public interface UpdateNoteUseCase {
    Note execute(UpdateNoteCommand command);
}
