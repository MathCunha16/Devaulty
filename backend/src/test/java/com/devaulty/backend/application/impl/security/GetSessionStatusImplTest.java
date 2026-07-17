package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.in.security.SessionStatus;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Duration;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetSessionStatusImplTest {

    @Mock
    private MasterKeySessionPort sessionHolder;

    @InjectMocks
    private GetSessionStatusImpl getSessionStatusUseCase;

    @Test
    void shouldReturnSessionStatusCorrectly() {
        // Arrange
        Duration timeout = Duration.ofMinutes(15);
        when(sessionHolder.hasKey()).thenReturn(true);
        when(sessionHolder.getSecondsRemaining(timeout)).thenReturn(900L);

        // Act
        SessionStatus status = getSessionStatusUseCase.execute();

        // Assert
        assertNotNull(status);
        assertTrue(status.active());
        assertEquals(900L, status.secondsLeft());

        verify(sessionHolder, times(1)).hasKey();
        verify(sessionHolder, times(1)).getSecondsRemaining(timeout);
    }
}
