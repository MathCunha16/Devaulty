package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.Credential;
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
class DeleteCredentialImplTest {

    @Mock
    private CredentialRepositoryPort credentialRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private MasterKeySessionPort masterKeySessionPort;

    @Mock
    private CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    @InjectMocks
    private DeleteCredentialImpl deleteCredentialUseCase;

    @Test
    void shouldDeleteCredentialSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        Credential credential = new Credential();
        credential.setId(credentialId);
        credential.setProjectId(projectId);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findById(credentialId)).thenReturn(Optional.of(credential));

        // Act
        deleteCredentialUseCase.execute(projectId, credentialId);

        // Assert
        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(credentialRepository, times(1)).deleteById(credentialId);
    }

    @Test
    void shouldThrowVaultLockedExceptionWhenKeyIsNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();

        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(masterKeySessionPort.getKey()).thenReturn(null);

        // Act & Assert
        assertThrows(VaultLockedException.class, () -> {
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findById(any());
        verify(credentialRepository, never()).deleteById(any());
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenSetupIsRequired() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();

        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, never()).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findById(any());
        verify(credentialRepository, never()).deleteById(any());
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
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, never()).findById(any());
        verify(credentialRepository, never()).deleteById(any());
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
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(credentialRepository, never()).deleteById(any());
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
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findById(credentialId);
        verify(credentialRepository, never()).deleteById(any());
    }
}
