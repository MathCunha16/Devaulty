package com.devaulty.backend.adapter.in.scheduler;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;

import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class VaultAutoLockSchedulerTest {

    @Mock
    private MasterKeySessionPort masterKeySessionPort;

    @InjectMocks
    private VaultAutoLockScheduler scheduler;

    @Test
    void purgeExpiredSession_shouldDoNothing_whenSessionHasNoKey() {
        // Arrange
        when(masterKeySessionPort.hasKey()).thenReturn(false);

        // Act
        scheduler.purgeExpiredSession();

        // Assert
        verify(masterKeySessionPort, times(1)).hasKey();
        verify(masterKeySessionPort, never()).getSecondsRemaining(any());
        verify(masterKeySessionPort, never()).clear();
    }

    @Test
    void purgeExpiredSession_shouldDoNothing_whenSessionIsValid() {
        // Arrange
        when(masterKeySessionPort.hasKey()).thenReturn(true);
        when(masterKeySessionPort.getSecondsRemaining(any(Duration.class))).thenReturn(100L);

        // Act
        scheduler.purgeExpiredSession();

        // Assert
        verify(masterKeySessionPort, times(1)).hasKey();
        verify(masterKeySessionPort, times(1)).getSecondsRemaining(Duration.ofMinutes(15));
        verify(masterKeySessionPort, never()).clear();
    }

    @Test
    void purgeExpiredSession_shouldClearSession_whenSessionIsExpired() {
        // Arrange
        when(masterKeySessionPort.hasKey()).thenReturn(true);
        when(masterKeySessionPort.getSecondsRemaining(any(Duration.class))).thenReturn(0L);

        // Act
        scheduler.purgeExpiredSession();

        // Assert
        verify(masterKeySessionPort, times(1)).hasKey();
        verify(masterKeySessionPort, times(1)).getSecondsRemaining(Duration.ofMinutes(15));
        verify(masterKeySessionPort, times(1)).clear();
    }

    @Test
    void purgeExpiredSession_shouldClearSession_whenSessionIsExpiredNegative() {
        // Arrange
        when(masterKeySessionPort.hasKey()).thenReturn(true);
        when(masterKeySessionPort.getSecondsRemaining(any(Duration.class))).thenReturn(-5L);

        // Act
        scheduler.purgeExpiredSession();

        // Assert
        verify(masterKeySessionPort, times(1)).hasKey();
        verify(masterKeySessionPort, times(1)).getSecondsRemaining(Duration.ofMinutes(15));
        verify(masterKeySessionPort, times(1)).clear();
    }
}
