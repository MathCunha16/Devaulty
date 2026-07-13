package com.devaulty.backend.application.impl.note;

import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.port.out.persistence.NoteRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.domain.model.Note;
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
class GetAllNotesByProjectImplTest {

    @Mock
    private NoteRepositoryPort noteRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @InjectMocks
    private GetAllNotesByProjectImpl getAllNotesUseCase;

    @Test
    void shouldReturnPageOfNotes() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        int page = 0;
        int size = 10;
        Note note = new Note(UUID.randomUUID(), projectId, "Title", "Content", false);
        List<Note> list = Collections.singletonList(note);
        Page<Note> expectedPage = new PageImpl<>(list);

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(noteRepository.findAllByProject(projectId, page, size)).thenReturn(expectedPage);

        // Act
        Page<Note> result = getAllNotesUseCase.execute(projectId, page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(note, result.getContent().getFirst());

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, times(1)).findAllByProject(projectId, page, size);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllNotesUseCase.execute(projectId, 0, 10);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(noteRepository, never()).findAllByProject(any(UUID.class), anyInt(), anyInt());
    }
}
