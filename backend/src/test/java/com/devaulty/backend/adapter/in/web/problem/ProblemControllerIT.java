package com.devaulty.backend.adapter.in.web.problem;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.in.web.problem.dto.*;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.problem.ProblemEntity;
import com.devaulty.backend.adapter.out.persistence.problem.SpringDataProblemRepository;
import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;
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

class ProblemControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataProblemRepository problemRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    private ProjectEntity savedProject;

    @BeforeEach
    void setUpData() {
        problemRepository.deleteAll();
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
    void createProblem_shouldReturnCreated() throws Exception {
        CreateProblemRequest request = new CreateProblemRequest(
                "My Problem",
                "Exception stacktrace",
                "Fix configuration",
                ProblemStatus.OPEN,
                ProblemSeverity.HIGH
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/problems", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(header().exists("Location"))
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").value("My Problem"))
                .andExpect(jsonPath("$.errorDescription").value("Exception stacktrace"))
                .andExpect(jsonPath("$.solution").value("Fix configuration"))
                .andExpect(jsonPath("$.status").value("OPEN"))
                .andExpect(jsonPath("$.severity").value("HIGH"))
                .andExpect(jsonPath("$.createdAt").isNotEmpty())
                .andExpect(jsonPath("$.tags").isEmpty());

        assertEquals(1, problemRepository.count());
    }

    @Test
    void createProblem_shouldReturnBadRequest_whenValidationFails() throws Exception {
        // Title blank (min 2), Status null, Severity null
        CreateProblemRequest request = new CreateProblemRequest(
                "   ",
                "Error description",
                "Solution",
                null,
                null
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/problems", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors", hasSize(3)));
    }

    @Test
    void createProblem_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        CreateProblemRequest request = new CreateProblemRequest(
                "Problem Title",
                "description",
                "solution",
                ProblemStatus.OPEN,
                ProblemSeverity.MEDIUM
        );

        mockMvc.perform(post("/api/v1/projects/{projectId}/problems", nonExistentProjectId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404))
                .andExpect(jsonPath("$.message").value("Project not found with identifier " + nonExistentProjectId));
    }

    @Test
    void getAllProblems_shouldReturnPagedProblems() throws Exception {
        ProblemEntity p1 = new ProblemEntity(UUID.randomUUID(), savedProject, "Problem 1", "Error 1", "Sol 1", ProblemStatus.OPEN, ProblemSeverity.MEDIUM);
        p1.setCreatedAt(LocalDateTime.now().minusHours(1));
        ProblemEntity p2 = new ProblemEntity(UUID.randomUUID(), savedProject, "Problem 2", "Error 2", "Sol 2", ProblemStatus.RESOLVED, ProblemSeverity.HIGH);
        p2.setCreatedAt(LocalDateTime.now());
        problemRepository.save(p1);
        problemRepository.save(p2);

        mockMvc.perform(get("/api/v1/projects/{projectId}/problems", savedProject.getId())
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].title").value("Problem 2"))
                .andExpect(jsonPath("$.content[0].tags").isEmpty())
                .andExpect(jsonPath("$.content[1].title").value("Problem 1"))
                .andExpect(jsonPath("$.content[1].tags").isEmpty())
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getAllProblems_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/problems", nonExistentProjectId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404));
    }

    @Test
    void getProblemById_shouldReturnProblem() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "Get Problem", "Error", "Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        mockMvc.perform(get("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), id))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(id.toString()))
                .andExpect(jsonPath("$.title").value("Get Problem"))
                .andExpect(jsonPath("$.tags").isEmpty());
    }

    @Test
    void getProblemById_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/problems/{problemId}", nonExistentProjectId, problemId))
                .andExpect(status().isNotFound());
    }

    @Test
    void getProblemById_shouldReturnNotFound_whenProblemDoesNotExist() throws Exception {
        UUID nonExistentProblemId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), nonExistentProblemId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Problem not found with identifier " + nonExistentProblemId));
    }

    @Test
    void getProblemById_shouldReturnNotFound_whenProblemDoesNotBelongToProject() throws Exception {
        // Create another project
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        // Create problem for other project
        UUID problemId = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(problemId, otherProject, "Other Problem", "Error", "Sol", ProblemStatus.OPEN, ProblemSeverity.MEDIUM);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        // Request other project's problem under savedProject's path
        mockMvc.perform(get("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), problemId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Problem not found with identifier " + problemId));
    }

    @Test
    void updateProblem_shouldUpdateAndReturnProblem() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "Old Title", "Old Error", "Old Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        UpdateProblemRequest request = new UpdateProblemRequest(
                "New Title",
                "New Error",
                "New Sol",
                ProblemSeverity.HIGH
        );

        mockMvc.perform(patch("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.title").value("New Title"))
                .andExpect(jsonPath("$.errorDescription").value("New Error"))
                .andExpect(jsonPath("$.solution").value("New Sol"))
                .andExpect(jsonPath("$.severity").value("HIGH"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        ProblemEntity updated = problemRepository.findById(id).orElseThrow();
        assertEquals("New Title", updated.getTitle());
    }

    @Test
    void updateProblem_shouldReturnBadRequest_whenValidationFails() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "Old Title", "Old Error", "Old Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        // Blank title
        UpdateProblemRequest request = new UpdateProblemRequest("   ", null, null, null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"));
    }

    @Test
    void updateProblemStatus_shouldUpdateAndReturnProblem() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "Title", "Error", "Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        UpdateProblemStatusRequest request = new UpdateProblemStatusRequest(ProblemStatus.RESOLVED);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/problems/{problemId}/status", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").value("RESOLVED"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        ProblemEntity updated = problemRepository.findById(id).orElseThrow();
        assertEquals(ProblemStatus.RESOLVED, updated.getStatus());
    }

    @Test
    void updateProblemStatus_shouldReturnBadRequest_whenStatusIsNull() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "Title", "Error", "Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        UpdateProblemStatusRequest request = new UpdateProblemStatusRequest(null);

        mockMvc.perform(patch("/api/v1/projects/{projectId}/problems/{problemId}/status", savedProject.getId(), id)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    @Test
    void deleteProblem_shouldDeleteAndReturnNoContent() throws Exception {
        UUID id = UUID.randomUUID();
        ProblemEntity problem = new ProblemEntity(id, savedProject, "To Delete", "Error", "Sol", ProblemStatus.OPEN, ProblemSeverity.LOW);
        problem.setCreatedAt(LocalDateTime.now());
        problemRepository.save(problem);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), id))
                .andExpect(status().isNoContent());

        assertFalse(problemRepository.existsById(id));
    }

    @Test
    void deleteProblem_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        UUID problemId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/problems/{problemId}", nonExistentProjectId, problemId))
                .andExpect(status().isNotFound());
    }

    @Test
    void deleteProblem_shouldReturnNotFound_whenProblemDoesNotExist() throws Exception {
        UUID nonExistentProblemId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/problems/{problemId}", savedProject.getId(), nonExistentProblemId))
                .andExpect(status().isNotFound());
    }
}
