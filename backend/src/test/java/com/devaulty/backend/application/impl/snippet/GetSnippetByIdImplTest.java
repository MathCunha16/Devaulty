package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
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
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetSnippetByIdImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetSnippetByIdImpl getSnippetByIdUseCase;

    @Test
    void shouldReturnSnippetWhenFoundAndBelongsToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet expectedSnippet = new Snippet(snippetId, projectId, "Title", "Desc", "ls", SnippetLanguage.BASH, SnippetType.COMMAND);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(expectedSnippet));

        // Act
        Snippet result = getSnippetByIdUseCase.execute(projectId, snippetId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedSnippet, result);

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getSnippetByIdUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getSnippetByIdUseCase.execute(projectId, snippetId);
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

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(snippetOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getSnippetByIdUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
    }
}
