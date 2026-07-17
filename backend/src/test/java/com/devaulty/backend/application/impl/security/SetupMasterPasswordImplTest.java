package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.exception.MasterPasswordAlreadyConfiguredException;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.AppSetting;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import javax.crypto.SecretKey;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class SetupMasterPasswordImplTest {

    @Mock
    private AppSettingRepositoryPort appSettingRepositoryPort;

    @Mock
    private MasterKeySessionPort sessionHolder;

    @Mock
    private KeyDerivationPort keyDerivationPort;

    @InjectMocks
    private SetupMasterPasswordImpl setupMasterPasswordUseCase;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";

    @Test
    void shouldSetupMasterPasswordSuccessfully() {
        // Arrange
        char[] password = "mySuperSecretMasterPassword".toCharArray();
        SecretKey mockSecretKey = mock(SecretKey.class);

        when(appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY)).thenReturn(false);
        when(keyDerivationPort.deriveKey(eq(password), any(byte[].class))).thenReturn(mockSecretKey);

        // Act
        setupMasterPasswordUseCase.execute(password);

        // Assert
        verify(appSettingRepositoryPort, times(1)).existsByKey(MASTER_PASSWORD_HASH_KEY);
        verify(keyDerivationPort, times(1)).deriveKey(eq(password), any(byte[].class));
        verify(appSettingRepositoryPort, times(2)).save(any(AppSetting.class)); // Salt and Hash
        verify(sessionHolder, times(1)).setKey(mockSecretKey);

        // Input password should be zeroed out in the finally block
        assertEquals('\0', password[0]);
    }

    @Test
    void shouldThrowMasterPasswordAlreadyConfiguredExceptionWhenAlreadyConfigured() {
        // Arrange
        char[] password = "mySuperSecretMasterPassword".toCharArray();

        when(appSettingRepositoryPort.existsByKey(MASTER_PASSWORD_HASH_KEY)).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordAlreadyConfiguredException.class, () -> {
            setupMasterPasswordUseCase.execute(password);
        });

        verify(appSettingRepositoryPort, times(1)).existsByKey(MASTER_PASSWORD_HASH_KEY);
        verify(keyDerivationPort, never()).deriveKey(any(), any());
        verify(appSettingRepositoryPort, never()).save(any());
        verify(sessionHolder, never()).setKey(any());

        // Password is zeroed out successfully because validations are now inside the try-finally block
        assertEquals('\0', password[0]);
    }
}
