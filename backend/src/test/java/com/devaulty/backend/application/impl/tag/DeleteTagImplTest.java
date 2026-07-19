package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DeleteTagImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private DeleteTagImpl deleteTagUseCase;

    @Test
    void shouldDeleteTagSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(true);

        // Act
        deleteTagUseCase.execute(projectId, tagId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(tagRepository, times(1)).delete(tagId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteTagUseCase.execute(projectId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).existsByIdAndProjectId(any(), any());
        verify(tagRepository, never()).delete(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteTagUseCase.execute(projectId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(tagRepository, never()).delete(any());
    }
}
