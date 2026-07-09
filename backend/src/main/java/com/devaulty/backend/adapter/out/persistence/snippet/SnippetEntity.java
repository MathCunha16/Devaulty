package com.devaulty.backend.adapter.out.persistence.snippet;

import com.devaulty.backend.adapter.out.persistence.base.BaseJpaEntity;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;
import jakarta.persistence.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.util.Objects;
import java.util.UUID;

@Table(name = "snippets")
@Entity
public class SnippetEntity extends BaseJpaEntity {
    @Id
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_id", nullable = false)
    private ProjectEntity project;

    @Column(nullable = false)
    private String title;

    private String description;

    @Column(nullable = false)
    private String content;

    @Enumerated(EnumType.STRING)
    private SnippetLanguage language;

    @Column(name = "snippet_type", nullable = false)
    @Enumerated(EnumType.STRING)
    private SnippetType snippetType = SnippetType.COMMAND;

    public SnippetEntity() {
        // Empty constructor
    }

    public SnippetEntity(UUID id, ProjectEntity project, String title, String description, String content, SnippetLanguage language, SnippetType snippetType) {
        this.id = id;
        this.project = project;
        this.title = title;
        this.description = description;
        this.content = content;
        this.language = language;
        this.snippetType = snippetType;
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

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public SnippetLanguage getLanguage() {
        return language;
    }

    public void setLanguage(SnippetLanguage language) {
        this.language = language;
    }

    public SnippetType getSnippetType() {
        return snippetType;
    }

    public void setSnippetType(SnippetType snippetType) {
        this.snippetType = snippetType;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        SnippetEntity that = (SnippetEntity) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
