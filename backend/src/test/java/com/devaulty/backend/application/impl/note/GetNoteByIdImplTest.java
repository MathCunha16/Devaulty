package com.devaulty.backend.application.impl.note;

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
class GetNoteByIdImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetNoteByIdImpl getNoteByIdUseCase;

    @Test
    void shouldReturnNoteWhenFoundAndBelongsToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note expectedNote = new Note(noteId, projectId, "Title", "Content", false);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(expectedNote));

        // Act
        Note result = getNoteByIdUseCase.execute(projectId, noteId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedNote, result);

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getNoteByIdUseCase.execute(projectId, noteId);
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
            getNoteByIdUseCase.execute(projectId, noteId);
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
        Note noteOfOtherProject = new Note(noteId, otherProjectId, "Title", "Content", false);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(noteOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getNoteByIdUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
    }
}
