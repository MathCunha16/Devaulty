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

import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetAllTagsByProjectImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetAllTagsByProjectImpl getAllTagsUseCase;

    @Test
    void shouldReturnListofTags() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        Tag tag = new Tag();
        tag.setProjectId(projectId);
        tag.setName("docker");
        List<Tag> expectedList = Collections.singletonList(tag);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.findAllByProject(projectId)).thenReturn(expectedList);

        // Act
        List<Tag> result = getAllTagsUseCase.execute(projectId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedList, result);
        assertEquals(1, result.size());
        assertEquals("docker", result.getFirst().getName());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).findAllByProject(projectId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllTagsUseCase.execute(projectId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).findAllByProject(any());
    }
}
