package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.GetNoteByIdUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class GetNoteByIdImpl implements GetNoteByIdUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public GetNoteByIdImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public Note execute(UUID projectId, UUID id) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        Note note = noteRepositoryPort.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Note", id));

        if(!note.getProjectId().equals(projectId)) throw new ResourceNotFoundException("Note", id);

        return note;
    }
}
