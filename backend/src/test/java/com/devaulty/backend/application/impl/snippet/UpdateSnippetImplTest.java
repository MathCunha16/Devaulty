package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.UpdateSnippetCommand;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
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
class UpdateSnippetImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private UpdateSnippetImpl updateSnippetUseCase;

    @Test
    void shouldUpdateFieldsWhenValidCommand() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet existingSnippet = new Snippet(snippetId, projectId, "Old Title", "Old Desc", "old content", SnippetLanguage.JAVA, SnippetType.CODE);
        UpdateSnippetCommand command = new UpdateSnippetCommand(
                snippetId,
                projectId,
                "New Title",
                "New Desc",
                "new content",
                SnippetLanguage.KOTLIN,
                SnippetType.CODE
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(existingSnippet));
        when(snippetRepository.save(any(Snippet.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Snippet result = updateSnippetUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("New Title", result.getTitle());
        assertEquals("New Desc", result.getDescription());
        assertEquals("new content", result.getContent());
        assertEquals(SnippetLanguage.KOTLIN, result.getLanguage());
        assertEquals(SnippetType.CODE, result.getSnippetType());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, times(1)).save(existingSnippet);
    }

    @Test
    void shouldNotUpdateFieldsWhenNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet existingSnippet = new Snippet(snippetId, projectId, "Old Title", "Old Desc", "old content", SnippetLanguage.JAVA, SnippetType.CODE);
        UpdateSnippetCommand command = new UpdateSnippetCommand(
                snippetId,
                projectId,
                null,
                null,
                null,
                null,
                null
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(existingSnippet));
        when(snippetRepository.save(any(Snippet.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Snippet result = updateSnippetUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("Old Title", result.getTitle());
        assertEquals("Old Desc", result.getDescription());
        assertEquals("old content", result.getContent());
        assertEquals(SnippetLanguage.JAVA, result.getLanguage());
        assertEquals(SnippetType.CODE, result.getSnippetType());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, times(1)).save(existingSnippet);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        UpdateSnippetCommand command = new UpdateSnippetCommand(snippetId, projectId, "Title", null, null, null, null);

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateSnippetUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        UpdateSnippetCommand command = new UpdateSnippetCommand(snippetId, projectId, "Title", null, null, null, null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateSnippetUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet snippetOfOtherProject = new Snippet(snippetId, otherProjectId, "Title", "Desc", "ls", SnippetLanguage.BASH, SnippetType.COMMAND);
        UpdateSnippetCommand command = new UpdateSnippetCommand(snippetId, projectId, "Title", null, null, null, null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(snippetOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateSnippetUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, never()).save(any(Snippet.class));
    }
}
