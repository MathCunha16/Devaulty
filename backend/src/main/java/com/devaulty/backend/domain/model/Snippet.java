package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;
import com.devaulty.backend.domain.model.enums.SnippetLanguage;
import com.devaulty.backend.domain.model.enums.SnippetType;

import java.util.Objects;
import java.util.UUID;

public class Snippet extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String title;
    private String description;
    private String content;
    private SnippetLanguage language;
    private SnippetType snippetType = SnippetType.COMMAND;

    public Snippet() {
        // Empty constructor
    }

    public Snippet(UUID id, UUID projectId, String title, String description, String content, SnippetLanguage language, SnippetType snippetType) {
        this.id = id;
        this.projectId = projectId;
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
        Snippet snippet = (Snippet) o;
        return Objects.equals(id, snippet.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
