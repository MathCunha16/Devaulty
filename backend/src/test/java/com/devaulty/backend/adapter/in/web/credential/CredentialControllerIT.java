package com.devaulty.backend.adapter.in.web.credential;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.out.persistence.project.ProjectEntity;
import com.devaulty.backend.adapter.out.persistence.project.SpringDataProjectRepository;
import com.devaulty.backend.adapter.out.persistence.credential.CredentialEntity;
import com.devaulty.backend.adapter.out.persistence.credential.SpringDataCredentialRepository;
import com.devaulty.backend.adapter.out.persistence.setting.AppSettingEntity;
import com.devaulty.backend.adapter.out.persistence.setting.SpringDataAppSettingRepository;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.UUID;

import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

class CredentialControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataCredentialRepository credentialRepository;

    @Autowired
    private SpringDataProjectRepository projectRepository;

    @Autowired
    private SpringDataAppSettingRepository appSettingRepository;

    @Autowired
    private MasterKeySessionPort masterKeySessionPort;

    @Autowired
    private CryptoPort cryptoPort;

    private ProjectEntity savedProject;
    private SecretKey mockSecretKey;

    @BeforeEach
    void setUpData() {
        credentialRepository.deleteAll();
        projectRepository.deleteAll();
        appSettingRepository.deleteAll();

        // 1. Setup global configuration so setup is not required
        appSettingRepository.save(new AppSettingEntity("master_password_hash", "mockHash"));
        appSettingRepository.save(new AppSettingEntity("master_password_salt", "mockSalt"));

        // 2. Unlock vault with dummy AES-256 key (32 bytes of zeros)
        mockSecretKey = new SecretKeySpec(new byte[32], "AES");
        masterKeySessionPort.setKey(mockSecretKey);

        // 3. Create a project
        ProjectEntity project = new ProjectEntity(
                UUID.randomUUID(),
                "Integration Project",
                "Description",
                "folder",
                "#FFF",
                false
        );
        project.setCreatedAt(LocalDateTime.now());
        savedProject = projectRepository.save(project);
    }

    @Test
    void createCredential_shouldReturnCreated_whenLoginRequestIsValid() throws Exception {
        String json = """
                {
                  "title": "My Login DB",
                  "secretType": "LOGIN",
                  "username": "dbuser",
                  "password": "dbpassword",
                  "notes": "some notes",
                  "relatedUrl": "http://db-server"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.projectId").value(savedProject.getId().toString()))
                .andExpect(jsonPath("$.title").value("My Login DB"))
                .andExpect(jsonPath("$.secretType").value("LOGIN"))
                .andExpect(jsonPath("$.decryptedPayload.username").value("dbuser"))
                .andExpect(jsonPath("$.decryptedPayload.password").value("dbpassword"))
                .andExpect(jsonPath("$.notes").value("some notes"))
                .andExpect(jsonPath("$.relatedUrl").value("http://db-server"))
                .andExpect(jsonPath("$.createdAt").isNotEmpty())
                .andExpect(jsonPath("$.tags").isEmpty());

        assertEquals(1, credentialRepository.count());

        // Reload the entity to assert at-rest protection and context-binding AAD decryption
        CredentialEntity reloaded = credentialRepository.findAll().getFirst();
        assertNotNull(reloaded.getPayloadEncrypted());
        assertFalse(new String(reloaded.getPayloadEncrypted()).contains("dbuser"));
        assertFalse(new String(reloaded.getPayloadEncrypted()).contains("dbpassword"));

        byte[] decrypted = cryptoPort.decrypt(
                reloaded.getPayloadEncrypted(),
                reloaded.getEncryptionIv(),
                reloaded.getEncryptionAuthTag(),
                mockSecretKey,
                computeAad(savedProject.getId(), reloaded.getId())
        );

        String decryptedStr = new String(decrypted);
        assertTrue(decryptedStr.contains("dbuser"));
        assertTrue(decryptedStr.contains("dbpassword"));
    }

    @Test
    void createCredential_shouldReturnCreated_whenApiKeyRequestIsValid() throws Exception {
        String json = """
                {
                  "title": "My Api Key",
                  "secretType": "API_KEY",
                  "apiKey": "key_12345"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.secretType").value("API_KEY"))
                .andExpect(jsonPath("$.decryptedPayload.apiKey").value("key_12345"));
    }

    @Test
    void createCredential_shouldReturnCreated_whenRawTextRequestIsValid() throws Exception {
        String json = """
                {
                  "title": "My Raw text",
                  "secretType": "RAW_TEXT",
                  "rawTextContent": "private certificate text"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.secretType").value("RAW_TEXT"))
                .andExpect(jsonPath("$.decryptedPayload.rawText").value("private certificate text"));
    }

    @Test
    void createCredential_shouldReturnBadRequest_whenSecretFieldsAreMissingForLogin() throws Exception {
        String json = """
                {
                  "title": "Invalid Login",
                  "secretType": "LOGIN",
                  "username": "admin"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors[0].field").value("validSecretPayload"))
                .andExpect(jsonPath("$.errors[0].message").value("Required secret fields are missing for the selected secret type"));
    }

    @Test
    void createCredential_shouldReturnBadRequest_whenTitleIsTooShort() throws Exception {
        String json = """
                {
                  "title": "A",
                  "secretType": "API_KEY",
                  "apiKey": "key"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors[0].field").value("title"))
                .andExpect(jsonPath("$.errors[0].message").value("Title must be between 2 and 255 characters"));
    }

    @Test
    void createCredential_shouldReturnNotFound_whenProjectDoesNotExist() throws Exception {
        UUID nonExistentProjectId = UUID.randomUUID();
        String json = """
                {
                  "title": "Title",
                  "secretType": "API_KEY",
                  "apiKey": "key"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", nonExistentProjectId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.status").value(404));
    }

    @Test
    void createCredential_shouldReturnLocked_whenVaultIsLocked() throws Exception {
        masterKeySessionPort.clear(); // Lock vault

        String json = """
                {
                  "title": "Title",
                  "secretType": "API_KEY",
                  "apiKey": "key"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isLocked())
                .andExpect(jsonPath("$.message").value("Vault is locked, please unlock it first"));
    }

    @Test
    void createCredential_shouldReturnForbidden_whenSetupIsRequired() throws Exception {
        appSettingRepository.deleteAll(); // Force setup required

        String json = """
                {
                  "title": "Title",
                  "secretType": "API_KEY",
                  "apiKey": "key"
                }
                """;

        mockMvc.perform(post("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.message").value("Master password not configured"));
    }

    @Test
    void getAllCredentials_shouldReturnPagedSummaries() throws Exception {
        CredentialEntity c1 = new CredentialEntity(
                UUID.randomUUID(),
                savedProject,
                "Cred 1",
                CredentialSecretType.LOGIN,
                new byte[]{1, 2},
                new byte[]{3},
                new byte[]{4},
                "notes",
                "http://url1"
        );
        c1.setCreatedAt(LocalDateTime.now().minusHours(1));

        CredentialEntity c2 = new CredentialEntity(
                UUID.randomUUID(),
                savedProject,
                "Cred 2",
                CredentialSecretType.API_KEY,
                new byte[]{5, 6},
                new byte[]{7},
                new byte[]{8},
                "notes 2",
                "http://url2"
        );
        c2.setCreatedAt(LocalDateTime.now());

        credentialRepository.save(c1);
        credentialRepository.save(c2);

        mockMvc.perform(get("/api/v1/projects/{projectId}/credentials", savedProject.getId())
                        .param("page", "0")
                        .param("size", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(2)))
                .andExpect(jsonPath("$.content[0].title").value("Cred 2"))
                .andExpect(jsonPath("$.content[0].tags").isEmpty())
                .andExpect(jsonPath("$.content[1].title").value("Cred 1"))
                .andExpect(jsonPath("$.content[1].tags").isEmpty())
                .andExpect(jsonPath("$.page.totalElements").value(2));
    }

    @Test
    void getCredentialById_shouldReturnDecryptedCredential() throws Exception {
        UUID credentialId = UUID.randomUUID();
        byte[] payload = "{\"username\":\"admin\",\"password\":\"secret\"}".getBytes();
        CryptoResultDto cryptoResult = cryptoPort.encrypt(payload, mockSecretKey, computeAad(savedProject.getId(), credentialId));

        CredentialEntity entity = new CredentialEntity(
                credentialId,
                savedProject,
                "Get Credential",
                CredentialSecretType.LOGIN,
                cryptoResult.cipherText(),
                cryptoResult.iv(),
                cryptoResult.authTag(),
                "notes",
                "http://url"
        );
        entity.setCreatedAt(LocalDateTime.now());
        credentialRepository.save(entity);

        mockMvc.perform(get("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), credentialId))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(credentialId.toString()))
                .andExpect(jsonPath("$.title").value("Get Credential"))
                .andExpect(jsonPath("$.decryptedPayload.username").value("admin"))
                .andExpect(jsonPath("$.decryptedPayload.password").value("secret"))
                .andExpect(jsonPath("$.tags").isEmpty());
    }

    @Test
    void getCredentialById_shouldReturnNotFound_whenCredentialDoesNotExist() throws Exception {
        UUID nonExistentId = UUID.randomUUID();

        mockMvc.perform(get("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), nonExistentId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Credential not found with identifier " + nonExistentId));
    }

    @Test
    void getCredentialById_shouldReturnNotFound_whenCredentialDoesNotBelongToProject() throws Exception {
        // Create another project
        ProjectEntity otherProject = new ProjectEntity(UUID.randomUUID(), "Other", "Desc", "folder", "#000", false);
        otherProject.setCreatedAt(LocalDateTime.now());
        projectRepository.save(otherProject);

        UUID credentialId = UUID.randomUUID();
        byte[] payload = "{\"apiKey\":\"secret_key\"}".getBytes();
        CryptoResultDto cryptoResult = cryptoPort.encrypt(payload, mockSecretKey, computeAad(otherProject.getId(), credentialId));

        CredentialEntity entity = new CredentialEntity(
                credentialId,
                otherProject,
                "Other project's credential",
                CredentialSecretType.API_KEY,
                cryptoResult.cipherText(),
                cryptoResult.iv(),
                cryptoResult.authTag(),
                null,
                null
        );
        entity.setCreatedAt(LocalDateTime.now());
        credentialRepository.save(entity);

        // Fetch using savedProject's ID (not otherProject's ID)
        mockMvc.perform(get("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), credentialId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Credential not found with identifier " + credentialId));
    }

    @Test
    void updateCredential_shouldUpdateTitleAndMetadataOnly_whenNoPayloadProvided() throws Exception {
        UUID credentialId = UUID.randomUUID();
        byte[] payload = "{\"apiKey\":\"orig_key\"}".getBytes();
        CryptoResultDto cryptoResult = cryptoPort.encrypt(payload, mockSecretKey, computeAad(savedProject.getId(), credentialId));

        CredentialEntity entity = new CredentialEntity(
                credentialId,
                savedProject,
                "Old Title",
                CredentialSecretType.API_KEY,
                cryptoResult.cipherText(),
                cryptoResult.iv(),
                cryptoResult.authTag(),
                "Old notes",
                "http://old"
        );
        entity.setCreatedAt(LocalDateTime.now());
        credentialRepository.save(entity);

        String json = """
                {
                  "title": "New Title",
                  "notes": "New notes",
                  "relatedUrl": "http://new"
                }
                """;

        mockMvc.perform(patch("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), credentialId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.title").value("New Title"))
                .andExpect(jsonPath("$.decryptedPayload.apiKey").value("orig_key"))
                .andExpect(jsonPath("$.notes").value("New notes"))
                .andExpect(jsonPath("$.relatedUrl").value("http://new"))
                .andExpect(jsonPath("$.updatedAt").isNotEmpty());

        // Reload to assert metadata updated but ciphertext preserved unchanged!
        CredentialEntity reloaded = credentialRepository.findById(credentialId).orElseThrow();
        assertEquals("New Title", reloaded.getTitle());
        assertEquals("New notes", reloaded.getNotes());
        assertEquals("http://new", reloaded.getRelatedUrl());
        assertArrayEquals(cryptoResult.cipherText(), reloaded.getPayloadEncrypted());
        assertArrayEquals(cryptoResult.iv(), reloaded.getEncryptionIv());
        assertArrayEquals(cryptoResult.authTag(), reloaded.getEncryptionAuthTag());
    }

    @Test
    void updateCredential_shouldUpdatePayloadAndReencrypt_whenNewPayloadProvided() throws Exception {
        UUID credentialId = UUID.randomUUID();
        byte[] payload = "{\"apiKey\":\"orig_key\"}".getBytes();
        CryptoResultDto cryptoResult = cryptoPort.encrypt(payload, mockSecretKey, computeAad(savedProject.getId(), credentialId));

        CredentialEntity entity = new CredentialEntity(
                credentialId,
                savedProject,
                "Title",
                CredentialSecretType.API_KEY,
                cryptoResult.cipherText(),
                cryptoResult.iv(),
                cryptoResult.authTag(),
                null,
                null
        );
        entity.setCreatedAt(LocalDateTime.now());
        credentialRepository.save(entity);

        String json = """
                {
                  "secretType": "RAW_TEXT",
                  "rawTextContent": "new raw text content"
                }
                """;

        mockMvc.perform(patch("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), credentialId)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.secretType").value("RAW_TEXT"))
                .andExpect(jsonPath("$.decryptedPayload.rawText").value("new raw text content"));

        // Reload to assert ciphertext, IV, and tag changed, and decrypts successfully to new payload using AAD
        CredentialEntity reloaded = credentialRepository.findById(credentialId).orElseThrow();
        assertEquals(CredentialSecretType.RAW_TEXT, reloaded.getSecretType());
        assertFalse(Arrays.equals(cryptoResult.cipherText(), reloaded.getPayloadEncrypted()));
        assertFalse(Arrays.equals(cryptoResult.iv(), reloaded.getEncryptionIv()));
        assertFalse(Arrays.equals(cryptoResult.authTag(), reloaded.getEncryptionAuthTag()));

        byte[] decrypted = cryptoPort.decrypt(
                reloaded.getPayloadEncrypted(),
                reloaded.getEncryptionIv(),
                reloaded.getEncryptionAuthTag(),
                mockSecretKey,
                computeAad(savedProject.getId(), credentialId)
        );
        assertEquals("{\"rawText\":\"new raw text content\"}", new String(decrypted));
    }

    @Test
    void deleteCredential_shouldRemoveCredentialAndReturnNoContent() throws Exception {
        UUID credentialId = UUID.randomUUID();
        CredentialEntity entity = new CredentialEntity(
                credentialId,
                savedProject,
                "To delete",
                CredentialSecretType.API_KEY,
                new byte[]{1},
                new byte[]{2},
                new byte[]{3},
                null,
                null
        );
        entity.setCreatedAt(LocalDateTime.now());
        credentialRepository.save(entity);

        mockMvc.perform(delete("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), credentialId))
                .andExpect(status().isNoContent());

        assertFalse(credentialRepository.existsById(credentialId));
    }

    @Test
    void deleteCredential_shouldReturnNotFound_whenCredentialDoesNotExist() throws Exception {
        UUID nonExistentId = UUID.randomUUID();

        mockMvc.perform(delete("/api/v1/projects/{projectId}/credentials/{credentialId}", savedProject.getId(), nonExistentId))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message").value("Credential not found with identifier " + nonExistentId));
    }

    private byte[] computeAad(UUID projectId, UUID credentialId) {
        java.nio.ByteBuffer buffer = java.nio.ByteBuffer.allocate(32);
        buffer.putLong(projectId.getMostSignificantBits());
        buffer.putLong(projectId.getLeastSignificantBits());
        buffer.putLong(credentialId.getMostSignificantBits());
        buffer.putLong(credentialId.getLeastSignificantBits());
        return buffer.array();
    }
}
