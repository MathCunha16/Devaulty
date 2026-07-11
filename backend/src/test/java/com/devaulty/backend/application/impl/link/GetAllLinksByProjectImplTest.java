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
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;

import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetAllLinksByProjectImplTest {

    @Mock
    private LinkRepositoryPort linkRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetAllLinksByProjectImpl getAllLinksUseCase;

    @Test
    void shouldReturnPageOfLinks() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        int page = 0;
        int size = 10;
        Link link = new Link(UUID.randomUUID(), projectId, "Title", "https://url.com", "Desc");
        List<Link> list = Collections.singletonList(link);
        Page<Link> expectedPage = new PageImpl<>(list);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(linkRepository.findAllByProject(projectId, page, size)).thenReturn(expectedPage);

        // Act
        Page<Link> result = getAllLinksUseCase.execute(projectId, page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(link, result.getContent().getFirst());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, times(1)).findAllByProject(projectId, page, size);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllLinksUseCase.execute(projectId, 0, 10);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(linkRepository, never()).findAllByProject(any(UUID.class), anyInt(), anyInt());
    }
}
