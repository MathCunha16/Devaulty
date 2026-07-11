package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;

import java.util.Objects;
import java.util.UUID;

public class Link extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String title;
    private String url;
    private String description;

    public Link() {
        // Empty constructor
    }

    public Link(UUID id, UUID projectId, String title, String url, String description) {
        this.id = id;
        this.projectId = projectId;
        this.title = title;
        this.url = url;
        this.description = description;
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

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Link link = (Link) o;
        return Objects.equals(id, link.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
