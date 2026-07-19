package com.devaulty.backend.adapter.in.web.snippet;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.snippet.dto.CreateSnippetRequest;
import com.devaulty.backend.adapter.in.web.snippet.dto.UpdateSnippetRequest;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.snippet.SnippetEntity;
import com.devaulty.backend.adapter.out.persistence.snippet.SpringDataSnippetRepository;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
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

class SnippetControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataSnippetRepository snippetRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;

    @BeforeEach
    void setUpData() {
        snippetRepository.deleteAll();
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
    void createSnippet_shouldReturnCreated() throws Exception {
        CreateSnippetRequest request = new CreateSnippetRequest(
                "My Snippet",
                "Snippet description",
                "echo 'hello world'",
                SnippetLanguage.BASH,
                SnippetType.COMMAND
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/snippets", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").value("My Snippet"))
                .andExpect(jsonPath("$.description").value("Snippet description"))
                .andExpect(jsonPath("$.content").value("echo 'hello world'"))
                .andExpect(jsonPath("$.language").value("BASH"))
                .andExpect(jsonPath("$.snippetType").value("COMMAND"))
                .andExpect(jsonPath("$.createdAt").isNotEmpty())
                .andExpect(jsonPath("$.tags").isEmpty());

        assertEquals(1, snippetRepository.count());
    }

    @Test
    void createSnippet_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Title too short (min 2), content blank
        CreateSnippetRequest request = new CreateSnippetRequest(
                "A",
                "Snippet description",
                "   ",
                SnippetLanguage.BASH,
                SnippetType.COMMAND
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/snippets", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(2)));
    }

    @Test
    void createSnippet_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        CreateSnippetRequest request = new CreateSnippetRequest(
                "Snippet Title",
                "Snippet description",
                "echo 'hello'",
                SnippetLanguage.BASH,
                SnippetType.COMMAND
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/snippets", nonExistentProjectId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404))
                .andExpect(jsonPath("$.message").value("Project not found with identifier " + nonExistentProjectId));
    }

    @Test
    void getAllSnippets_shouldReturnPagedSnippets() throws Exception {
        SnippetEntity s1 = new SnippetEntity(UUID.randomUUID(), savedProject, "Snippet 1", "Desc 1", "content 1", SnippetLanguage.JAVA, SnippetType.CODE);
        s1.setCreatedAt(LocalDateTime.now().minusHours(1));
        SnippetEntity s2 = new SnippetEntity(UUID.randomUUID(), savedProject, "Snippet 2", "Desc 2", "content 2", SnippetLanguage.PYTHON, SnippetType.CODE);
        s2.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(s1);
        snippetRepository.save(s2);

        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets", savedProject.getId())
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].title").value("Snippet 2"))
                .andExpect(jsonPath("$.content[1].title").value("Snippet 1"))
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getAllSnippets_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets", nonExistentProjectId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404));
    }

    @Test
    void getSnippetById_shouldReturnSnippet() throws Exception {
        UUID id = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(id, savedProject, "Get Snippet", "Desc", "print(1)", SnippetLanguage.PYTHON, SnippetType.CODE);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), id))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(id.toString()))
                .andExpect(jsonPath("$.title").value("Get Snippet"));
    }

    @Test
    void getSnippetById_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets/{snippetId}", nonExistentProjectId, snippetId))
                .andExpect(status().isNotFound());
    }

    @Test
    void getSnippetById_shouldReturnNotFound_whenSnippetDoesNotExist() throws Exception {
        UUID nonExistentSnippetId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), nonExistentSnippetId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Snippet not found with identifier " + nonExistentSnippetId));
    }

    @Test
    void getSnippetById_shouldReturnNotFound_whenSnippetDoesNotBelongToProject() throws Exception {
        // Create another project
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        // Create snippet for other project
        UUID snippetId = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(snippetId, otherProject, "Other Snippet", "Desc", "ls", SnippetLanguage.BASH, SnippetType.COMMAND);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        // Perform request requesting other project's snippet under savedProject's path
        mockMvc.perform(get("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), snippetId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Snippet not found with identifier " + snippetId));
    }

    @Test
    void updateSnippet_shouldUpdateAndReturnSnippet() throws Exception {
        UUID id = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(id, savedProject, "Old Title", "Old Desc", "old content", SnippetLanguage.JAVA, SnippetType.CODE);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        UpdateSnippetRequest request = new UpdateSnippetRequest(
                "New Title",
                "New Desc",
                "new content",
                SnippetLanguage.KOTLIN,
                SnippetType.CODE
        );

        mockMvc.perform(patch("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.title").value("New Title"))
                .andExpect(jsonPath("$.description").value("New Desc"))
                .andExpect(jsonPath("$.content").value("new content"))
                .andExpect(jsonPath("$.language").value("KOTLIN"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        SnippetEntity updated = snippetRepository.findById(id).orElseThrow();
        assertEquals("New Title", updated.getTitle());
    }

    @Test
    void updateSnippet_shouldReturnBadRequest_whenValidationFails() throws Exception {
        UUID id = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(id, savedProject, "Old Title", "Old Desc", "old content", SnippetLanguage.JAVA, SnippetType.CODE);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        // Title too short
        UpdateSnippetRequest request = new UpdateSnippetRequest("A", null, null, null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"));
    }

    @Test
    void updateSnippet_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();
        UpdateSnippetRequest request = new UpdateSnippetRequest("New Title", null, null, null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/snippets/{snippetId}", nonExistentProjectId, snippetId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void updateSnippet_shouldReturnNotFound_whenSnippetDoesNotExist() throws Exception {
        UUID nonExistentSnippetId = UUID.randomUUID();
        UpdateSnippetRequest request = new UpdateSnippetRequest("New Title", null, null, null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), nonExistentSnippetId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void updateSnippet_shouldReturnNotFound_whenSnippetDoesNotBelongToProject() throws Exception {
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        UUID snippetId = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(snippetId, otherProject, "Other Snippet", "Desc", "ls", SnippetLanguage.BASH, SnippetType.COMMAND);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        UpdateSnippetRequest request = new UpdateSnippetRequest("New Title", null, null, null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), snippetId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteSnippet_shouldDeleteAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        SnippetEntity snippet = new SnippetEntity(id, savedProject, "To Delete", "Desc", "rm", SnippetLanguage.BASH, SnippetType.COMMAND);
        snippet.setCreatedAt(LocalDateTime.now());
        snippetRepository.save(snippet);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/snippets/{snippetId}", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        assertFalse(snippetRepository.existsById(id));
    }

    @Test
    void deleteSnippet_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID snippetId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/snippets/{snippetId}", nonExistentProjectId, snippetId))
                .andExpect(status().isNotFound());
    }
}
