package com.devaulty.backend.adapter.out.persistence.note;

import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Component
public class NotePersistenceAdapter implements NoteRepositoryPort {

    private final SpringDataNoteRepository noteRepository;
    private final SpringDataProjectRepository projectRepository;
    private final NoteMapper noteMapper;

    public NotePersistenceAdapter(SpringDataNoteRepository noteRepository, SpringDataProjectRepository projectRepository, NoteMapper noteMapper) {
        this.noteRepository = noteRepository;
        this.projectRepository = projectRepository;
        this.noteMapper = noteMapper;
    }

    @Override
    public Note save(Note note) {
        NoteEntity entity = noteMapper.toEntity(note);
        entity.setProject(projectRepository.getReferenceById(note.getProjectId()));
        NoteEntity savedEntity = noteRepository.save(entity);
        return noteMapper.toDomain(savedEntity);
    }

    @Override
    public Optional<Note> findById(UUID id) {
        Optional<NoteEntity> entity = noteRepository.findById(id);
        return entity.map(noteMapper::toDomain);
    }

    @Override
    public Page<Note> findAllByProject(UUID projectId, int page, int size) {
        Pageable pageable = PageRequest.of(page, size, Sort.by(Sort.Direction.DESC, "createdAt"));
        Page<NoteEntity> entities = noteRepository.findAllByProject_Id(projectId, pageable);
        return entities.map(noteMapper::toDomain);
    }


    @Override
    public void deleteById(UUID id) {
        noteRepository.deleteById(id);
    }

    @Override
    public boolean existsByIdAndProjectId(UUID id, UUID projectId) {
        return noteRepository.existsByIdAndProject_Id(id, projectId);
    }

    @Override
    public String getSupportedType() {
        return "note";
    }

    @Override
    public List<UUID> findExistingIdsByProject(List<UUID> ids, UUID projectId) {
        if (ids == null || ids.isEmpty()) {
            return List.of();
        }
        return noteRepository.findExistingIdsByProjectId(ids, projectId);
    }
}
