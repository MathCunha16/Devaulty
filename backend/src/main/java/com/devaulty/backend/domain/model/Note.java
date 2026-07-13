package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;

import java.util.Objects;
import java.util.UUID;

public class Note extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String title;
    private String content;
    private boolean archived;

    public Note() {
        // Empty constructor
    }

    public Note(UUID id, UUID projectId, String title, String content, boolean archived) {
        this.id = id;
        this.projectId = projectId;
        this.title = title;
        this.content = content;
        this.archived = archived;
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

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public boolean isArchived() {
        return archived;
    }

    public void setArchived(boolean archived) {
        this.archived = archived;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Note note = (Note) o;
        return Objects.equals(id, note.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
