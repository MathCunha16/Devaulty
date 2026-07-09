package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.port.in.project.CreateProjectCommand;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Project;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateProjectImplTest {

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateProjectImpl createProjectUseCase;

    @Test
    void shouldCreateProjectSuccessfully() {
        // Arrange
        CreateProjectCommand command = new CreateProjectCommand(
                "New Project",
                "Description of New Project",
                "folder-icon",
                "#FF5733"
        );

        when(projectRepository.save(any(Project.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Project result = createProjectUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals("New Project", result.getName());
        assertEquals("Description of New Project", result.getDescription());
        assertEquals("folder-icon", result.getIcon());
        assertEquals("#FF5733", result.getColor());
        assertFalse(result.isArchived());
        assertNotNull(result.getCreatedAt());

        verify(projectRepository, times(1)).save(any(Project.class));
    }
}
