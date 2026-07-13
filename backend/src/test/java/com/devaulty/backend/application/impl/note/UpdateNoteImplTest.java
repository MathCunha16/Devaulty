package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.UpdateNoteCommand;
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
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UpdateNoteImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private UpdateNoteImpl updateNoteUseCase;

    @Test
    void shouldUpdateFieldsWhenValidCommand() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note existingNote = new Note(noteId, projectId, "Old Title", "Old Content", false);
        UpdateNoteCommand command = new UpdateNoteCommand(
                noteId,
                projectId,
                "New Title",
                "New Content"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(existingNote));
        when(noteRepository.save(any(Note.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Note result = updateNoteUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("New Title", result.getTitle());
        assertEquals("New Content", result.getContent());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, times(1)).save(existingNote);
    }

    @Test
    void shouldNotUpdateFieldsWhenNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        Note existingNote = new Note(noteId, projectId, "Old Title", "Old Content", false);
        UpdateNoteCommand command = new UpdateNoteCommand(
                noteId,
                projectId,
                null,
                null
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(existingNote));
        when(noteRepository.save(any(Note.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Note result = updateNoteUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("Old Title", result.getTitle());
        assertEquals("Old Content", result.getContent());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, times(1)).save(existingNote);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        UpdateNoteCommand command = new UpdateNoteCommand(noteId, projectId, "Title", null);

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateNoteUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNoteDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();
        UpdateNoteCommand command = new UpdateNoteCommand(noteId, projectId, "Title", null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateNoteUseCase.execute(command);
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
        UpdateNoteCommand command = new UpdateNoteCommand(noteId, projectId, "Title", null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findById(noteId)).thenReturn(Optional.of(noteOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateNoteUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findById(noteId);
        verify(noteRepository, never()).save(any(Note.class));
    }
}
