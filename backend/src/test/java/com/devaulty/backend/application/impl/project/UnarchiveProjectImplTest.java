package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.BusinessRuleException;
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
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UnarchiveProjectImplTest {

    @Mock
    private ProjectRepositoryPort projectRepositoryPort;

    @InjectMocks
    private UnarchiveProjectImpl unarchiveProjectUseCase;

    @Test
    void shouldUnarchiveProjectSuccessfully() {
        // Arrange
        UUID id = UUID.randomUUID();
        Project project = new Project(id, "Test Project", "Desc", "icon", "#fff", true);

        when(projectRepositoryPort.findById(id)).thenReturn(Optional.of(project));
        when(projectRepositoryPort.save(any(Project.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        unarchiveProjectUseCase.execute(id);

        // Assert
        assertFalse(project.isArchived());
        assertNotNull(project.getUpdatedAt());

        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, times(1)).save(project);
    }

    @Test
    void shouldThrowBusinessRuleExceptionWhenAlreadyUnarchived() {
        // Arrange
        UUID id = UUID.randomUUID();
        Project project = new Project(id, "Test Project", "Desc", "icon", "#fff", false);

        when(projectRepositoryPort.findById(id)).thenReturn(Optional.of(project));

        // Act & Assert
        BusinessRuleException exception = assertThrows(BusinessRuleException.class, () -> {
            unarchiveProjectUseCase.execute(id);
        });

        assertEquals("Project already unarchived", exception.getMessage());
        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, never()).save(any(Project.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID id = UUID.randomUUID();
        when(projectRepositoryPort.findById(id)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            unarchiveProjectUseCase.execute(id);
        });

        verify(projectRepositoryPort, times(1)).findById(id);
        verify(projectRepositoryPort, never()).save(any(Project.class));
    }
}
