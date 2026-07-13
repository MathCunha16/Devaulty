package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.BusinessRuleException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.UnarchiveNoteUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class UnarchiveNoteImpl implements UnarchiveNoteUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UnarchiveNoteImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {

        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);

        Note note = noteRepositoryPort.findById(id)
                .filter(n -> projectId.equals(n.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Note", id));

        if(!note.isArchived()){
            throw new BusinessRuleException("Note not archived");
        }

        note.setArchived(false);
        note.setUpdatedAt(LocalDateTime.now());

        noteRepositoryPort.save(note);
    }
}
