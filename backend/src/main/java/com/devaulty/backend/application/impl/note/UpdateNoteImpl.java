package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.UpdateNoteCommand;
import com.devaulty.backend.application.port.in.note.UpdateNoteUseCase;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;

public class UpdateNoteImpl implements UpdateNoteUseCase {

    private final NoteRepositoryPort noteRepositoryPort;
    private final ProjectRepositoryPort projectRepositoryPort;

    public UpdateNoteImpl(NoteRepositoryPort noteRepositoryPort, ProjectRepositoryPort projectRepositoryPort) {
        this.noteRepositoryPort = noteRepositoryPort;
        this.projectRepositoryPort = projectRepositoryPort;
    }

    @Override
    @Transactional
    public Note execute(UpdateNoteCommand command) {

        if(!projectRepositoryPort.existsById(command.projectId())) throw new ResourceNotFoundException("Project", command.projectId());

        Note note = noteRepositoryPort.findById(command.id())
                .filter(n -> command.projectId().equals(n.getProjectId()))
                .orElseThrow(() -> new ResourceNotFoundException("Note", command.id()));

        if (command.title() != null ) note.setTitle(command.title());
        if (command.content() != null ) note.setContent(command.content());
        note.setUpdatedAt(LocalDateTime.now());

        return noteRepositoryPort.save(note);
    }
}
