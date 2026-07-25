package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import com.devaulty.backend.application.port.in.release.AppUpdateInfo;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.in.release.InstallUpdateUseCase;
import com.devaulty.backend.application.port.in.release.UpdateProgressInfo;
import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.core.io.buffer.DefaultDataBufferFactory;
import reactor.core.publisher.Flux;

import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DownloadUpdateImplTest {

    @Mock
    private ReleasePort releasePort;

    @Mock
    private CheckForUpdatesUseCase checkForUpdatesUseCase;

    @Mock
    private InstallUpdateUseCase installUpdateUseCase;

    @InjectMocks
    private DownloadUpdateImpl downloadUpdateUseCase;

    @Test
    @DisplayName("Should throw UpdateNotAvailableException when updateAvailable is false")
    void shouldThrowUpdateNotAvailableException_whenNoUpdateAvailable() {
        // Arrange
        AppUpdateInfo noUpdateInfo = new AppUpdateInfo(
                false,
                "0.1.0-alpha",
                "v0.1.0-alpha",
                "Title",
                "Notes",
                null,
                0L,
                Instant.now()
        );

        when(checkForUpdatesUseCase.execute()).thenReturn(noUpdateInfo);

        // Act & Assert
        assertThrows(UpdateNotAvailableException.class, () -> downloadUpdateUseCase.execute());
        verify(checkForUpdatesUseCase, times(1)).execute();
        verifyNoInteractions(releasePort);
    }

    @Test
    @DisplayName("Should throw UpdateNotAvailableException when downloadUrl is null")
    void shouldThrowUpdateNotAvailableException_whenDownloadUrlIsNull() {
        // Arrange
        AppUpdateInfo nullUrlInfo = new AppUpdateInfo(
                true,
                "0.1.0-alpha",
                "v0.2.0",
                "Title",
                "Notes",
                null,
                0L,
                Instant.now()
        );

        when(checkForUpdatesUseCase.execute()).thenReturn(nullUrlInfo);

        // Act & Assert
        assertThrows(UpdateNotAvailableException.class, () -> downloadUpdateUseCase.execute());
        verify(checkForUpdatesUseCase, times(1)).execute();
        verifyNoInteractions(releasePort);
    }

    @Test
    @DisplayName("Should stream download progress and trigger installation when download succeeds")
    void shouldStreamProgressAndTriggerInstallation_whenDownloadSucceeds() {
        // Arrange
        String downloadUrl = "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0.deb";
        long totalBytes = 11L;

        AppUpdateInfo validUpdateInfo = new AppUpdateInfo(
                true,
                "0.1.0-alpha",
                "v0.2.0",
                "Title",
                "Notes",
                downloadUrl,
                totalBytes,
                Instant.now()
        );

        when(checkForUpdatesUseCase.execute()).thenReturn(validUpdateInfo);

        DefaultDataBufferFactory factory = new DefaultDataBufferFactory();
        var buffer = factory.wrap("hello world".getBytes(StandardCharsets.UTF_8));
        when(releasePort.downloadAsset(downloadUrl)).thenReturn(Flux.just(buffer));
        doNothing().when(installUpdateUseCase).execute(any());

        // Act
        Flux<UpdateProgressInfo> progressFlux = downloadUpdateUseCase.execute();
        List<UpdateProgressInfo> results = progressFlux.collectList().block();

        // Assert
        assertNotNull(results);
        assertFalse(results.isEmpty());
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.DOWNLOADING));
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.INSTALLING));
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.COMPLETED && p.percentage() == 100));

        verify(checkForUpdatesUseCase, times(1)).execute();
        verify(releasePort, times(1)).downloadAsset(downloadUrl);
        verify(installUpdateUseCase, times(1)).execute(any());
    }
}
