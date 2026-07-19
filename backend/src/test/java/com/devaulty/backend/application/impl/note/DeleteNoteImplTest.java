package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DeleteNoteImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @InjectMocks
    private DeleteNoteImpl deleteNoteUseCase;

    @Test
    void shouldDeleteNoteSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.existsByIdAndProjectId(noteId, projectId)).thenReturn(true);

        // Act
        deleteNoteUseCase.execute(projectId, noteId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).existsByIdAndProjectId(noteId, projectId);
        verify(itemTagRepository, times(1)).removeAllTagsFromItem("note", noteId);
        verify(noteRepository, times(1)).deleteById(noteId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, never()).existsByIdAndProjectId(any(UUID.class), any(UUID.class));
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(noteRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNoteDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.existsByIdAndProjectId(noteId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteNoteUseCase.execute(projectId, noteId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).existsByIdAndProjectId(noteId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(noteRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNoteDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        // existsByIdAndProjectId returns false because it belongs to another project
        when(noteRepository.existsByIdAndProjectId(noteId, projectId)).thenReturn(false);

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            deleteNoteUseCase.execute(projectId, noteId);
        });
        assertEquals("Note not found with identifier " + noteId, exception.getMessage());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).existsByIdAndProjectId(noteId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(noteRepository, never()).deleteById(any(UUID.class));
    }
}
