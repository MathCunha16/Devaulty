package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.CreateNoteCommand;
import com.devaulty.backend.application.port.in.note.CreateNoteUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

public class CreateNoteImpl implements CreateNoteUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public CreateNoteImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Note execute(CreateNoteCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Note note = new Note();
        note.setId(UUID.randomUUID());
        note.setProjectId(command.projectId());
        note.setTitle(command.title());
        note.setContent(command.content());
        note.setArchived(false);
        note.setCreatedAt(LocalDateTime.now());

        return noteRepositoryPort.save(note);
    }
}
