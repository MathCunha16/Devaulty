package com.devaulty.backend.application.impl.project;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DeleteProjectImplTest {

    @Mock
    private ProjectRepositoryPort projectRepositoryPort;

    @InjectMocks
    private DeleteProjectImpl deleteProjectUseCase;

    @Test
    void shouldDeleteProjectSuccessfully() {
        // Arrange
        UUID id = UUID.randomUUID();
        when(projectRepositoryPort.existsById(id)).thenReturn(true);

        // Act
        deleteProjectUseCase.execute(id);

        // Assert
        verify(projectRepositoryPort, times(1)).existsById(id);
        verify(projectRepositoryPort, times(1)).deleteById(id);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID id = UUID.randomUUID();
        when(projectRepositoryPort.existsById(id)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteProjectUseCase.execute(id);
        });

        verify(projectRepositoryPort, times(1)).existsById(id);
        verify(projectRepositoryPort, never()).deleteById(any(UUID.class));
    }
}
