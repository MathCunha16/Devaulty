package com.devaulty.backend.application.impl.link;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.in.link.CreateLinkCommand;
import com.devaulty.backend.application.port.out.persistence.LinkRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Link;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateLinkImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private CreateLinkImpl createLinkUseCase;

    @Test
    void shouldCreateLinkSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateLinkCommand command = new CreateLinkCommand(
                projectId,
                "Test Title",
                "https://example.com",
                "Test Desc"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.save(any(Link.class))).thenAnswer(invocation -> invocation.getArgument(0));

        // Act
        Link result = createLinkUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.getId());
        assertEquals(projectId, result.getProjectId());
        assertEquals("Test Title", result.getTitle());
        assertEquals("https://example.com", result.getUrl());
        assertEquals("Test Desc", result.getDescription());
        assertNotNull(result.getCreatedAt());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).save(any(Link.class));
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        CreateLinkCommand command = new CreateLinkCommand(
                projectId,
                "Test Title",
                "https://example.com",
                "Test Desc"
        );

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createLinkUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, never()).save(any(Link.class));
    }
}
