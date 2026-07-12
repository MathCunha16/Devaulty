package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProblemRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Problem;
import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;
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
class GetProblemByIdImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetProblemByIdImpl getProblemByIdUseCase;

    @Test
    void shouldReturnProblemWhenFoundAndBelongsToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        Problem expectedProblem = new Problem();
        expectedProblem.setId(problemId);
        expectedProblem.setProjectId(projectId);
        expectedProblem.setTitle("Title");
        expectedProblem.setStatus(ProblemStatus.OPEN);
        expectedProblem.setSeverity(ProblemSeverity.MEDIUM);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(expectedProblem));

        // Act
        Problem result = getProblemByIdUseCase.execute(projectId, problemId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedProblem, result);

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getProblemByIdUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, never()).findById(any(UUID.class));
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
            getProblemByIdUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        Problem problemOfOtherProject = new Problem();
        problemOfOtherProject.setId(problemId);
        problemOfOtherProject.setProjectId(otherProjectId);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(problemOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getProblemByIdUseCase.execute(projectId, problemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
    }
}
