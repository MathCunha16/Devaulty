package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ItemTagRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import javax.crypto.SecretKey;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
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

    @Mock
    private ItemTagRepositoryPort itemTagRepository;

    @InjectMocks
    private DeleteCredentialImpl deleteCredentialUseCase;

    @Test
    void shouldDeleteCredentialSuccessfully() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.existsByIdAndProjectId(credentialId, projectId)).thenReturn(true);

        // Act
        deleteCredentialUseCase.execute(projectId, credentialId);

        // Assert
        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).existsByIdAndProjectId(credentialId, projectId);
        verify(itemTagRepository, times(1)).removeAllTagsFromItem("credential", credentialId);
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
        verify(credentialRepository, never()).existsByIdAndProjectId(any(), any());
        verify(credentialRepository, never()).deleteById(any());
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
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
        verify(credentialRepository, never()).existsByIdAndProjectId(any(), any());
        verify(credentialRepository, never()).deleteById(any());
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
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
        verify(credentialRepository, never()).existsByIdAndProjectId(any(), any());
        verify(credentialRepository, never()).deleteById(any());
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
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
        when(credentialRepository.existsByIdAndProjectId(credentialId, projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            deleteCredentialUseCase.execute(projectId, credentialId);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).existsByIdAndProjectId(credentialId, projectId);
        verify(credentialRepository, never()).deleteById(any());
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenCredentialDoesNotBelongToProject() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        UUID credentialId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        // existsByIdAndProjectId returns false because it belongs to another project
        when(credentialRepository.existsByIdAndProjectId(credentialId, projectId)).thenReturn(false);

        // Act & Assert
        ResourceNotFoundException exception = assertThrows(ResourceNotFoundException.class, () -> {
            deleteCredentialUseCase.execute(projectId, credentialId);
        });
        assertEquals("Credential not found with identifier " + credentialId, exception.getMessage());

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).existsByIdAndProjectId(credentialId, projectId);
        verify(credentialRepository, never()).deleteById(any());
        verify(itemTagRepository, never()).removeAllTagsFromItem(any(), any());
    }
}
