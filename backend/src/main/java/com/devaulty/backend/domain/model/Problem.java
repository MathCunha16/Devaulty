package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;
import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;

import java.util.Objects;
import java.util.UUID;

public class Problem extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String title;
    private String errorDescription;
    private String solution;
    private ProblemStatus status = ProblemStatus.OPEN;
    private ProblemSeverity severity = ProblemSeverity.MEDIUM; // Medium by default

    public Problem() {
        // Empty constructor
    }

    public Problem(UUID id, UUID projectId, String title, String errorDescription, String solution, ProblemStatus status, ProblemSeverity severity) {
        this.id = id;
        this.projectId = projectId;
        this.title = title;
        this.errorDescription = errorDescription;
        this.solution = solution;
        this.status = status;
        this.severity = severity;
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public UUID getProjectId() {
        return projectId;
    }

    public void setProjectId(UUID projectId) {
        this.projectId = projectId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getErrorDescription() {
        return errorDescription;
    }

    public void setErrorDescription(String errorDescription) {
        this.errorDescription = errorDescription;
    }

    public String getSolution() {
        return solution;
    }

    public void setSolution(String solution) {
        this.solution = solution;
    }

    public ProblemStatus getStatus() {
        return status;
    }

    public void setStatus(ProblemStatus status) {
        this.status = status;
    }

    public ProblemSeverity getSeverity() {
        return severity;
    }

    public void setSeverity(ProblemSeverity severity) {
        this.severity = severity;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Problem problem = (Problem) o;
        return Objects.equals(id, problem.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
