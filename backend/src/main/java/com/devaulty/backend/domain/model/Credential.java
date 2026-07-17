package com.devaulty.backend.domain.model;

import com.devaulty.backend.domain.model.base.BaseEntity;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;

import java.util.Objects;
import java.util.UUID;

public class Credential extends BaseEntity {

    private UUID id;
    private UUID projectId;
    private String title;
    private CredentialSecretType secretType;
    private byte[] payloadEncrypted;
    private byte[] encryptionIv;
    private byte[] encryptionAuthTag;
    private String notes;
    private String relatedUrl;

    public Credential() {
        // Empty constructor
    }

    public Credential(UUID id, UUID projectId, String title, CredentialSecretType secretType, byte[] payloadEncrypted, byte[] encryptionIv, byte[] encryptionAuthTag, String notes, String relatedUrl) {
        this.id = id;
        this.projectId = projectId;
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
        Credential that = (Credential) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(id);
    }
}
