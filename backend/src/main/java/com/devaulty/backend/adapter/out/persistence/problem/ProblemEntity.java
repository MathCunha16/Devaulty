package com.devaulty.backend.adapter.out.persistence.problem;

import com.devaulty.backend.adapter.out.persistence.base.BaseJpaEntity;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.domain.model.enums.ProblemSeverity;
import com.devaulty.backend.domain.model.enums.ProblemStatus;
import jakarta.persistence.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.util.Objects;
import java.util.UUID;

@Table(name = "problems")
@Entity
public class ProblemEntity extends BaseJpaEntity {

    @Id
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_id", nullable = false)
    private ProjectEntity project;

    @Column(nullable = false)
    private String title;

    @Column(name = "error_description")
    private String errorDescription;

    private String solution;

    @Enumerated(EnumType.STRING)
    private ProblemStatus status = ProblemStatus.OPEN;

    @Enumerated(EnumType.STRING)
    private ProblemSeverity severity = ProblemSeverity.MEDIUM;

    public ProblemEntity() {
        // Empty constructor
    }

    public ProblemEntity(UUID id, ProjectEntity project, String title, String errorDescription, String solution, ProblemStatus status, ProblemSeverity severity) {
        this.id = id;
        this.project = project;
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

    public ProjectEntity getProject() {
        return project;
    }

    public void setProject(ProjectEntity project) {
        this.project = project;
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
        ProblemEntity that = (ProblemEntity) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
