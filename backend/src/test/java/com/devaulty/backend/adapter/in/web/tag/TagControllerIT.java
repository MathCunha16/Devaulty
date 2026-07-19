package com.devaulty.backend.adapter.in.web.tag;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.tag.dto.CreateTagRequest;
import com.devaulty.backend.adapter.in.web.tag.dto.UpdateTagRequest;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.tag.TagEntity;
import com.devaulty.backend.adapter.out.persistence.tag.SpringDataTagRepository;
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

class TagControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataTagRepository tagRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;

    @BeforeEach
    void setUpData() {
        tagRepository.deleteAll();
        projectRepository.deleteAll();

        ProjectEntity project = new ProjectEntity(
                UUID.randomUUID(),
                "Tag Integration Project",
                "Description",
                "folder",
                "#FFF",
                false
        );
        project.setCreatedAt(LocalDateTime.now());
        savedProject = projectRepository.save(project);
    }

    private TagEntity createAndSaveTag(String name, String color) {
        TagEntity tag = new TagEntity(UUID.randomUUID(), savedProject, name, color);
        tag.setCreatedAt(LocalDateTime.now());
        tag.setUpdatedAt(LocalDateTime.now());
        return tagRepository.save(tag);
    }

    @Test
    void createTag_shouldReturnCreated() throws Exception {
        CreateTagRequest request = new CreateTagRequest("java", "#3A3A3A");

        mockMvc.perform(post("/api/v1/projects/{projectId}/tags", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.name").value("java"))
                .andExpect(jsonPath("$.color").value("#3A3A3A"));

        assertEquals(1, tagRepository.count());
    }

    @Test
    void createTag_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Name blank, invalid hex color
        CreateTagRequest request = new CreateTagRequest("   ", "not-a-color");

        mockMvc.perform(post("/api/v1/projects/{projectId}/tags", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(2)));
    }

    @Test
    void getAllTags_shouldReturnList() throws Exception {
        createAndSaveTag("docker", "#111");
        createAndSaveTag("spring", "#222");

        mockMvc.perform(get("/api/v1/projects/{projectId}/tags", savedProject.getId()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(2)))
                .andExpect(jsonPath("$[*].name", containsInAnyOrder("docker", "spring")));
    }

    @Test
    void getTagById_shouldReturnTag() throws Exception {
        TagEntity tag = createAndSaveTag("kubernetes", "#444");

        mockMvc.perform(get("/api/v1/projects/{projectId}/tags/{tagId}", savedProject.getId(), tag.getId()))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(tag.getId().toString()))
                .andExpect(jsonPath("$.name").value("kubernetes"))
                .andExpect(jsonPath("$.color").value("#444"));
    }

    @Test
    void getTagById_shouldReturnNotFound_whenTagDoesNotExist() throws Exception {
        UUID nonExistentId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/tags/{tagId}", savedProject.getId(), nonExistentId))
                .andExpect(status().isNotFound());
    }

    @Test
    void searchByName_shouldReturnMatchingTags() throws Exception {
        createAndSaveTag("react-js", "#111");
        createAndSaveTag("vue-js", "#222");
        createAndSaveTag("angular", "#333");

        mockMvc.perform(get("/api/v1/projects/{projectId}/tags/search", savedProject.getId())
                        .param("name", "js"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(2)))
                .andExpect(jsonPath("$[*].name", containsInAnyOrder("react-js", "vue-js")));
    }

    @Test
    void updateTag_shouldUpdateFields() throws Exception {
        TagEntity tag = createAndSaveTag("old-name", "#111");

        UpdateTagRequest request = new UpdateTagRequest("new-name", "#555");

        mockMvc.perform(patch("/api/v1/projects/{projectId}/tags/{tagId}", savedProject.getId(), tag.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.name").value("new-name"))
                .andExpect(jsonPath("$.color").value("#555"));

        TagEntity updated = tagRepository.findById(tag.getId()).orElseThrow();
        assertEquals("new-name", updated.getName());
        assertEquals("#555", updated.getColor());
    }

    @Test
    void deleteTag_shouldRemoveTag() throws Exception {
        TagEntity tag = createAndSaveTag("temp", "#999");

        mockMvc.perform(delete("/api/v1/projects/{projectId}/tags/{tagId}", savedProject.getId(), tag.getId()))
                .andExpect(status().isNoContent());

        assertFalse(tagRepository.existsById(tag.getId()));
    }
}
