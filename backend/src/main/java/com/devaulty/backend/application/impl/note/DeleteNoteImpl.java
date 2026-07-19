package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.DeleteNoteUseCase;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

public class DeleteNoteImpl implements DeleteNoteUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;
    private final ItemTagRepositoryPort itemTagRepositoryPort;

    public DeleteNoteImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort, ItemTagRepositoryPort itemTagRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
        this.itemTagRepositoryPort = itemTagRepositoryPort;
    }

    @Override
    @Transactional
    public void execute(UUID projectId, UUID id) {
        if(!projectRepositoryPort.existsById(projectId)) throw new ResourceNotFoundException("Project", projectId);
        if (!noteRepositoryPort.existsByIdAndProjectId(id, projectId)) throw new ResourceNotFoundException("Note", id);

        itemTagRepositoryPort.removeAllTagsFromItem("note", id);
        noteRepositoryPort.deleteById(id);
    }
}
