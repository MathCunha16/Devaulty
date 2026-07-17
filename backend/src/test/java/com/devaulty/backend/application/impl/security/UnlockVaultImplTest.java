package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.exception.InvalidMasterPasswordException;
import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.port.out.persistence.AppSettingRepositoryPort;
import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.AppSetting;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.argon2.Argon2PasswordEncoder;
import org.springframework.test.util.ReflectionTestUtils;

import javax.crypto.SecretKey;
import java.nio.CharBuffer;
import java.util.Base64;
import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UnlockVaultImplTest {

    @Mock
    private AppSettingRepositoryPort appSettingRepositoryPort;

    @Mock
    private MasterKeySessionPort sessionHolder;

    @Mock
    private KeyDerivationPort keyDerivationPort;

    @InjectMocks
    private UnlockVaultImpl unlockVaultUseCase;

    private static final String MASTER_PASSWORD_HASH_KEY = "master_password_hash";
    private static final String MASTER_PASSWORD_SALT_KEY = "master_password_salt";

    @Test
    void shouldUnlockVaultSuccessfully() {
        // Arrange
        char[] password = "myMasterPassword123".toCharArray();
        String mockHash = "mockHashEncodedWithArgon2";
        byte[] saltBytes = new byte[]{1, 2, 3};
        String mockSaltBase64 = Base64.getEncoder().encodeToString(saltBytes);
        SecretKey mockSecretKey = mock(SecretKey.class);

        // Mock Argon2PasswordEncoder inside UnlockVaultImpl using reflection
        Argon2PasswordEncoder mockEncoder = mock(Argon2PasswordEncoder.class);
        ReflectionTestUtils.setField(unlockVaultUseCase, "argon2PasswordEncoder", mockEncoder);

        when(appSettingRepositoryPort.findByKey(MASTER_PASSWORD_HASH_KEY))
                .thenReturn(Optional.of(new AppSetting(MASTER_PASSWORD_HASH_KEY, mockHash)));
        when(mockEncoder.matches(any(CharBuffer.class), eq(mockHash))).thenReturn(true);
        when(appSettingRepositoryPort.findByKey(MASTER_PASSWORD_SALT_KEY))
                .thenReturn(Optional.of(new AppSetting(MASTER_PASSWORD_SALT_KEY, mockSaltBase64)));
        when(keyDerivationPort.deriveKey(eq(password), any(byte[].class))).thenReturn(mockSecretKey);

        // Act
        boolean result = unlockVaultUseCase.execute(password);

        // Assert
        assertTrue(result);
        verify(appSettingRepositoryPort, times(1)).findByKey(MASTER_PASSWORD_HASH_KEY);
        verify(mockEncoder, times(1)).matches(any(CharBuffer.class), eq(mockHash));
        verify(appSettingRepositoryPort, times(1)).findByKey(MASTER_PASSWORD_SALT_KEY);
        verify(keyDerivationPort, times(1)).deriveKey(eq(password), any(byte[].class));
        verify(sessionHolder, times(1)).setKey(mockSecretKey);

        // Password is zeroed out in finally
        assertEquals('\0', password[0]);
    }

    @Test
    void shouldThrowInvalidMasterPasswordExceptionWhenPasswordDoesNotMatch() {
        // Arrange
        char[] password = "wrongMasterPassword".toCharArray();
        String mockHash = "mockHashEncodedWithArgon2";

        Argon2PasswordEncoder mockEncoder = mock(Argon2PasswordEncoder.class);
        ReflectionTestUtils.setField(unlockVaultUseCase, "argon2PasswordEncoder", mockEncoder);

        when(appSettingRepositoryPort.findByKey(MASTER_PASSWORD_HASH_KEY))
                .thenReturn(Optional.of(new AppSetting(MASTER_PASSWORD_HASH_KEY, mockHash)));
        when(mockEncoder.matches(any(CharBuffer.class), eq(mockHash))).thenReturn(false);

        // Act & Assert
        assertThrows(InvalidMasterPasswordException.class, () -> {
            unlockVaultUseCase.execute(password);
        });

        verify(appSettingRepositoryPort, times(1)).findByKey(MASTER_PASSWORD_HASH_KEY);
        verify(mockEncoder, times(1)).matches(any(CharBuffer.class), eq(mockHash));
        verify(appSettingRepositoryPort, never()).findByKey(MASTER_PASSWORD_SALT_KEY);
        verify(keyDerivationPort, never()).deriveKey(any(), any());
        verify(sessionHolder, never()).setKey(any());

        // Password is zeroed out
        assertEquals('\0', password[0]);
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenHashIsMissing() {
        // Arrange
        char[] password = "myMasterPassword123".toCharArray();

        when(appSettingRepositoryPort.findByKey(MASTER_PASSWORD_HASH_KEY)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            unlockVaultUseCase.execute(password);
        });

        verify(appSettingRepositoryPort, times(1)).findByKey(MASTER_PASSWORD_HASH_KEY);
        verify(appSettingRepositoryPort, never()).findByKey(MASTER_PASSWORD_SALT_KEY);
        verify(keyDerivationPort, never()).deriveKey(any(), any());
        verify(sessionHolder, never()).setKey(any());

        // Password is zeroed out
        assertEquals('\0', password[0]);
    }
}
