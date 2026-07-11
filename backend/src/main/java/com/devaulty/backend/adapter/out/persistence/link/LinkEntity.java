package com.devaulty.backend.adapter.out.persistence.link;

import com.devaulty.backend.adapter.out.persistence.base.BaseJpaEntity;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import jakarta.persistence.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.util.Objects;
import java.util.UUID;

@Table(name = "links")
@Entity
public class LinkEntity extends BaseJpaEntity {

    @Id
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_id", nullable = false)
    private ProjectEntity project;

    @Column(nullable = false)
    private String title;

    @Column(nullable = false)
    private String url;

    private String description;

    public LinkEntity() {
        // Empty constructor
    }

    public LinkEntity(UUID id, ProjectEntity project, String title, String url, String description) {
        this.id = id;
        this.project = project;
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
        LinkEntity that = (LinkEntity) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
