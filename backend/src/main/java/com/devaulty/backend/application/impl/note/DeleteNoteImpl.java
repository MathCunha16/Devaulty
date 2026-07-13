package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.DeleteNoteUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteNoteImpl implements DeleteNoteUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public DeleteNoteImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional(readOnly = true)
    public void execute(UUID projectId, UUID id) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        if (noteRepositoryPort.findById(id)
                .filter(note -> projectId.equals(note.getProjectId()))
                .isEmpty()) {
            throw new ResourceNotFoundException("Note", id);
        }

        noteRepositoryPort.deleteById(id);
    }
}
