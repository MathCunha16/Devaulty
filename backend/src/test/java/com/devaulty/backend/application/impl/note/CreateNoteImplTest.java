package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.note.CreateNoteCommand;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateNoteImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateNoteImpl createNoteUseCase;

    @Test
    void shouldCreateNoteSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateNoteCommand command = new CreateNoteCommand(
                projectId,
                "Test Note",
                "Test Content"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.save(any(Note.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Note result = createNoteUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals(projectId, result.getProjectId());
        assertEquals("Test Note", result.getTitle());
        assertEquals("Test Content", result.getContent());
        assertFalse(result.isArchived());
        assertNotNull(result.getCreatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).save(any(Note.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateNoteCommand command = new CreateNoteCommand(
                projectId,
                "Test Note",
                "Test Content"
        );

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createNoteUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, never()).save(any(Note.class));
    }
}
