package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetTagsForItemsImplTest {

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetTagsForItemsImpl getTagsForItemsUseCase;

    @Test
    void shouldReturnTagsForItems() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        List<UUID> itemIds = Collections.singletonList(itemId);
        String itemType = "snippet";

        Tag tag = new Tag();
        tag.setProjectId(projectId);
        tag.setName("docker");
        List<Tag> tagsList = Collections.singletonList(tag);

        Map<UUID, List<Tag>> expectedMap = new HashMap<>();
        expectedMap.put(itemId, tagsList);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(itemTagRepository.findTagsForItems(itemType, itemIds)).thenReturn(expectedMap);

        // Act
        Map<UUID, List<Tag>> result = getTagsForItemsUseCase.execute(itemType, projectId, itemIds);

        // Assert
        assertNotNull(result);
        assertEquals(expectedMap, result);
        assertTrue(result.containsKey(itemId));
        assertEquals(1, result.get(itemId).size());
        assertEquals("docker", result.get(itemId).getFirst().getName());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(itemTagRepository, times(1)).findTagsForItems(itemType, itemIds);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        List<UUID> itemIds = Collections.singletonList(itemId);
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getTagsForItemsUseCase.execute(itemType, projectId, itemIds);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(itemTagRepository, never()).findTagsForItems(any(), any());
    }
}
