package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;

import java.util.Objects;
import java.util.UUID;

public class Project extends BaseEntity {

    private UUID id;
    private String name;
    private String description;
    private String icon;
    private String color;
    private boolean archived;

    public Project() {
        // Empty constructor
    }

    public Project(UUID id, String name, String description, String icon, String color, boolean archived) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.icon = icon;
        this.color = color;
        this.archived = archived;
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getIcon() {
        return icon;
    }

    public void setIcon(String icon) {
        this.icon = icon;
    }

    public String getColor() {
        return color;
    }

    public void setColor(String color) {
        this.color = color;
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
        Project project = (Project) o;
        return Objects.equals(id, project.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
