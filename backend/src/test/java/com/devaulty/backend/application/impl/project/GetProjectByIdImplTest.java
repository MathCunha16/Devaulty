package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
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
class GetProjectByIdImplTest {

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetProjectByIdImpl getProjectByIdUseCase;

    @Test
    void shouldReturnProjectWhenFound() {
        // Arrange
        UUID id = UUID.randomUUID();
        Project expectedProject = new Project(id, "Test Project", "Desc", "icon", "#fff", false);

        when(projectRepository.findById(id)).thenReturn(Optional.of(expectedProject));

        // Act
        Project result = getProjectByIdUseCase.execute(id);

        // Assert
        assertNotNull(result);
        assertEquals(expectedProject, result);
        verify(projectRepository, times(1)).findById(id);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenNotFound() {
        // Arrange
        UUID id = UUID.randomUUID();
        when(projectRepository.findById(id)).thenReturn(Optional.empty());

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            getProjectByIdUseCase.execute(id);
        });

        assertEquals("Project not found with identifier " + id, exception.getMessage());
        verify(projectRepository, times(1)).findById(id);
    }
}
