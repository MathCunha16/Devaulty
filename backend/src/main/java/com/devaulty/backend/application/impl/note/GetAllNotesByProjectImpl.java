package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.GetAllNotesByProjectUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.data.domain.Page;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetAllNotesByProjectImpl implements GetAllNotesByProjectUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetAllNotesByProjectImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Page<Note> execute(UUID projectId, int page, int size) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        return noteRepositoryPort.findAllByProject(projectId, page, size);
    }
}
