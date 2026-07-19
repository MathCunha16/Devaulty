package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AssociateTagToItemImplTest {

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @Mock
    private TagRepositoryPort tagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private AssociateTagToItemImpl associateTagToItemUseCase;

    @Test
    void shouldAssociateTagToItemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(true);

        // Act
        associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(itemTagRepository, times(1)).associateTagToItem(tagId, itemType, itemId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, never()).existsByIdAndProjectId(any(), any());
        verify(itemTagRepository, never()).associateTagToItem(any(), any(), any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenTagDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(itemTagRepository, never()).associateTagToItem(any(), any(), any());
    }
}
