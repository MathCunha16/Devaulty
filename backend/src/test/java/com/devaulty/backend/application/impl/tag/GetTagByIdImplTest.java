package com.devaulty.backend.application.impl.tag;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
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
class GetTagByIdImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetTagByIdImpl getTagByIdUseCase;

    @Test
    void shouldReturnTagSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        Tag tag = new Tag();
        tag.setId(tagId);
        tag.setProjectId(projectId);
        tag.setName("docker");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.of(tag));

        // Act
        Tag result = getTagByIdUseCase.execute(projectId, tagId);

        // Assert
        assertNotNull(result);
        assertEquals(tag, result);
        assertEquals("docker", result.getName());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagByIdUseCase.execute(projectId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).findById(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagByIdUseCase.execute(projectId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagBelongsToAnotherProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        Tag tag = new Tag();
        tag.setId(tagId);
        tag.setProjectId(otherProjectId);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findById(tagId)).thenReturn(Optional.of(tag));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagByIdUseCase.execute(projectId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findById(tagId);
    }
}
