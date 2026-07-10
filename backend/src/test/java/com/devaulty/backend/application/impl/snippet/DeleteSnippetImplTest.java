package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

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

        when(projectRepository.existsById(projectId)).thenReturn(true);

        // Act
        deleteSnippetUseCase.execute(projectId, snippetId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
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
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }
}
