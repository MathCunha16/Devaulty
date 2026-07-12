package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
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
class DeleteProblemImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private DeleteProblemImpl deleteProblemUseCase;

    @Test
    void shouldDeleteProblemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        Problem problem = new Problem();
        problem.setId(problemId);
        problem.setProjectId(projectId);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(problem));

        // Act
        deleteProblemUseCase.execute(projectId, problemId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
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
        verify(problemRepository, never()).findById(any(UUID.class));
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteProblemUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        Problem problem = new Problem();
        problem.setId(problemId);
        problem.setProjectId(otherProjectId);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(problem));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteProblemUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
        verify(problemRepository, never()).deleteById(any(UUID.class));
    }
}
