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
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;

import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetAllProblemsByProjectImplTest {

    @Mock
    private ProblemRepositoryPort problemRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetAllProblemsByProjectImpl getAllProblemsUseCase;

    @Test
    void shouldReturnPageOfProblems() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        int page = 0;
        int size = 10;
        Problem problem = new Problem();
        problem.setId(UUID.randomUUID());
        problem.setProjectId(projectId);
        problem.setTitle("Title");
        problem.setStatus(ProblemStatus.OPEN);
        problem.setSeverity(ProblemSeverity.MEDIUM);
        List<Problem> list = Collections.singletonList(problem);
        Page<Problem> expectedPage = new PageImpl<>(list);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(problemRepository.findAllByProject(projectId, page, size)).thenReturn(expectedPage);

        // Act
        Page<Problem> result = getAllProblemsUseCase.execute(projectId, page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(problem, result.getContent().getFirst());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, times(1)).findAllByProject(projectId, page, size);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllProblemsUseCase.execute(projectId, 0, 10);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(problemRepository, never()).findAllByProject(any(UUID.class), anyInt(), anyInt());
    }
}
