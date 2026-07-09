package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.project.CreateProjectCommand;
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
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UpdateProjectImplTest {

    @Mock
    private ProjectRepositoryPort projectRepositoryPort;

    @InjectMocks
    private UpdateProjectImpl updateProjectUseCase;

    @Test
    void shouldUpdateAllFieldsWhenCommandHasAllValues() {
        // Arrange
        UUID id = UUID.randomUUID();
        Project existingProject = new Project(id, "Old Name", "Old Desc", "old-icon", "#000", false);
        CreateProjectCommand command = new CreateProjectCommand("New Name", "New Desc", "new-icon", "#fff");

        when(projectRepositoryPort.findById(id)).thenReturn(Optional.of(existingProject));
        when(projectRepositoryPort.save(any(Project.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Project result = updateProjectUseCase.execute(id, command);

        // Assert
        assertNotNull(result);
        assertEquals("New Name", result.getName());
        assertEquals("New Desc", result.getDescription());
        assertEquals("new-icon", result.getIcon());
        assertEquals("#fff", result.getColor());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, times(1)).save(existingProject);
    }

    @Test
    void shouldNotUpdateFieldsWhenCommandHasNullValues() {
        // Arrange
        UUID id = UUID.randomUUID();
        Project existingProject = new Project(id, "Old Name", "Old Desc", "old-icon", "#000", false);
        CreateProjectCommand command = new CreateProjectCommand(null, null, null, null);

        when(projectRepositoryPort.findById(id)).thenReturn(Optional.of(existingProject));
        when(projectRepositoryPort.save(any(Project.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Project result = updateProjectUseCase.execute(id, command);

        // Assert
        assertNotNull(result);
        assertEquals("Old Name", result.getName());
        assertEquals("Old Desc", result.getDescription());
        assertEquals("old-icon", result.getIcon());
        assertEquals("#000", result.getColor());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, times(1)).save(existingProject);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID id = UUID.randomUUID();
        CreateProjectCommand command = new CreateProjectCommand("New Name", "New Desc", "new-icon", "#fff");

        when(projectRepositoryPort.findById(id)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateProjectUseCase.execute(id, command);
        });

        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, never()).save(any(Project.class));
    }
}
