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
class SearchTagByNameImplTest {

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private SearchTagByNameImpl searchTagByNameUseCase;

    @Test
    void shouldReturnMatchingTags() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        Tag tag = new Tag();
        tag.setProjectId(projectId);
        tag.setName("java-development");
        List<Tag> expectedList = Collections.singletonList(tag);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.searchByNameAndProjectId(projectId, "java")).thenReturn(expectedList);

        // Act
        List<Tag> result = searchTagByNameUseCase.execute(projectId, "java");

        // Assert
        assertNotNull(result);
        assertEquals(expectedList, result);
        assertEquals(1, result.size());
        assertEquals("java-development", result.getFirst().getName());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).searchByNameAndProjectId(projectId, "java");
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            searchTagByNameUseCase.execute(projectId, "java");
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).searchByNameAndProjectId(any(), any());
    }
}
