package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
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
class DeleteProblemImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @InjectMocks
    private DeleteProblemImpl deleteProblemUseCase;

    @Test
    void shouldDeleteProblemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.existsByIdAndProjectId(problemId, projectId)).thenReturn(true);

        // Act
        deleteProblemUseCase.execute(projectId, problemId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).existsByIdAndProjectId(problemId, projectId);
        verify(itemTagRepository, times(1)).removeAllTagsFromItem("problem", problemId);
        verify(problemRepository, times(1)).deleteById(problemId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteProblemUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, never()).existsByIdAndProjectId(any(UUID.class), any(UUID.class));
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.existsByIdAndProjectId(problemId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteProblemUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).existsByIdAndProjectId(problemId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        // existsByIdAndProjectId returns false because it belongs to another project
        when(problemRepository.existsByIdAndProjectId(problemId, projectId)).thenReturn(false);

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            deleteProblemUseCase.execute(projectId, problemId);
        });
        assertEquals("Problem not found with identifier " + problemId, exception.getMessage());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).existsByIdAndProjectId(problemId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }
}
