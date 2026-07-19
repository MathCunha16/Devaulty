package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
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
class DeleteSnippetImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @InjectMocks
    private DeleteSnippetImpl deleteSnippetUseCase;

    @Test
    void shouldDeleteSnippetSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.existsByIdAndProjectId(snippetId, projectId)).thenReturn(true);

        // Act
        deleteSnippetUseCase.execute(projectId, snippetId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(snippetId, projectId);
        verify(itemTagRepository, times(1)).removeAllTagsFromItem("snippet", snippetId);
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
        verify(snippetRepository, never()).existsByIdAndProjectId(any(UUID.class), any(UUID.class));
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.existsByIdAndProjectId(snippetId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteSnippetUseCase.execute(projectId, snippetId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(snippetId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenSnippetDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        // existsByIdAndProjectId returns false because it belongs to another project
        when(snippetRepository.existsByIdAndProjectId(snippetId, projectId)).thenReturn(false);

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            deleteSnippetUseCase.execute(projectId, snippetId);
        });
        assertEquals("Snippet not found with identifier " + snippetId, exception.getMessage());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(snippetId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(snippetRepository, never()).deleteById(any(UUID.class));
    }
}
