package com.devaulty.backend.adapter.in.web.project;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.project.dto.CreateProjectRequest;
import com.devaulty.backend.adapter.in.web.project.dto.UpdateProjectRequest;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
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

class ProjectControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataProjectRepository projectRepository;

    @BeforeEach
    void cleanDatabase() {
        projectRepository.deleteAll();
    }

    @Test
    void createProject_shouldReturnCreated() throws Exception {
        CreateProjectRequest request = new CreateProjectRequest(
                "My Project",
                "Project description",
                "folder",
                "#FF5733"
        );

        mockMvc.perform(post("/api/v1/projects")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.name").value("My Project"))
                .andExpect(jsonPath("$.description").value("Project description"))
                .andExpect(jsonPath("$.icon").value("folder"))
                .andExpect(jsonPath("$.color").value("#FF5733"))
                .andExpect(jsonPath("$.archived").value(false))
                .andExpect(jsonPath("$.createdAt").isNotEmpty());

        assertEquals(1, projectRepository.count());
    }

    @Test
    void createProject_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Name too short (min 3)
        CreateProjectRequest request = new CreateProjectRequest(
                "A",
                "Project description",
                "folder",
                "#invalidColor"
        );

        mockMvc.perform(post("/api/v1/projects")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(2)));
    }

    @Test
    void getAllProjects_shouldReturnPagedProjects() throws Exception {
        ProjectEntity p1 = new ProjectEntity(UUID.randomUUID(), "Project 1", "Desc 1", "icon1", "#000", false);
        p1.setCreatedAt(LocalDateTime.now().minusHours(1)); // Older
        ProjectEntity p2 = new ProjectEntity(UUID.randomUUID(), "Project 2", "Desc 2", "icon2", "#fff", false);
        p2.setCreatedAt(LocalDateTime.now()); // Newer
        projectRepository.save(p1);
        projectRepository.save(p2);

        mockMvc.perform(get("/api/v1/projects")
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].name").value("Project 2")) // Sorted DESC by createdAt
                .andExpect(jsonPath("$.content[1].name").value("Project 1"))
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getProjectById_shouldReturnProject() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "Unique Project", "Description", "star", "#123", false);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(get("/api/v1/projects/{id}", id))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(id.toString()))
                .andExpect(jsonPath("$.name").value("Unique Project"));
    }

    @Test
    void getProjectById_shouldReturnNotFound_whenDoesNotExist() throws Exception {
        UUID randomId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{id}", randomId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404))
                .andExpect(jsonPath("$.message").value("Project not found with identifier " + randomId));
    }

    @Test
    void updateProject_shouldUpdateAndReturnProject() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "Old Name", "Old Desc", "old-icon", "#000", false);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        UpdateProjectRequest request = new UpdateProjectRequest(
                "New Name",
                "New Desc",
                "new-icon",
                "#FFF"
        );

        mockMvc.perform(patch("/api/v1/projects/{id}", id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.name").value("New Name"))
                .andExpect(jsonPath("$.description").value("New Desc"))
                .andExpect(jsonPath("$.icon").value("new-icon"))
                .andExpect(jsonPath("$.color").value("#FFF"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        ProjectEntity updated = projectRepository.findById(id).orElseThrow();
        assertEquals("New Name", updated.getName());
    }

    @Test
    void updateProject_shouldReturnNotFound_whenDoesNotExist() throws Exception {
        UUID randomId = UUID.randomUUID();
        UpdateProjectRequest request = new UpdateProjectRequest("Name", "Desc", "icon", "#fff");

        mockMvc.perform(patch("/api/v1/projects/{id}", randomId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void archiveProject_shouldArchiveAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "To Archive", "Desc", "icon", "#000", false);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(patch("/api/v1/projects/{id}/archive", id))
                .andExpect(status().isNoContent());

        ProjectEntity archived = projectRepository.findById(id).orElseThrow();
        assertTrue(archived.isArchived());
        assertNotNull(archived.getUpdatedAt());
    }

    @Test
    void archiveProject_shouldReturnBadRequest_whenAlreadyArchived() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "Archived", "Desc", "icon", "#000", true);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(patch("/api/v1/projects/{id}/archive", id))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Project already archived"));
    }

    @Test
    void unarchiveProject_shouldUnarchiveAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "To Unarchive", "Desc", "icon", "#000", true);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(patch("/api/v1/projects/{id}/unarchive", id))
                .andExpect(status().isNoContent());

        ProjectEntity unarchived = projectRepository.findById(id).orElseThrow();
        assertFalse(unarchived.isArchived());
        assertNotNull(unarchived.getUpdatedAt());
    }

    @Test
    void unarchiveProject_shouldReturnBadRequest_whenAlreadyActive() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "Active", "Desc", "icon", "#000", false);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(patch("/api/v1/projects/{id}/unarchive", id))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Project already unarchived"));
    }

    @Test
    void deleteProject_shouldDeleteAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        ProjectEntity project = new ProjectEntity(id, "To Delete", "Desc", "icon", "#000", false);
        project.setCreatedAt(LocalDateTime.now());
        projectRepository.save(project);

        mockMvc.perform(delete("/api/v1/projects/{id}", id))
                .andExpect(status().isNoContent());

        assertFalse(projectRepository.existsById(id));
    }

    @Test
    void deleteProject_shouldReturnNotFound_whenDoesNotExist() throws Exception {
        UUID randomId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{id}", randomId))
                .andExpect(status().isNotFound());
    }
}
