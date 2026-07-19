package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceAlreadyExistsException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.CreateTagCommand;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateTagImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateTagImpl createTagUseCase;

    @Test
    void shouldCreateTagSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateTagCommand command = new CreateTagCommand(projectId, "java", "#3A3A3A");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByNameAndProjectId(projectId, "java")).thenReturn(false);

        Tag expectedTag = new Tag();
        expectedTag.setProjectId(projectId);
        expectedTag.setName("java");
        expectedTag.setColor("#3A3A3A");
        when(tagRepository.save(any(Tag.class))).thenReturn(expectedTag);

        // Act
        Tag result = createTagUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("java", result.getName());
        assertEquals("#3A3A3A", result.getColor());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByNameAndProjectId(projectId, "java");
        verify(tagRepository, times(1)).save(any(Tag.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateTagCommand command = new CreateTagCommand(projectId, "java", "#3A3A3A");

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createTagUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).existsByNameAndProjectId(any(), any());
        verify(tagRepository, never()).save(any());
    }

    @Test
    void shouldThrowResourceAlreadyExistsExceptionWhenTagAlreadyExistsInProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateTagCommand command = new CreateTagCommand(projectId, "java", "#3A3A3A");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByNameAndProjectId(projectId, "java")).thenReturn(true);

        // Act & Assert
        assertThrows(ResourceAlreadyExistsException.class, () -> {
            createTagUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByNameAndProjectId(projectId, "java");
        verify(tagRepository, never()).save(any());
    }
}
