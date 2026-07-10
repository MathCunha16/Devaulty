package com.devaulty.backend.application.impl.snippet;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.SnippetRepositoryPort;
import com.devaulty.backend.domain.model.Snippet;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
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
class GetAllSnippetsByProjectImplTest {

    @Mock
    private SnippetRepositoryPort snippetRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetAllSnippetsByProjectImpl getAllSnippetsUseCase;

    @Test
    void shouldReturnPageOfSnippets() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        int page = 0;
        int size = 10;
        Snippet snippet = new Snippet(UUID.randomUUID(), projectId, "Title", "Desc", "ls", SnippetLanguage.BASH, SnippetType.COMMAND);
        List<Snippet> list = Collections.singletonList(snippet);
        Page<Snippet> expectedPage = new PageImpl<>(list);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(snippetRepository.findAllByProject(projectId, page, size)).thenReturn(expectedPage);

        // Act
        Page<Snippet> result = getAllSnippetsUseCase.execute(projectId, page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(snippet, result.getContent().getFirst());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, times(1)).findAllByProject(projectId, page, size);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllSnippetsUseCase.execute(projectId, 0, 10);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(snippetRepository, never()).findAllByProject(any(UUID.class), anyInt(), anyInt());
    }
}
