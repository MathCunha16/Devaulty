package com.devaulty.backend.application.impl.credential;

import com.devaulty.backend.application.exception.MasterPasswordNotConfiguredException;
import com.devaulty.backend.application.exception.ResourceNotFoundException;
import com.devaulty.backend.application.exception.VaultLockedException;
import com.devaulty.backend.application.port.in.credential.CredentialSummary;
import com.devaulty.backend.application.port.in.security.CheckMasterPasswordSetupUseCase;
import com.devaulty.backend.application.port.out.persistence.CredentialRepositoryPort;
import com.devaulty.backend.application.port.out.persistence.ProjectRepositoryPort;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import com.devaulty.backend.domain.model.enums.CredentialSecretType;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;

import javax.crypto.SecretKey;
import java.util.Collections;
import java.util.List;
import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetAllCredentialsByProjectImplTest {

    @Mock
    private CredentialRepositoryPort credentialRepository;

    @Mock
    private ProjectRepositoryPort projectRepository;

    @Mock
    private MasterKeySessionPort masterKeySessionPort;

    @Mock
    private CheckMasterPasswordSetupUseCase checkMasterPasswordSetupUseCase;

    @InjectMocks
    private GetAllCredentialsByProjectImpl getAllCredentialsUseCase;

    @Test
    void shouldReturnPageOfCredentials() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        int page = 0;
        int size = 10;
        SecretKey mockKey = mock(SecretKey.class);

        CredentialSummary summary = new CredentialSummary(
                UUID.randomUUID(),
                projectId,
                "Title",
                CredentialSecretType.LOGIN,
                "http://url",
                java.time.LocalDateTime.now(),
                null
        );
        List<CredentialSummary> list = Collections.singletonList(summary);
        Page<CredentialSummary> expectedPage = new PageImpl<>(list);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(true);
        when(credentialRepository.findAllByProject(projectId, page, size)).thenReturn(expectedPage);

        // Act
        Page<CredentialSummary> result = getAllCredentialsUseCase.execute(projectId, page, size);

        // Assert
        assertNotNull(result);
        assertEquals(expectedPage, result);
        assertEquals(1, result.getTotalElements());
        assertEquals(summary, result.getContent().getFirst());

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, times(1)).findAllByProject(projectId, page, size);
    }

    @Test
    void shouldThrowVaultLockedExceptionWhenKeyIsNull() {
        // Arrange
        UUID projectId = UUID.randomUUID();

        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(masterKeySessionPort.getKey()).thenReturn(null);

        // Act & Assert
        assertThrows(VaultLockedException.class, () -> {
            getAllCredentialsUseCase.execute(projectId, 0, 10);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findAllByProject(any(), anyInt(), anyInt());
    }

    @Test
    void shouldThrowMasterPasswordNotConfiguredExceptionWhenSetupIsRequired() {
        // Arrange
        UUID projectId = UUID.randomUUID();

        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(true);

        // Act & Assert
        assertThrows(MasterPasswordNotConfiguredException.class, () -> {
            getAllCredentialsUseCase.execute(projectId, 0, 10);
        });

        verify(masterKeySessionPort, never()).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, never()).existsById(any());
        verify(credentialRepository, never()).findAllByProject(any(), anyInt(), anyInt());
    }

    @Test
    void shouldThrowResourceNotFoundExceptionWhenProjectDoesNotExist() {
        // Arrange
        UUID projectId = UUID.randomUUID();
        SecretKey mockKey = mock(SecretKey.class);

        when(masterKeySessionPort.getKey()).thenReturn(mockKey);
        when(checkMasterPasswordSetupUseCase.isSetupRequired()).thenReturn(false);
        when(projectRepository.existsById(projectId)).thenReturn(false);

        // Act & Assert
        assertThrows(ResourceNotFoundException.class, () -> {
            getAllCredentialsUseCase.execute(projectId, 0, 10);
        });

        verify(masterKeySessionPort, times(1)).getKey();
        verify(checkMasterPasswordSetupUseCase, times(1)).isSetupRequired();
        verify(projectRepository, times(1)).existsById(projectId);
        verify(credentialRepository, never()).findAllByProject(any(), anyInt(), anyInt());
    }
}
