package com.devaulty.backend.adapter.in.web.security;

import com.devaulty.backend.infrastructure.BaseIntegrationTest;
import com.devaulty.backend.adapter.out.persistence.setting.AppSettingEntity;
import com.devaulty.backend.adapter.out.persistence.setting.SpringDataAppSettingRepository;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;

import javax.crypto.spec.SecretKeySpec;

import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

class SecurityControllerIT extends BaseIntegrationTest {

    @Autowired
    private SpringDataAppSettingRepository appSettingRepository;

    @Autowired
    private MasterKeySessionPort masterKeySessionPort;

    @BeforeEach
    void setUpData() {
        appSettingRepository.deleteAll();
        masterKeySessionPort.clear();
    }

    @Test
    void setupMasterPassword_shouldReturnNoContent_whenFirstTime() throws Exception {
        String json = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/master-password")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isNoContent());

        assertTrue(appSettingRepository.existsById("master_password_hash"));
        assertTrue(appSettingRepository.existsById("master_password_salt"));

        AppSettingEntity hashEntity = appSettingRepository.findById("master_password_hash").orElseThrow();
        AppSettingEntity saltEntity = appSettingRepository.findById("master_password_salt").orElseThrow();

        assertFalse(hashEntity.getValue().isEmpty());
        assertFalse(saltEntity.getValue().isEmpty());
        assertNotEquals("mySuperSecretPassword123", hashEntity.getValue());

        assertTrue(masterKeySessionPort.hasKey());
    }

    @Test
    void setupMasterPassword_shouldReturnBadRequest_whenPasswordTooShort() throws Exception {
        String json = """
                {
                  "masterPassword": "short"
                }
                """;

        mockMvc.perform(post("/api/v1/security/master-password")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors[0].field").value("masterPassword"))
                .andExpect(jsonPath("$.errors[0].message").value("Password must be between 8 and 255 characters"));

        assertFalse(appSettingRepository.existsById("master_password_hash"));
        assertFalse(masterKeySessionPort.hasKey());
    }

    @Test
    void setupMasterPassword_shouldReturnConflict_whenAlreadyConfigured() throws Exception {
        appSettingRepository.save(new AppSettingEntity("master_password_hash", "mockHash"));
        appSettingRepository.save(new AppSettingEntity("master_password_salt", "mockSalt"));

        String json = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/master-password")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isConflict())
                .andExpect(jsonPath("$.message").value("Master password already configured"));
    }

    @Test
    void checkMasterPasswordSetup_shouldReturnTrue_whenNotConfigured() throws Exception {
        mockMvc.perform(get("/api/v1/security/master-password/required-status"))
                .andExpect(status().isOk())
                .andExpect(content().string("true"));
    }

    @Test
    void checkMasterPasswordSetup_shouldReturnFalse_whenAlreadyConfigured() throws Exception {
        appSettingRepository.save(new AppSettingEntity("master_password_hash", "mockHash"));
        appSettingRepository.save(new AppSettingEntity("master_password_salt", "mockSalt"));

        mockMvc.perform(get("/api/v1/security/master-password/required-status"))
                .andExpect(status().isOk())
                .andExpect(content().string("false"));
    }

    @Test
    void getSessionStatus_shouldReturnInactiveAndZeroSeconds_whenLocked() throws Exception {
        mockMvc.perform(get("/api/v1/security/vault/status"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.active").value(false))
                .andExpect(jsonPath("$.secondsLeft").value(0));
    }

    @Test
    void getSessionStatus_shouldReturnActiveAndRemainingSeconds_whenUnlocked() throws Exception {
        masterKeySessionPort.setKey(new SecretKeySpec(new byte[32], "AES"));

        mockMvc.perform(get("/api/v1/security/vault/status"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.active").value(true))
                .andExpect(jsonPath("$.secondsLeft").value(greaterThan(0)));
    }

    @Test
    void unlockVault_shouldReturnTrue_whenPasswordIsCorrect() throws Exception {
        // 1. Setup a master password first
        String setupJson = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/master-password")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(setupJson))
                .andExpect(status().isNoContent());

        // Lock vault for testing unlock explicitly
        masterKeySessionPort.clear();
        assertFalse(masterKeySessionPort.hasKey());

        // 2. Unlock vault
        String unlockJson = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/vault/unlock")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(unlockJson))
                .andExpect(status().isOk())
                .andExpect(content().string("true"));

        assertTrue(masterKeySessionPort.hasKey());
    }

    @Test
    void unlockVault_shouldReturnUnauthorized_whenPasswordIsWrong() throws Exception {
        // 1. Setup a master password
        String setupJson = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/master-password")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(setupJson))
                .andExpect(status().isNoContent());

        masterKeySessionPort.clear();

        // 2. Try to unlock with wrong password
        String unlockJson = """
                {
                  "masterPassword": "wrongPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/vault/unlock")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(unlockJson))
                .andExpect(status().isUnauthorized())
                .andExpect(jsonPath("$.message").value("Invalid MasterPassword!"));

        assertFalse(masterKeySessionPort.hasKey());
    }

    @Test
    void unlockVault_shouldReturnForbidden_whenNotConfigured() throws Exception {
        String json = """
                {
                  "masterPassword": "mySuperSecretPassword123"
                }
                """;

        mockMvc.perform(post("/api/v1/security/vault/unlock")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isForbidden())
                .andExpect(jsonPath("$.message").value("Master password not configured"));
    }

    @Test
    void unlockVault_shouldReturnBadRequest_whenPasswordTooShort() throws Exception {
        String json = """
                {
                  "masterPassword": "short"
                }
                """;

        mockMvc.perform(post("/api/v1/security/vault/unlock")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(json))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.message").value("Validation error"))
                .andExpect(jsonPath("$.errors[0].field").value("masterPassword"));
    }

    @Test
    void lockVault_shouldClearSessionAndReturnNoContent() throws Exception {
        masterKeySessionPort.setKey(new SecretKeySpec(new byte[32], "AES"));
        assertTrue(masterKeySessionPort.hasKey());

        mockMvc.perform(post("/api/v1/security/vault/lock"))
                .andExpect(status().isNoContent());

        assertFalse(masterKeySessionPort.hasKey());
    }
}
