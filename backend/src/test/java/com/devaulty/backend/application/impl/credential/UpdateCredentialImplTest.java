package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.DecryptedCredential;
import com.devaulty.backend.application.port.in.credential.UpdateCredentialCommand;
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
import java.util.Optional;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UpdateCredentialImplTest {

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
    private UpdateCredentialImpl updateCredentialUseCase;

    @Test
    void shouldUpdateCredentialSuccessfully_includingPayload() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        char[] newPayload = "newPayload".toCharArray();
        UpdateCredentialCommand command = new UpdateCredentialCommand(
                credentialId,
                projectId,
                "New Title",
                CredentialSecretType.API_KEY,
                newPayload,
                "New Notes",
                "http://new"
        );

        byte[] oldEncryptedPayload = new byte[]{1, 2};
        byte[] oldIv = new byte[]{3};
        byte[] oldAuthTag = new byte[]{4};

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(projectId);
        credential.setTitle("Old Title");
        credential.setSecretType(CredentialSecretType.LOGIN);
        credential.setPayloadEncrypted(oldEncryptedPayload);
        credential.setEncryptionIv(oldIv);
        credential.setEncryptionAuthTag(oldAuthTag);

        byte[] newEncrypted = new byte[]{5, 6};
        byte[] newIv = new byte[]{7};
        byte[] newAuthTag = new byte[]{8};
        byte[] decryptedPayload = "newPayload".getBytes();

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));
        when(cryptoPort.encrypt(any(byte[].class), eq(mockKey), any(byte[].class)))
                .thenReturn(new CryptoResultDto(newEncrypted, newIv, newAuthTag));
        when(credentialRepository.save(any(Credential.class))).thenAnswer(invocation -> invocation.getArgument(0));
        when(cryptoPort.decrypt(eq(newEncrypted), eq(newIv), eq(newAuthTag), eq(mockKey), any(byte[].class))).thenReturn(decryptedPayload);

        // Act
        DecryptedCredential result = updateCredentialUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("New Title", result.title());
        assertEquals(CredentialSecretType.API_KEY, result.secretType());
        assertEquals("New Notes", result.notes());
        assertEquals("http://new", result.relatedUrl());
        assertArrayEquals(decryptedPayload, result.decryptedPayload());
        assertNotNull(result.updatedAt());

        // Payload zeroed out in command
        for (char c : command.payload()) {
            assertEquals('\0', c);
        }

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(cryptoPort, times(1)).encrypt(any(byte[].class), eq(mockKey), any(byte[].class));
        verify(credentialRepository, times(1)).save(credential);
        verify(cryptoPort, times(1)).decrypt(eq(newEncrypted), eq(newIv), eq(newAuthTag), eq(mockKey), any(byte[].class));
    }

    @Test
    void shouldUpdateCredentialSuccessfully_excludingPayload() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        UpdateCredentialCommand command = new UpdateCredentialCommand(
                credentialId,
                projectId,
                null,
                null,
                null,
                null,
                null
        );

        byte[] oldEncryptedPayload = new byte[]{1, 2};
        byte[] oldIv = new byte[]{3};
        byte[] oldAuthTag = new byte[]{4};
        byte[] decryptedPayload = "oldPayload".getBytes();

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(projectId);
        credential.setTitle("Old Title");
        credential.setSecretType(CredentialSecretType.LOGIN);
        credential.setPayloadEncrypted(oldEncryptedPayload);
        credential.setEncryptionIv(oldIv);
        credential.setEncryptionAuthTag(oldAuthTag);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));
        when(credentialRepository.save(any(Credential.class))).thenAnswer(invocation -> invocation.getArgument(0));
        when(cryptoPort.decrypt(eq(oldEncryptedPayload), eq(oldIv), eq(oldAuthTag), eq(mockKey), any(byte[].class))).thenReturn(decryptedPayload);

        // Act
        DecryptedCredential result = updateCredentialUseCase.execute(command);

        // Assert
        assertNotNull(result);
        assertEquals("Old Title", result.title());
        assertEquals(CredentialSecretType.LOGIN, result.secretType());
        assertArrayEquals(decryptedPayload, result.decryptedPayload());
        assertNotNull(result.updatedAt());

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(cryptoPort, never()).encrypt(any(), any(), any());
        verify(credentialRepository, times(1)).save(credential);
        verify(cryptoPort, times(1)).decrypt(eq(oldEncryptedPayload), eq(oldIv), eq(oldAuthTag), eq(mockKey), any(byte[].class));
    }

    @Test
    void shouldThrowVaultLockedExceptionWhenKeyIsNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        UpdateCredentialCommand command = new UpdateCredentialCommand(credentialId, projectId, null, null, null, null, null);

        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(masterKeySessionPort.getKey()).thenReturn(null);

        // Act & Assert
        assertThrows(VaultLockedException.class, () -> {
            updateCredentialUseCase.execute(command);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findById(any());
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenSetupIsRequired() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);
        UpdateCredentialCommand command = new UpdateCredentialCommand(credentialId, projectId, null, null, null, null, null);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            updateCredentialUseCase.execute(command);
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
        UpdateCredentialCommand command = new UpdateCredentialCommand(credentialId, projectId, null, null, null, null, null);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateCredentialUseCase.execute(command);
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
        UpdateCredentialCommand command = new UpdateCredentialCommand(credentialId, projectId, null, null, null, null, null);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.empty());

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateCredentialUseCase.execute(command);
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
        UpdateCredentialCommand command = new UpdateCredentialCommand(credentialId, projectId, null, null, null, null, null);

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(otherProjectId);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            updateCredentialUseCase.execute(command);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(credentialRepository, never()).save(any());
    }
}
