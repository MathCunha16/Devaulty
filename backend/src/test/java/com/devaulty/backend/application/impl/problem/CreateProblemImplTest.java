package com.devaulty.backend.application.impl.problem;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.problem.CreateProblemCommand;
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

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateProblemImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateProblemImpl createProblemUseCase;

    @Test
    void shouldCreateProblemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateProblemCommand command = new CreateProblemCommand(
                projectId,
                "Test Title",
                "Error Description",
                "Solution",
                ProblemStatus.OPEN,
                ProblemSeverity.HIGH
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.save(any(Problem.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Problem result = createProblemUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals(projectId, result.getProjectId());
        assertEquals("Test Title", result.getTitle());
        assertEquals("Error Description", result.getErrorDescription());
        assertEquals("Solution", result.getSolution());
        assertEquals(ProblemStatus.OPEN, result.getStatus());
        assertEquals(ProblemSeverity.HIGH, result.getSeverity());
        assertNotNull(result.getCreatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).save(any(Problem.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateProblemCommand command = new CreateProblemCommand(
                projectId,
                "Test Title",
                "Error Description",
                "Solution",
                ProblemStatus.OPEN,
                ProblemSeverity.HIGH
        );

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createProblemUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, never()).save(any(Problem.class));
    }
}
