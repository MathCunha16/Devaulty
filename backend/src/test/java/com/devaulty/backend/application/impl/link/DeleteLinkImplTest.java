package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DeleteLinkImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @InjectMocks
    private DeleteLinkImpl deleteLinkUseCase;

    @Test
    void shouldDeleteLinkSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.existsByIdAndProjectId(linkId, projectId)).thenReturn(true);

        // Act
        deleteLinkUseCase.execute(projectId, linkId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).existsByIdAndProjectId(linkId, projectId);
        verify(itemTagRepository, times(1)).removeAllTagsFromItem("link", linkId);
        verify(linkRepository, times(1)).deleteById(linkId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteLinkUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, never()).existsByIdAndProjectId(any(UUID.class), any(UUID.class));
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.existsByIdAndProjectId(linkId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteLinkUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).existsByIdAndProjectId(linkId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        // existsByIdAndProjectId returns false because it belongs to another project
        when(linkRepository.existsByIdAndProjectId(linkId, projectId)).thenReturn(false);

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            deleteLinkUseCase.execute(projectId, linkId);
        });
        assertEquals("Link not found with identifier " + linkId, exception.getMessage());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).existsByIdAndProjectId(linkId, projectId);
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }
}
