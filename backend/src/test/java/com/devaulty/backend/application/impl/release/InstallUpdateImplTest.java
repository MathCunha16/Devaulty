package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
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

    @TempDir
    Path tempFolder;

    @Mock
    private ConfigurableApplicationContext applicationContext;

    @Spy
    private InstallUpdateImpl installUpdateUseCase = new InstallUpdateImpl(applicationContext);

    @BeforeEach
    void setUp() {
        System.setProperty("devaulty.temp.dir", tempFolder.toString());
    }

    @AfterEach
    void tearDown() {
        System.clearProperty("devaulty.temp.dir");
    }

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
        File tempDir = tempFolder.toFile();
        File testInstaller = new File(tempDir, "test_installer.deb");
        testInstaller.createNewFile();

        // Prevent real pkexec / msiexec execution during unit test by mocking startDetached
        doNothing().when(installUpdateUseCase).startDetached(anyList());

        // Act & Assert
        assertDoesNotThrow(() -> installUpdateUseCase.execute(testInstaller.toPath()));
        verify(installUpdateUseCase, times(1)).startDetached(anyList());
    }
}
