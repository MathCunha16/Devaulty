package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
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
class DeleteLinkImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private DeleteLinkImpl deleteLinkUseCase;

    @Test
    void shouldDeleteLinkSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link link = new Link(linkId, projectId, "Title", "https://url.com", "Desc");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(link));

        // Act
        deleteLinkUseCase.execute(projectId, linkId);

        // Assert
        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
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
        verify(linkRepository, never()).findById(any(UUID.class));
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteLinkUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link link = new Link(linkId, otherProjectId, "Title", "https://url.com", "Desc");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(link));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteLinkUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
        verify(linkRepository, never()).deleteById(any(UUID.class));
    }
}
