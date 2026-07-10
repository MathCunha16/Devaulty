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
class DeleteSnippetImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private DeleteSnippetImpl deleteSnippetUseCase;

    @Test
    void shouldDeleteSnippetSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet snippet = new Snippet(snippetId, projectId, "Title", "Desc", "content", SnippetLanguage.JAVA, SnippetType.CODE);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(snippet));

        // Act
        deleteSnippetUseCase.execute(projectId, snippetId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, times(1)).deleteById(snippetId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteSnippetUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, never()).findById(any(UUID.class));
        verify(snippetRepository, never()).deleteById(any(UUID.class));
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
            deleteSnippetUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        Snippet snippet = new Snippet(snippetId, otherProjectId, "Title", "Desc", "content", SnippetLanguage.JAVA, SnippetType.CODE);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findById(snippetId)).thenReturn(Optional.of(snippet));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteSnippetUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findById(snippetId);
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }
}
