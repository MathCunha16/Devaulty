package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.snippet.CreateSnippetCommand;
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

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateSnippetImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateSnippetImpl createSnippetUseCase;

    @Test
    void shouldCreateSnippetSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateSnippetCommand command = new CreateSnippetCommand(
                projectId,
                "Test Title",
                "Test Desc",
                "echo 'test'",
                SnippetLanguage.BASH,
                SnippetType.COMMAND
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.save(any(Snippet.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Snippet result = createSnippetUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals(projectId, result.getProjectId());
        assertEquals("Test Title", result.getTitle());
        assertEquals("Test Desc", result.getDescription());
        assertEquals("echo 'test'", result.getContent());
        assertEquals(SnippetLanguage.BASH, result.getLanguage());
        assertEquals(SnippetType.COMMAND, result.getSnippetType());
        assertNotNull(result.getCreatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).save(any(Snippet.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateSnippetCommand command = new CreateSnippetCommand(
                projectId,
                "Test Title",
                "Test Desc",
                "echo 'test'",
                SnippetLanguage.BASH,
                SnippetType.COMMAND
        );

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createSnippetUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, never()).save(any(Snippet.class));
    }
}
