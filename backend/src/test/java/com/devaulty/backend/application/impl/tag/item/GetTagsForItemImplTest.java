package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetTagsForItemImplTest {

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private ProjectScopedRepositoryPort snippetRepository;

    private GetTagsForItemImpl getTagsForItemUseCase;

    @BeforeEach
    void setUp() {
        getTagsForItemUseCase = new GetTagsForItemImpl(
                itemTagRepository,
                projectRepository,
                Collections.singletonList(snippetRepository)
        );
    }

    @Test
    void shouldReturnTagsForItem_whenProjectAndItemExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        Tag tag = new Tag();
        tag.setProjectId(projectId);
        tag.setName("docker");
        List<Tag> expectedList = Collections.singletonList(tag);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.getSupportedType()).thenReturn(itemType);
        when(snippetRepository.existsByIdAndProjectId(itemId, projectId)).thenReturn(true);
        when(itemTagRepository.findTagsForItem(itemType, projectId, itemId)).thenReturn(expectedList);

        // Act
        List<Tag> result = getTagsForItemUseCase.execute(itemType, projectId, itemId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedList, result);
        assertEquals(1, result.size());
        assertEquals("docker", result.getFirst().getName());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(itemId, projectId);
        verify(itemTagRepository, times(1)).findTagsForItem(itemType, projectId, itemId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagsForItemUseCase.execute(itemType, projectId, itemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(itemTagRepository, never()).findTagsForItem(any(), any(), any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenItemDoesNotExistOrDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.getSupportedType()).thenReturn(itemType);
        when(snippetRepository.existsByIdAndProjectId(itemId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagsForItemUseCase.execute(itemType, projectId, itemId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(itemId, projectId);
        verify(itemTagRepository, never()).findTagsForItem(any(), any(), any());
    }
}
