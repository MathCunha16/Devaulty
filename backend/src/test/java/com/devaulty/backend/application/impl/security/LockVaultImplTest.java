package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class LockVaultImplTest {

    @Mock
    private MasterKeySessionPort sessionHolder;

    @InjectMocks
    private LockVaultImpl lockVaultUseCase;

    @Test
    void shouldClearSessionWhenExecutingLock() {
        // Act
        lockVaultUseCase.execute();

        // Assert
        verify(sessionHolder, times(1)).clear();
    }
}
