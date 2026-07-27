package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import com.devaulty.backend.application.port.in.release.AppUpdateInfo;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.in.release.InstallUpdateUseCase;
import com.devaulty.backend.application.port.in.release.UpdateProgressInfo;
import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class DownloadUpdateImplTest {

    @TempDir
    Path tempFolder;

    @Mock
    private ReleasePort releasePort;

    @Mock
    private CheckForUpdatesUseCase checkForUpdatesUseCase;

    @Mock
    private InstallUpdateUseCase installUpdateUseCase;

    @InjectMocks
    private DownloadUpdateImpl downloadUpdateUseCase;

    @BeforeEach
    void setUp() {
        System.setProperty("devaulty.temp.dir", tempFolder.toString());
    }

    @AfterEach
    void tearDown() {
        System.clearProperty("devaulty.temp.dir");
    }

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
        assertThrows(UpdateNotAvailableException.class, () -> downloadUpdateUseCase.execute(progress -> {}));
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
        assertThrows(UpdateNotAvailableException.class, () -> downloadUpdateUseCase.execute(progress -> {}));
        verify(checkForUpdatesUseCase, times(1)).execute();
        verifyNoInteractions(releasePort);
    }

    @Test
    @DisplayName("Should emit download progress and trigger installation when download succeeds")
    void shouldEmitProgressAndTriggerInstallation_whenDownloadSucceeds() {
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

        InputStream fakeAssetStream = new ByteArrayInputStream("hello world".getBytes(StandardCharsets.UTF_8));
        when(releasePort.downloadAsset(downloadUrl)).thenReturn(fakeAssetStream);
        doNothing().when(installUpdateUseCase).execute(any());

        // Act
        List<UpdateProgressInfo> results = new ArrayList<>();
        downloadUpdateUseCase.execute(results::add);

        // Assert
        assertFalse(results.isEmpty());
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.DOWNLOADING));
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.INSTALLING));
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.COMPLETED && p.percentage() == 100));

        verify(checkForUpdatesUseCase, times(1)).execute();
        verify(releasePort, times(1)).downloadAsset(downloadUrl);
        verify(installUpdateUseCase, times(1)).execute(any());
    }

    @Test
    @DisplayName("Should emit FAILED progress and clean up file when download stream throws")
    void shouldEmitFailedProgress_whenDownloadStreamThrows() {
        // Arrange
        String downloadUrl = "https://github.com/MathCunha16/Devaulty/releases/download/v0.2.0/devaulty_0.2.0.deb";

        AppUpdateInfo validUpdateInfo = new AppUpdateInfo(
                true,
                "0.1.0-alpha",
                "v0.2.0",
                "Title",
                "Notes",
                downloadUrl,
                11L,
                Instant.now()
        );

        when(checkForUpdatesUseCase.execute()).thenReturn(validUpdateInfo);

        InputStream brokenStream = new InputStream() {
            @Override
            public int read() throws java.io.IOException {
                throw new java.io.IOException("Connection reset");
            }
        };
        when(releasePort.downloadAsset(downloadUrl)).thenReturn(brokenStream);

        // Act
        List<UpdateProgressInfo> results = new ArrayList<>();
        downloadUpdateUseCase.execute(results::add);

        // Assert
        assertTrue(results.stream().anyMatch(p -> p.status() == UpdateStatus.FAILED));
        verify(installUpdateUseCase, never()).execute(any());
    }
}