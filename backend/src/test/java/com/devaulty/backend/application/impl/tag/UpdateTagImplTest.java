package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.tag.UpdateTagCommand;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
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
class UpdateTagImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private UpdateTagImpl updateTagUseCase;

    @Test
    void shouldUpdateTagSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UpdateTagCommand command = new UpdateTagCommand(tagId, projectId, "new-name", "#000");

        Tag tag = new Tag();
        tag.setId(tagId);
        tag.setProjectId(projectId);
        tag.setName("old-name");
        tag.setColor("#FFF");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.of(tag));
        when(tagRepository.save(any(Tag.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Tag result = updateTagUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("new-name", result.getName());
        assertEquals("#000", result.getColor());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
        verify(tagRepository, times(1)).save(any(Tag.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UpdateTagCommand command = new UpdateTagCommand(tagId, projectId, "new-name", "#000");

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateTagUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).findById(any());
        verify(tagRepository, never()).save(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UpdateTagCommand command = new UpdateTagCommand(tagId, projectId, "new-name", "#000");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateTagUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
        verify(tagRepository, never()).save(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagBelongsToAnotherProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UpdateTagCommand command = new UpdateTagCommand(tagId, projectId, "new-name", "#000");

        Tag tag = new Tag();
        tag.setId(tagId);
        tag.setProjectId(otherProjectId); // Belongs to otherProjectId

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.of(tag));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateTagUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
        verify(tagRepository, never()).save(any());
    }
}
