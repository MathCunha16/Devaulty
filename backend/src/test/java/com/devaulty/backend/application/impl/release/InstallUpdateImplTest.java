package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.Spy;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.context.ConfigurableApplicationContext;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.anyList;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class InstallUpdateImplTest {

    @Mock
    private ConfigurableApplicationContext applicationContext;

    @Spy
    private InstallUpdateImpl installUpdateUseCase = new InstallUpdateImpl(applicationContext);

    @Test
    @DisplayName("Should throw UpdateNotAvailableException when installer file does not exist")
    void shouldThrowUpdateNotAvailableException_whenInstallerFileDoesNotExist() {
        // Arrange
        Path nonExistentPath = Path.of("non_existent_file.deb");

        // Act & Assert
        assertThrows(UpdateNotAvailableException.class, () -> installUpdateUseCase.execute(nonExistentPath));
    }

    @Test
    @DisplayName("Should throw IllegalStateException when installer file path is outside temp directory")
    void shouldThrowIllegalStateException_whenInstallerFileOutsideTempDir() throws IOException {
        // Arrange
        File outsideTempFile = File.createTempFile("outside_test", ".deb");
        outsideTempFile.deleteOnExit();

        // Act & Assert
        assertThrows(IllegalStateException.class, () -> installUpdateUseCase.execute(outsideTempFile.toPath()));
    }

    @Test
    @DisplayName("Should launch installer process without executing real system commands during unit tests")
    void shouldLaunchInstallerProcess_withoutExecutingRealSystemCommands() throws IOException {
        // Arrange
        String userHome = System.getProperty("user.home");
        String os = System.getProperty("os.name").toLowerCase();
        File tempDir;
        if (os.contains("win")) {
            String localAppData = System.getenv("LOCALAPPDATA");
            File base = localAppData != null ? new File(localAppData) : new File(userHome, "AppData/Local");
            tempDir = new File(base, "devaulty/temp");
        } else if (os.contains("mac")) {
            tempDir = new File(userHome, "Library/Caches/devaulty/temp");
        } else {
            tempDir = new File(userHome, ".config/devaulty/temp");
        }

        if (!tempDir.exists()) {
            tempDir.mkdirs();
        }

        File testInstaller = new File(tempDir, "test_installer.deb");
        testInstaller.createNewFile();
        testInstaller.deleteOnExit();

        // Prevent real pkexec / msiexec execution during unit test by mocking startDetached
        doNothing().when(installUpdateUseCase).startDetached(anyList());

        // Act & Assert
        assertDoesNotThrow(() -> installUpdateUseCase.execute(testInstaller.toPath()));
        verify(installUpdateUseCase, times(1)).startDetached(anyList());
    }
}
