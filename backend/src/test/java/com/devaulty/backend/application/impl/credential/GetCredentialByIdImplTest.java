package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.Credential;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import javax.crypto.SecretKey;
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetCredentialByIdImplTest {

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
    private GetCredentialByIdImpl getCredentialByIdUseCase;

    @Test
    void shouldReturnCredentialWhenFoundAndBelongsToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        byte[] encryptedPayload = new byte[]{1, 2, 3};
        byte[] iv = new byte[]{4};
        byte[] authTag = new byte[]{5};
        byte[] decryptedPayload = "decrypted".getBytes();

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(projectId);
        credential.setTitle("Title");
        credential.setSecretType(CredentialSecretType.LOGIN);
        credential.setPayloadEncrypted(encryptedPayload);
        credential.setEncryptionIv(iv);
        credential.setEncryptionAuthTag(authTag);
        credential.setCreatedAt(java.time.LocalDateTime.now());

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));
        when(cryptoPort.decrypt(encryptedPayload, iv, authTag, mockKey)).thenReturn(decryptedPayload);

        // Act
        DecryptedCredential result = getCredentialByIdUseCase.execute(projectId, credentialId);

        // Assert
        assertNotNull(result);
        assertEquals(credentialId, result.id());
        assertEquals(projectId, result.projectId());
        assertEquals("Title", result.title());
        assertEquals(CredentialSecretType.LOGIN, result.secretType());
        assertArrayEquals(decryptedPayload, result.decryptedPayload());

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(cryptoPort, times(1)).decrypt(encryptedPayload, iv, authTag, mockKey);
    }

    @Test
    void shouldThrowVaultLockedExceptionWhenKeyIsNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();

        when(masterKeySessionPort.getKey()).thenReturn(null);

        // Act & Assert
        assertThrows(VaultLockedException.class, () -> {
            getCredentialByIdUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, never()).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findById(any());
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenSetupIsRequired() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            getCredentialByIdUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findById(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getCredentialByIdUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, never()).findById(any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenCredentialDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getCredentialByIdUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenCredentialDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID otherProjectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(otherProjectId);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getCredentialByIdUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(cryptoPort, never()).decrypt(any(), any(), any(), any());
    }
}
