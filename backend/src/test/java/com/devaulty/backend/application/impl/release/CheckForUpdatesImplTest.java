package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.port.in.release.AppUpdateInfo;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import com.devaulty.backend.application.port.out.external.release.dto.ReleaseAssetInfo;
import com.devaulty.backend.infrastructure.properties.DevaultyProperties;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.time.Instant;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class CheckForUpdatesImplTest {

    @Mock
    private ReleasePort releasePort;

    @Mock
    private DevaultyProperties devaultyProperties;

    @InjectMocks
    private CheckForUpdatesImpl checkForUpdatesUseCase;

    @BeforeEach
    void setUp() {
        lenient().when(devaultyProperties.getVersion()).thenReturn("0.1.0-alpha");
    }

    @Test
    @DisplayName("Should return updateAvailable = true when GitHub has a newer version tag")
    void shouldReturnUpdateAvailable_whenNewerReleaseExistsOnGitHub() {
        // Arrange
        String currentOs = System.getProperty("os.name").toLowerCase();
        String extension = currentOs.contains("win") ? ".msi" : currentOs.contains("mac") ? ".dmg" : ".deb";

        ReleaseAssetInfo matchingAsset = new ReleaseAssetInfo(
                "devaulty-installer" + extension,
                "https://download.url/devaulty" + extension,
                50000000L,
                "application/octet-stream"
        );
        LatestReleaseInfo latestRelease = new LatestReleaseInfo(
                "v0.2.0",
                "Release 0.2.0",
                "Bug fixes and performance improvements",
                "https://github.com/MathCunha16/Devaulty/releases/tag/v0.2.0",
                Instant.now(),
                false,
                List.of(matchingAsset)
        );

        when(releasePort.getLatestRelease()).thenReturn(latestRelease);

        // Act
        AppUpdateInfo result = checkForUpdatesUseCase.execute();

        // Assert
        assertNotNull(result);
        assertTrue(result.updateAvailable());
        assertEquals("0.1.0-alpha", result.currentVersion());
        assertEquals("v0.2.0", result.latestVersion());
        assertEquals("Release 0.2.0", result.releaseTitle());
        assertEquals("Bug fixes and performance improvements", result.releaseNotes());
        assertEquals("https://download.url/devaulty" + extension, result.downloadUrl());
        assertEquals(50000000L, result.downloadSizeInBytes());

        verify(releasePort, times(1)).getLatestRelease();
    }

    @Test
    @DisplayName("Should return updateAvailable = false when current version is equal to GitHub release tag")
    void shouldReturnNoUpdateAvailable_whenCurrentVersionIsUpToDate() {
        // Arrange
        LatestReleaseInfo latestRelease = new LatestReleaseInfo(
                "v0.1.0-alpha",
                "Release 0.1.0-alpha",
                "Initial alpha release",
                "https://github.com/MathCunha16/Devaulty/releases/tag/v0.1.0-alpha",
                Instant.now(),
                false,
                List.of()
        );

        when(releasePort.getLatestRelease()).thenReturn(latestRelease);

        // Act
        AppUpdateInfo result = checkForUpdatesUseCase.execute();

        // Assert
        assertNotNull(result);
        assertFalse(result.updateAvailable());
        assertEquals("0.1.0-alpha", result.currentVersion());
        assertEquals("v0.1.0-alpha", result.latestVersion());

        verify(releasePort, times(1)).getLatestRelease();
    }

    @Test
    @DisplayName("Should return updateAvailable = false when GitHub release is null")
    void shouldReturnNoUpdateAvailable_whenNoReleaseExistsOnGitHub() {
        // Arrange
        when(releasePort.getLatestRelease()).thenReturn(null);

        // Act
        AppUpdateInfo result = checkForUpdatesUseCase.execute();

        // Assert
        assertNotNull(result);
        assertFalse(result.updateAvailable());
        assertEquals("0.1.0-alpha", result.currentVersion());
        assertNull(result.latestVersion());
        assertNull(result.downloadUrl());
        assertEquals(0L, result.downloadSizeInBytes());

        verify(releasePort, times(1)).getLatestRelease();
    }
}
