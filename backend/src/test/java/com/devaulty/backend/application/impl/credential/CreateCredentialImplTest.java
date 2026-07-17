package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.CreateCredentialCommand;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;
import com.devaulty.backend.domain.model.Credential;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import javax.crypto.SecretKey;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CreateCredentialImplTest {

    @Mock
    private CredentialRepositoryPort credentialRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private CryptoPort cryptoPort;

    @Mock
    private MasterKeySessionPort masterKeySessionPort;

    @Mock
    private CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    @InjectMocks
    private CreateCredentialImpl createCredentialUseCase;

    @Test
    void shouldCreateCredentialSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        char[] passwordPayload = "mySecretPassword".toCharArray();
        CreateCredentialCommand command = new CreateCredentialCommand(
                projectId,
                "DB Pass",
                CredentialSecretType.LOGIN,
                passwordPayload,
                "notes",
                "http://db"
        );

        SecretKey mockSecretKey = mock(SecretKey.class);
        byte[] mockEncryptedPayload = new byte[]{1, 2, 3};
        byte[] mockIv = new byte[]{4, 5};
        byte[] mockAuthTag = new byte[]{6, 7};
        byte[] mockDecryptedPayload = "mySecretPassword".getBytes();

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(masterKeySessionPort.getKey()).thenReturn(mockSecretKey);
        
        when(cryptoPort.encrypt(any(byte[].class), eq(mockSecretKey)))
                .thenReturn(new CryptoResultDto(mockEncryptedPayload, mockIv, mockAuthTag));

        when(credentialRepository.save(any(Credential.class))).thenAnswer(invocation -> {
            Credential c = invocation.getArgument(0);
            c.setCreatedAt(java.time.LocalDateTime.now());
            return c;
        });

        when(cryptoPort.decrypt(mockEncryptedPayload, mockIv, mockAuthTag, mockSecretKey))
                .thenReturn(mockDecryptedPayload);

        // Act
        DecryptedCredential result = createCredentialUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertNotNull(result.id());
        assertEquals(projectId, result.projectId());
        assertEquals("DB Pass", result.title());
        assertEquals(CredentialSecretType.LOGIN, result.secretType());
        assertArrayEquals(mockDecryptedPayload, result.decryptedPayload());
        assertEquals("notes", result.notes());
        assertEquals("http://db", result.relatedUrl());
        assertNotNull(result.createdAt());

        // The method should zero out the command payload
        assertEquals('\0', command.payload()[0]);

        verify(projectRepository, times(1)).existsById(projectId);
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(masterKeySessionPort, times(1)).getKey();
        verify(cryptoPort, times(1)).encrypt(any(byte[].class), eq(mockSecretKey));
        verify(credentialRepository, times(1)).save(any(Credential.class));
        verify(cryptoPort, times(1)).decrypt(mockEncryptedPayload, mockIv, mockAuthTag, mockSecretKey);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        char[] passwordPayload = "mySecretPassword".toCharArray();
        CreateCredentialCommand command = new CreateCredentialCommand(
                projectId,
                "DB Pass",
                CredentialSecretType.LOGIN,
                passwordPayload,
                "notes",
                "http://db"
        );

        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            createCredentialUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, never()).save(any());
        // Clean up: payload is zeroed here successfully because validations are now inside the try-finally block
        assertEquals('\0', command.payload()[0]);
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenSetupIsRequired() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        char[] passwordPayload = "mySecretPassword".toCharArray();
        CreateCredentialCommand command = new CreateCredentialCommand(
                projectId,
                "DB Pass",
                CredentialSecretType.LOGIN,
                passwordPayload,
                "notes",
                "http://db"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            createCredentialUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(credentialRepository, never()).save(any());
        // Clean up: payload is zeroed here successfully because validations are now inside the try-finally block
        assertEquals('\0', command.payload()[0]);
    }

    @Test
    void shouldThrowVaultLockedExceptionWhenMasterKeyIsMissing() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        char[] passwordPayload = "mySecretPassword".toCharArray();
        CreateCredentialCommand command = new CreateCredentialCommand(
                projectId,
                "DB Pass",
                CredentialSecretType.LOGIN,
                passwordPayload,
                "notes",
                "http://db"
        );

        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(masterKeySessionPort.getKey()).thenReturn(null);

        // Act & Assert
        assertThrows(VaultLockedException.class, () -> {
            createCredentialUseCase.execute(command);
        });

        verify(projectRepository, times(1)).existsById(projectId);
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(masterKeySessionPort, times(1)).getKey();
        verify(credentialRepository, never()).save(any());
        assertEquals('\0', command.payload()[0]);
    }
}
