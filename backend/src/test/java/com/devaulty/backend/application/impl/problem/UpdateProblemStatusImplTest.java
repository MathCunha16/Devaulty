package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.UpdateProblemStatusCommand;
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
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UpdateProblemStatusImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private UpdateProblemStatusImpl updateProblemStatusUseCase;

    @Test
    void shouldUpdateStatusSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        Problem existingProblem = new Problem();
        existingProblem.setId(problemId);
        existingProblem.setProjectId(projectId);
        existingProblem.setTitle("Title");
        existingProblem.setStatus(ProblemStatus.OPEN);
        existingProblem.setSeverity(ProblemSeverity.MEDIUM);

        UpdateProblemStatusCommand command = new UpdateProblemStatusCommand(problemId, projectId, ProblemStatus.RESOLVED);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(existingProblem));
        when(problemRepository.save(any(Problem.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Problem result = updateProblemStatusUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals(ProblemStatus.RESOLVED, result.getStatus());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
        verify(problemRepository, times(1)).save(existingProblem);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        UpdateProblemStatusCommand command = new UpdateProblemStatusCommand(problemId, projectId, ProblemStatus.RESOLVED);

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateProblemStatusUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProblemDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();
        UpdateProblemStatusCommand command = new UpdateProblemStatusCommand(problemId, projectId, ProblemStatus.RESOLVED);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateProblemStatusUseCase.execute(command);
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
        UpdateProblemStatusCommand command = new UpdateProblemStatusCommand(problemId, projectId, ProblemStatus.RESOLVED);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findById(problemId)).thenReturn(Optional.of(problemOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateProblemStatusUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findById(problemId);
        verify(problemRepository, never()).save(any(Problem.class));
    }
}
