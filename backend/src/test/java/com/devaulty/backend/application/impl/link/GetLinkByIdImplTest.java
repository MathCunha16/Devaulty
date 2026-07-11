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
class GetLinkByIdImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetLinkByIdImpl getLinkByIdUseCase;

    @Test
    void shouldReturnLinkWhenFoundAndBelongsToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link expectedLink = new Link(linkId, projectId, "Title", "https://url.com", "Desc");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(expectedLink));

        // Act
        Link result = getLinkByIdUseCase.execute(projectId, linkId);

        // Assert
        assertNotNull(result);
        assertEquals(expectedLink, result);

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getLinkByIdUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, never()).findById(any(UUID.class));
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
            getLinkByIdUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link linkOfOtherProject = new Link(linkId, otherProjectId, "Title", "https://url.com", "Desc");

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(linkOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getLinkByIdUseCase.execute(projectId, linkId);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
    }
}
