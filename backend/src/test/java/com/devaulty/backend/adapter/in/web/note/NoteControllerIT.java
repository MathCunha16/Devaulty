package com.devaulty.backend.adapter.in.web.note;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.note.dto.*;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.note.NoteEntity;
import com.devaulty.backend.adapter.out.persistence.note.SpringDataNoteRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;

import java.time.LocalDateTime;
import java.util.UUID;

import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

class NoteControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataNoteRepository noteRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;

    @BeforeEach
    void setUpData() {
        noteRepository.deleteAll();
        projectRepository.deleteAll();

        ProjectEntity project = new ProjectEntity(
                UUID.randomUUID(),
                "Integration Project",
                "Description",
                "folder",
                "#FFF",
                false
        );
        project.setCreatedAt(LocalDateTime.now());
        savedProject = projectRepository.save(project);
    }

    @Test
    void createNote_shouldReturnCreated() throws Exception {
        CreateNoteRequest request = new CreateNoteRequest(
                "My Note",
                "Note content here"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/notes", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").value("My Note"))
                .andExpect(jsonPath("$.content").value("Note content here"))
                .andExpect(jsonPath("$.archived").value(false))
                .andExpect(jsonPath("$.createdAt").isNotEmpty())
                .andExpect(jsonPath("$.tags").isEmpty());

        assertEquals(1, noteRepository.count());
    }

    @Test
    void createNote_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Title blank (min 2)
        CreateNoteRequest request = new CreateNoteRequest(
                "   ",
                "Content"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/notes", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(1)));
    }

    @Test
    void createNote_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        CreateNoteRequest request = new CreateNoteRequest(
                "Note Title",
                "content"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/notes", nonExistentProjectId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404))
                .andExpect(jsonPath("$.message").value("Project not found with identifier " + nonExistentProjectId));
    }

    @Test
    void getAllNotes_shouldReturnPagedNotes() throws Exception {
        NoteEntity n1 = new NoteEntity(UUID.randomUUID(), savedProject, "Note 1", "Content 1", false);
        n1.setCreatedAt(LocalDateTime.now().minusHours(1));
        NoteEntity n2 = new NoteEntity(UUID.randomUUID(), savedProject, "Note 2", "Content 2", false);
        n2.setCreatedAt(LocalDateTime.now());
        noteRepository.save(n1);
        noteRepository.save(n2);

        mockMvc.perform(get("/api/v1/projects/{projectId}/notes", savedProject.getId())
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].title").value("Note 2"))
                .andExpect(jsonPath("$.content[0].tags").isEmpty())
                .andExpect(jsonPath("$.content[1].title").value("Note 1"))
                .andExpect(jsonPath("$.content[1].tags").isEmpty())
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getAllNotes_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/notes", nonExistentProjectId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404));
    }

    @Test
    void getNoteById_shouldReturnNote() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Get Note", "Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(get("/api/v1/projects/{projectId}/notes/{noteID}", savedProject.getId(), id))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(id.toString()))
                .andExpect(jsonPath("$.title").value("Get Note"))
                .andExpect(jsonPath("$.tags").isEmpty());
    }

    @Test
    void getNoteById_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/notes/{noteID}", nonExistentProjectId, noteId))
                .andExpect(status().isNotFound());
    }

    @Test
    void getNoteById_shouldReturnNotFound_whenNoteDoesNotExist() throws Exception {
        UUID nonExistentNoteId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/notes/{noteID}", savedProject.getId(), nonExistentNoteId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Note not found with identifier " + nonExistentNoteId));
    }

    @Test
    void getNoteById_shouldReturnNotFound_whenNoteDoesNotBelongToProject() throws Exception {
        // Create another project
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        // Create note for other project
        UUID noteId = UUID.randomUUID();
        NoteEntity note = new NoteEntity(noteId, otherProject, "Other Note", "Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        // Request other project's note under savedProject's path
        mockMvc.perform(get("/api/v1/projects/{projectId}/notes/{noteID}", savedProject.getId(), noteId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Note not found with identifier " + noteId));
    }

    @Test
    void updateNote_shouldUpdateAndReturnNote() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Old Title", "Old Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        UpdateNoteRequest request = new UpdateNoteRequest(
                "New Title",
                "New Content"
        );

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.title").value("New Title"))
                .andExpect(jsonPath("$.content").value("New Content"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        NoteEntity updated = noteRepository.findById(id).orElseThrow();
        assertEquals("New Title", updated.getTitle());
    }

    @Test
    void updateNote_shouldReturnBadRequest_whenValidationFails() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Old Title", "Old Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        // Blank title
        UpdateNoteRequest request = new UpdateNoteRequest("   ", null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"));
    }

    @Test
    void archiveNote_shouldArchiveAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "To Archive", "Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}/archive", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        NoteEntity archived = noteRepository.findById(id).orElseThrow();
        assertTrue(archived.isArchived());
    }

    @Test
    void archiveNote_shouldReturnBadRequest_whenAlreadyArchived() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Archived", "Content", true);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}/archive", savedProject.getId(), id))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Note already archived"));
    }

    @Test
    void unarchiveNote_shouldUnarchiveAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Archived", "Content", true);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}/unarchive", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        NoteEntity unarchived = noteRepository.findById(id).orElseThrow();
        assertFalse(unarchived.isArchived());
    }

    @Test
    void unarchiveNote_shouldReturnBadRequest_whenNotArchived() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "Active", "Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/notes/{noteId}/unarchive", savedProject.getId(), id))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Note not archived"));
    }

    @Test
    void deleteNote_shouldDeleteAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        NoteEntity note = new NoteEntity(id, savedProject, "To Delete", "Content", false);
        note.setCreatedAt(LocalDateTime.now());
        noteRepository.save(note);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/notes/{noteId}", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        assertFalse(noteRepository.existsById(id));
    }

    @Test
    void deleteNote_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID noteId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/notes/{noteId}", nonExistentProjectId, noteId))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteNote_shouldReturnNotFound_whenNoteDoesNotExist() throws Exception {
        UUID nonExistentNoteId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/notes/{noteId}", savedProject.getId(), nonExistentNoteId))
                .andExpect(status().isNotFound());
    }
}
