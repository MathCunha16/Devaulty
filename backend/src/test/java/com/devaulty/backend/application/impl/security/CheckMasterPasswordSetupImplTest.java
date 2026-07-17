package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CheckMasterPasswordSetupImplTest {

    @Mock
    private AppSettingRepositoryPort appSettingRepositoryPort;

    @InjectMocks
    private CheckMasterPasswordSetupImpl checkMasterPasswordSetupUseCase;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";

    @Test
    void shouldReturnTrueWhenMasterPasswordHashDoesNotExist() {
        // Arrange
        when(appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY)).thenReturn(false);

        // Act
        boolean result = checkMasterPasswordSetupUseCase.isSetupRequired();

        // Assert
        assertTrue(result);
        verify(appSettingRepositoryPort, times(1)).existsByKey(MASTER_PASSWORD_HASH_KEY);
    }

    @Test
    void shouldReturnFalseWhenMasterPasswordHashExists() {
        // Arrange
        when(appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY)).thenReturn(true);

        // Act
        boolean result = checkMasterPasswordSetupUseCase.isSetupRequired();

        // Assert
        assertFalse(result);
        verify(appSettingRepositoryPort, times(1)).existsByKey(MASTER_PASSWORD_HASH_KEY);
    }
}
