package com.devaulty.backend.adapter.in.web.link;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.link.dto.CreateLinkRequest;
import com.devaulty.backend.adapter.in.web.link.dto.UpdateLinkRequest;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.link.LinkEntity;
import com.devaulty.backend.adapter.out.persistence.link.SpringDataLinkRepository;
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

class LinkControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataLinkRepository linkRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;

    @BeforeEach
    void setUpData() {
        linkRepository.deleteAll();
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
    void createLink_shouldReturnCreated() throws Exception {
        CreateLinkRequest request = new CreateLinkRequest(
                "My Link",
                "https://example.com",
                "Link description"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/links", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").value("My Link"))
                .andExpect(jsonPath("$.url").value("https://example.com"))
                .andExpect(jsonPath("$.description").value("Link description"))
                .andExpect(jsonPath("$.createdAt").isNotEmpty());

        assertEquals(1, linkRepository.count());
    }

    @Test
    void createLink_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Title blank (min 2), URL blank
        CreateLinkRequest request = new CreateLinkRequest(
                "   ",
                "   ",
                "Link description"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/links", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(2)));
    }

    @Test
    void createLink_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        CreateLinkRequest request = new CreateLinkRequest(
                "Link Title",
                "https://google.com",
                "description"
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/links", nonExistentProjectId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404))
                .andExpect(jsonPath("$.message").value("Project not found with identifier " + nonExistentProjectId));
    }

    @Test
    void getAllLinks_shouldReturnPagedLinks() throws Exception {
        LinkEntity l1 = new LinkEntity(UUID.randomUUID(), savedProject, "Link 1", "https://url1.com", "Desc 1");
        l1.setCreatedAt(LocalDateTime.now().minusHours(1));
        LinkEntity l2 = new LinkEntity(UUID.randomUUID(), savedProject, "Link 2", "https://url2.com", "Desc 2");
        l2.setCreatedAt(LocalDateTime.now());
        linkRepository.save(l1);
        linkRepository.save(l2);

        mockMvc.perform(get("/api/v1/projects/{projectId}/links", savedProject.getId())
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].title").value("Link 2")) // Sorted DESC by createdAt
                .andExpect(jsonPath("$.content[1].title").value("Link 1"))
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getAllLinks_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/links", nonExistentProjectId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404));
    }

    @Test
    void getLinkById_shouldReturnLink() throws Exception {
        UUID id = UUID.randomUUID();
        LinkEntity link = new LinkEntity(id, savedProject, "Get Link", "https://url.com", "Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        mockMvc.perform(get("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), id))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(id.toString()))
                .andExpect(jsonPath("$.title").value("Get Link"));
    }

    @Test
    void getLinkById_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/links/{linkId}", nonExistentProjectId, linkId))
                .andExpect(status().isNotFound());
    }

    @Test
    void getLinkById_shouldReturnNotFound_whenLinkDoesNotExist() throws Exception {
        UUID nonExistentLinkId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), nonExistentLinkId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Link not found with identifier " + nonExistentLinkId));
    }

    @Test
    void getLinkById_shouldReturnNotFound_whenLinkDoesNotBelongToProject() throws Exception {
        // Create another project
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        // Create link for other project
        UUID linkId = UUID.randomUUID();
        LinkEntity link = new LinkEntity(linkId, otherProject, "Other Link", "https://other.com", "Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        // Request other project's link under savedProject's path
        mockMvc.perform(get("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), linkId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Link not found with identifier " + linkId));
    }

    @Test
    void updateLink_shouldUpdateAndReturnLink() throws Exception {
        UUID id = UUID.randomUUID();
        LinkEntity link = new LinkEntity(id, savedProject, "Old Title", "https://old.com", "Old Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        UpdateLinkRequest request = new UpdateLinkRequest(
                "New Title",
                "https://new.com",
                "New Desc"
        );

        mockMvc.perform(patch("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.title").value("New Title"))
                .andExpect(jsonPath("$.url").value("https://new.com"))
                .andExpect(jsonPath("$.description").value("New Desc"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        LinkEntity updated = linkRepository.findById(id).orElseThrow();
        assertEquals("New Title", updated.getTitle());
    }

    @Test
    void updateLink_shouldReturnBadRequest_whenValidationFails() throws Exception {
        UUID id = UUID.randomUUID();
        LinkEntity link = new LinkEntity(id, savedProject, "Old Title", "https://old.com", "Old Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        // Blank/spaces title
        UpdateLinkRequest request = new UpdateLinkRequest("   ", null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"));
    }

    @Test
    void updateLink_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();
        UpdateLinkRequest request = new UpdateLinkRequest("New Title", null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/links/{linkId}", nonExistentProjectId, linkId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void updateLink_shouldReturnNotFound_whenLinkDoesNotExist() throws Exception {
        UUID nonExistentLinkId = UUID.randomUUID();
        UpdateLinkRequest request = new UpdateLinkRequest("New Title", null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), nonExistentLinkId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void updateLink_shouldReturnNotFound_whenLinkDoesNotBelongToProject() throws Exception {
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        UUID linkId = UUID.randomUUID();
        LinkEntity link = new LinkEntity(linkId, otherProject, "Other Link", "https://other.com", "Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        UpdateLinkRequest request = new UpdateLinkRequest("New Title", null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), linkId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteLink_shouldDeleteAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        LinkEntity link = new LinkEntity(id, savedProject, "To Delete", "https://delete.com", "Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        assertFalse(linkRepository.existsById(id));
    }

    @Test
    void deleteLink_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID linkId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/links/{linkId}", nonExistentProjectId, linkId))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteLink_shouldReturnNotFound_whenLinkDoesNotExist() throws Exception {
        UUID nonExistentLinkId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), nonExistentLinkId))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteLink_shouldReturnNotFound_whenLinkDoesNotBelongToProject() throws Exception {
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "`#000`", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        UUID linkId = UUID.randomUUID();
        LinkEntity link = new LinkEntity(linkId, otherProject, "Other Link", "https://other.com", "Desc");
        link.setCreatedAt(LocalDateTime.now());
        linkRepository.save(link);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/links/{linkId}", savedProject.getId(), linkId))
                .andExpect(status().isNotFound());
    }
}
