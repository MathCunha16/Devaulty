package com.devaulty.backend.application.impl.tag.item;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
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
class RemoveTagFromItemImplTest {

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private RemoveTagFromItemImpl removeTagFromItemUseCase;

    @Test
    void shouldRemoveTagFromItemSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID tagId = UUID.randomUUID();
        UUID itemId = UUID.randomUUID();
        String itemType = "snippet";

        when(projectRepository.existsById(projectId)).thenReturn(true);

        // Act
        removeTagFromItemUseCase.execute(projectId, itemType, itemId, tagId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(itemTagRepository, times(1)).disassembleTagFromItem(tagId, itemType, itemId);
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
            removeTagFromItemUseCase.execute(projectId, itemType, itemId, tagId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(itemTagRepository, never()).disassembleTagFromItem(any(), any(), any());
    }
}
