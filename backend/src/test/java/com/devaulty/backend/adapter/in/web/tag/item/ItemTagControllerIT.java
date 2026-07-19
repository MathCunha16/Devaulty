package com.devaulty.backend.adapter.in.web.tag.item;

import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.snippet.SnippetEntity;
import com.devaulty.backend.adapter.out.persistence.snippet.SpringDataSnippetRepository;
import com.devaulty.backend.adapter.out.persistence.tag.SpringDataTagRepository;
import com.devaulty.backend.adapter.out.persistence.tag.TagEntity;
import com.devaulty.backend.adapter.out.persistence.tag.item.ItemTagEntity;
import com.devaulty.backend.adapter.out.persistence.tag.item.ItemTagId;
import com.devaulty.backend.adapter.out.persistence.tag.item.SpringDataItemTagRepository;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;

import java.time.LocalDateTime;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.delete;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.put;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

class ItemTagControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataItemTagRepository itemTagRepository;

    @Autowired
    private SpringDataTagRepository tagRepository;

    @Autowired
    private SpringDataSnippetRepository snippetRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;
    private TagEntity savedTag;
    private SnippetEntity savedSnippet;

    @BeforeEach
    void setUpData() {
        itemTagRepository.deleteAll();
        tagRepository.deleteAll();
        snippetRepository.deleteAll();
        projectRepository.deleteAll();

        ProjectEntity project = new ProjectEntity(
                UUID.randomUUID(),
                "ItemTag Integration Project",
                "Description",
                "folder",
                "#FFF",
                false
        );
        project.setCreatedAt(LocalDateTime.now());
        savedProject = projectRepository.save(project);

        TagEntity tag = new TagEntity(UUID.randomUUID(), savedProject, "docker", "#123");
        tag.setCreatedAt(LocalDateTime.now());
        tag.setUpdatedAt(LocalDateTime.now());
        savedTag = tagRepository.save(tag);

        SnippetEntity snippet = new SnippetEntity(
                UUID.randomUUID(),
                savedProject,
                "My Docker Compose",
                "Compose desc",
                "version: '3'",
                SnippetLanguage.YAML,
                SnippetType.CODE
        );
        snippet.setCreatedAt(LocalDateTime.now());
        savedSnippet = snippetRepository.save(snippet);
    }

    @Test
    void associateTag_shouldCreateAssociation() throws Exception {
        mockMvc.perform(put("/api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}",
                        savedProject.getId(), "snippet", savedSnippet.getId(), savedTag.getId()))
                .andExpect(status().isNoContent());

        ItemTagId associationId = new ItemTagId(savedTag.getId(), "snippet", savedSnippet.getId());
        assertTrue(itemTagRepository.existsById(associationId));
    }

    @Test
    void associateTag_shouldReturnNotFound_whenTagDoesNotExist() throws Exception {
        UUID nonExistentTagId = UUID.randomUUID();

        mockMvc.perform(put("/api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}",
                        savedProject.getId(), "snippet", savedSnippet.getId(), nonExistentTagId))
                .andExpect(status().isNotFound());
    }

    @Test
    void associateTag_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();

        mockMvc.perform(put("/api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}",
                        nonExistentProjectId, "snippet", savedSnippet.getId(), savedTag.getId()))
                .andExpect(status().isNotFound());
    }

    @Test
    void disassociateTag_shouldRemoveAssociation() throws Exception {
        ItemTagEntity association = new ItemTagEntity(savedTag, "snippet", savedSnippet.getId());
        itemTagRepository.save(association);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/items/{itemType}/{itemId}/tags/{tagId}",
                        savedProject.getId(), "snippet", savedSnippet.getId(), savedTag.getId()))
                .andExpect(status().isNoContent());

        ItemTagId associationId = new ItemTagId(savedTag.getId(), "snippet", savedSnippet.getId());
        assertFalse(itemTagRepository.existsById(associationId));
    }
}
