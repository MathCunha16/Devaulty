package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;

import java.util.Objects;
import java.util.UUID;

public class Tag extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String name;
    private String color;

    public Tag() {
        // Empty constructor
    }

    public Tag(UUID id, UUID projectId, String name, String color) {
        this.id = id;
        this.projectId = projectId;
        this.name = name;
        this.color = color;
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

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getColor() {
        return color;
    }

    public void setColor(String color) {
        this.color = color;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        Tag tag = (Tag) o;
        return Objects.equals(id, tag.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
