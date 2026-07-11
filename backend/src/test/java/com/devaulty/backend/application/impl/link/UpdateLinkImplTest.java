package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.UpdateLinkCommand;
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
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UpdateLinkImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private UpdateLinkImpl updateLinkUseCase;

    @Test
    void shouldUpdateFieldsWhenValidCommand() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link existingLink = new Link(linkId, projectId, "Old Title", "https://old.com", "Old Desc");
        UpdateLinkCommand command = new UpdateLinkCommand(
                linkId,
                projectId,
                "New Title",
                "https://new.com",
                "New Desc"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(existingLink));
        when(linkRepository.save(any(Link.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Link result = updateLinkUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("New Title", result.getTitle());
        assertEquals("https://new.com", result.getUrl());
        assertEquals("New Desc", result.getDescription());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
        verify(linkRepository, times(1)).save(existingLink);
    }

    @Test
    void shouldNotUpdateFieldsWhenNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        Link existingLink = new Link(linkId, projectId, "Old Title", "https://old.com", "Old Desc");
        UpdateLinkCommand command = new UpdateLinkCommand(
                linkId,
                projectId,
                null,
                null,
                null
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(existingLink));
        when(linkRepository.save(any(Link.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Link result = updateLinkUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("Old Title", result.getTitle());
        assertEquals("https://old.com", result.getUrl());
        assertEquals("Old Desc", result.getDescription());
        assertNotNull(result.getUpdatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
        verify(linkRepository, times(1)).save(existingLink);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        UpdateLinkCommand command = new UpdateLinkCommand(linkId, projectId, "Title", null, null);

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateLinkUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, never()).findById(any(UUID.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenLinkDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        UpdateLinkCommand command = new UpdateLinkCommand(linkId, projectId, "Title", null, null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateLinkUseCase.execute(command);
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
        UpdateLinkCommand command = new UpdateLinkCommand(linkId, projectId, "Title", null, null);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findById(linkId)).thenReturn(Optional.of(linkOfOtherProject));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateLinkUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findById(linkId);
        verify(linkRepository, never()).save(any(Link.class));
    }
}
