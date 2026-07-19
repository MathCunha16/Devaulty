package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectScopedRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.TagRepositoryPort;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Collections;
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

    @Mock
    private ProjectScopedRepositoryPort snippetRepository;

    private AssociateTagToItemImpl associateTagToItemUseCase;

    @BeforeEach
    void setUp() {
        associateTagToItemUseCase = new AssociateTagToItemImpl(
                itemTagRepository,
                tagRepository,
                projectRepository,
                Collections.singletonList(snippetRepository)
        );
    }

    @Test
    void shouldAssociateTagToItemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(true);
        when(snippetRepository.getSupportedType()).thenReturn(itemType);
        when(snippetRepository.existsByIdAndProjectId(itemId, projectId)).thenReturn(true);

        // Act
        associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(itemId, projectId);
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
        verify(snippetRepository, never()).existsByIdAndProjectId(any(), any());
        verify(itemTagRepository, never()).associateTagToItem(any(), any(), any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenItemDoesNotExistOrDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(true);
        when(snippetRepository.getSupportedType()).thenReturn(itemType);
        when(snippetRepository.existsByIdAndProjectId(itemId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(snippetRepository, times(1)).existsByIdAndProjectId(itemId, projectId);
        verify(itemTagRepository, never()).associateTagToItem(any(), any(), any());
    }

    @Test
    void shouldThrowIllegalArgumentExceptionWhenItemTypeIsUnsupported() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "unsupported-type";

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(tagRepository.existsByIdAndProjectId(tagId, projectId)).thenReturn(true);
        when(snippetRepository.getSupportedType()).thenReturn("snippet");

        // Act & Assert
        assertThrows(IllegalArgumentException.class, () -> {
            associateTagToItemUseCase.execute(projectId, itemType, itemId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(tagRepository, times(1)).existsByIdAndProjectId(tagId, projectId);
        verify(itemTagRepository, never()).associateTagToItem(any(), any(), any());
    }
}
