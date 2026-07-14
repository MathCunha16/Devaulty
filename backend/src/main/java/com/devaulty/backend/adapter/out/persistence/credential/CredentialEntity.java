package com.devaulty.backend.adapter.out.persistence.credential;

import com.devaulty.backend.adapter.out.persistence.base.BaseJpaEntity;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import jakarta.persistence.*;
import org.hibernate.annotations.JdbcTypeCode;
import org.hibernate.type.SqlTypes;

import java.util.Objects;
import java.util.UUID;

@Table(name = "credentials")
@Entity
public class CredentialEntity extends BaseJpaEntity {

    @Id
    @JdbcTypeCode(SqlTypes.VARCHAR)
    private UUID id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "project_id", nullable = false)
    private ProjectEntity project;

    @Column(nullable = false)
    private String title;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private CredentialSecretType secretType;

    @Column(name = "payload_encrypted", nullable = false)
    private byte[] payloadEncrypted;

    @Column(name = "encryption_iv", nullable = false)
    private byte[] encryptionIv;

    @Column(name = "encryption_auth_tag", nullable = false)
    private byte[] encryptionAuthTag;

    private String notes;

    @Column(name = "related_url")
    private String relatedUrl;

    public CredentialEntity() {
        // Empty constructor
    }

    public CredentialEntity(UUID id, ProjectEntity project, String title, CredentialSecretType secretType, byte[] payloadEncrypted, byte[] encryptionIv, byte[] encryptionAuthTag, String notes, String relatedUrl) {
        this.id = id;
        this.project = project;
        this.title = title;
        this.secretType = secretType;
        this.payloadEncrypted = payloadEncrypted;
        this.encryptionIv = encryptionIv;
        this.encryptionAuthTag = encryptionAuthTag;
        this.notes = notes;
        this.relatedUrl = relatedUrl;
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

    public CredentialSecretType getSecretType() {
        return secretType;
    }

    public void setSecretType(CredentialSecretType secretType) {
        this.secretType = secretType;
    }

    public byte[] getPayloadEncrypted() {
        return payloadEncrypted;
    }

    public void setPayloadEncrypted(byte[] payloadEncrypted) {
        this.payloadEncrypted = payloadEncrypted;
    }

    public byte[] getEncryptionIv() {
        return encryptionIv;
    }

    public void setEncryptionIv(byte[] encryptionIv) {
        this.encryptionIv = encryptionIv;
    }

    public byte[] getEncryptionAuthTag() {
        return encryptionAuthTag;
    }

    public void setEncryptionAuthTag(byte[] encryptionAuthTag) {
        this.encryptionAuthTag = encryptionAuthTag;
    }

    public String getNotes() {
        return notes;
    }

    public void setNotes(String notes) {
        this.notes = notes;
    }

    public String getRelatedUrl() {
        return relatedUrl;
    }

    public void setRelatedUrl(String relatedUrl) {
        this.relatedUrl = relatedUrl;
    }

    @Override
    public boolean equals(Object o) {
        if (o == null || getClass() != o.getClass()) return false;
        CredentialEntity that = (CredentialEntity) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
