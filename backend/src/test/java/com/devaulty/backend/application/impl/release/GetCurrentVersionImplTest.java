package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.infrastructure.properties.DevaultyProperties;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class GetCurrentVersionImplTest {

    @Mock
    private DevaultyProperties devaultyProperties;

    @InjectMocks
    private GetCurrentVersionImpl getCurrentVersionUseCase;

    @Test
    @DisplayName("Should return current application version string from DevaultyProperties")
    void shouldReturnCurrentAppVersion_whenExecuted() {
        // Arrange
        when(devaultyProperties.getVersion()).thenReturn("0.1.0-alpha");

        // Act
        String result = getCurrentVersionUseCase.execute();

        // Assert
        assertNotNull(result);
        assertEquals("0.1.0-alpha", result);
        verify(devaultyProperties, times(1)).getVersion();
    }
}
