package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.BusinessRuleException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ArchiveNoteImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private ArchiveNoteImpl archiveNoteUseCase;

    @Test
    void shouldArchiveNoteSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note note = new Note(noteId, projectId, "Title", "Content", false);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(note));

        // Act
        archiveNoteUseCase.execute(projectId, noteId);

        // Assert
        assertTrue(note.isArchived());
        assertNotNull(note.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, times(1)).save(note);
    }

    @Test
    void shouldThrowBusinessRuleExceptionWhenAlreadyArchived() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note note = new Note(noteId, projectId, "Title", "Content", true);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(note));

        // Act & Assert
        assertThrows(BusinessRuleException.class, () -> {
            archiveNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, never()).save(any(Note.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            archiveNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNoteDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            archiveNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNoteDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note note = new Note(noteId, otherProjectId, "Title", "Content", false);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(note));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            archiveNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, never()).save(any(Note.class));
    }
}
