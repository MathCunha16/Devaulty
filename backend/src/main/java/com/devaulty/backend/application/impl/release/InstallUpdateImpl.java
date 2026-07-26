package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import com.devaulty.backend.application.port.in.release.InstallUpdateUseCase;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.SpringApplication;
import org.springframework.context.ConfigurableApplicationContext;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

/**
 * Launches the platform-specific installer for a previously downloaded update
 * and schedules the current application instance to shut down so the
 * installer can replace the running files.
 *
 * <p>Known limitation: on macOS the update flow only mounts the {@code .dmg}
 * image; it does not copy the app into {@code /Applications} or relaunch it
 * automatically. A full macOS auto-update typically requires a dedicated
 * mechanism (e.g. Sparkle) and is left as a follow-up.</p>
 */
public class InstallUpdateImpl implements InstallUpdateUseCase {

    private static final Logger logger = LoggerFactory.getLogger(InstallUpdateImpl.class);

    private static final long SHUTDOWN_DELAY_SECONDS = 5;
    private static final String WINDOWS_INSTALL_DIR_ENV = "DEVAULTY_INSTALL_DIR";
    private static final String WINDOWS_DEFAULT_INSTALL_DIR = "C:\\Program Files\\devaulty";

    private final ConfigurableApplicationContext applicationContext;
    private final ScheduledExecutorService scheduler = Executors.newSingleThreadScheduledExecutor();

    public InstallUpdateImpl(ConfigurableApplicationContext applicationContext) {
        this.applicationContext = applicationContext;
    }

    @Override
    public void execute(Path installerFile) {
        if (installerFile == null || !Files.exists(installerFile)) {
            throw new UpdateNotAvailableException("No downloaded installer file found for: " + installerFile);
        }

        File verifiedInstaller = verifyWithinTempDir(installerFile.toFile(), ReleaseTempFolder.resolve());

        logger.info("Starting installer for file: {}", verifiedInstaller.getAbsolutePath());
        launchInstallerProcess(verifiedInstaller);
        scheduleApplicationShutdown();
    }

    private void launchInstallerProcess(File installerFile) {
        String os = System.getProperty("os.name").toLowerCase();
        String filePath = installerFile.getAbsolutePath();
        long currentPid = ProcessHandle.current().pid();

        try {
            if (os.contains("win")) {
                launchWindowsInstaller(filePath, currentPid);
            } else if (os.contains("mac")) {
                launchMacInstaller(filePath, currentPid);
            } else {
                launchLinuxInstaller(filePath, currentPid);
            }
        } catch (IOException e) {
            logger.error("Failed to launch native installer process", e);
            throw new IllegalStateException("Failed to launch application installer process.", e);
        }
    }

    private void launchLinuxInstaller(String filePath, long currentPid) throws IOException {
        String installProgram = filePath.endsWith(".rpm") ? "rpm" : "dpkg";
        String installFlag = filePath.endsWith(".rpm") ? "-Uvh" : "-i";

        // Wait for the current process to exit, install the package with pkexec,
        // then relaunch the app fully detached (nohup + setsid) so it survives
        // after the installer script exits.
        // Uses ";" instead of "&&" so the app always relaunches, even if the
        // installation fails — the previous version is still installed and
        // should remain usable.
        String script =
                "PID=\"$1\"; FILE=\"$2\"; PROGRAM=\"$3\"; FLAG=\"$4\"; " +
                        "while kill -0 \"$PID\" 2>/dev/null; do sleep 0.2; done; " +
                        "pkexec \"$PROGRAM\" \"$FLAG\" \"$FILE\"; " +
                        "nohup setsid devaulty >/dev/null 2>&1 &";

        List<String> command = List.of(
                "bash", "-c", script, "bash",
                String.valueOf(currentPid), filePath, installProgram, installFlag
        );

        startDetached(command);
    }

    private void launchMacInstaller(String filePath, long currentPid) throws IOException {
        // Read the note about limitations at the top of the class.
        // After mounting the DMG the app is always relaunched so the user is
        // not left without a running instance.  The relaunch starts the
        // currently-installed version; the user still needs to manually copy
        // the new .app from the mounted image into /Applications.
        String script =
                "PID=\"$1\"; FILE=\"$2\"; " +
                        "while kill -0 \"$PID\" 2>/dev/null; do sleep 0.2; done; " +
                        "open \"$FILE\"; " +
                        "sleep 2; open -a Devaulty";

        List<String> command = List.of(
                "bash", "-c", script, "bash",
                String.valueOf(currentPid), filePath
        );

        startDetached(command);
    }

    private void launchWindowsInstaller(String filePath, long currentPid) throws IOException {
        String installDir = System.getenv(WINDOWS_INSTALL_DIR_ENV);
        if (installDir == null || installDir.isBlank()) {
            installDir = WINDOWS_DEFAULT_INSTALL_DIR;
            logger.warn("{} not set, falling back to default install path: {}",
                    WINDOWS_INSTALL_DIR_ENV, installDir);
        }
        String executablePath = installDir + "\\devaulty.exe";

        String script =
                "param($ProcId, $File, $ExePath) " +
                        "Wait-Process -Id $ProcId -ErrorAction SilentlyContinue; " +
                        "Start-Process msiexec.exe -ArgumentList '/i', $File, '/qb' -Wait; " +
                        "Start-Process $ExePath";

        List<String> command = List.of(
                "powershell", "-NoProfile", "-WindowStyle", "Hidden",
                "-Command", script,
                String.valueOf(currentPid), filePath, executablePath
        );

        startDetached(command);
    }

    protected void startDetached(List<String> command) throws IOException {
        ProcessBuilder builder = new ProcessBuilder(command);
        builder.redirectOutput(ProcessBuilder.Redirect.DISCARD);
        builder.redirectError(ProcessBuilder.Redirect.DISCARD);
        builder.start();
    }

    private File verifyWithinTempDir(File installerFile, File tempDir) {
        try {
            Path canonicalInstaller = installerFile.getCanonicalFile().toPath();
            Path canonicalTempDir = tempDir.getCanonicalFile().toPath();

            if (!canonicalInstaller.startsWith(canonicalTempDir)) {
                throw new IllegalStateException("Installer file resolved outside of the expected temp directory.");
            }
            return canonicalInstaller.toFile();
        } catch (IOException e) {
            throw new IllegalStateException("Failed to resolve installer file path.", e);
        }
    }

    // Temp folder resolution delegated to ReleaseTempFolder to ensure both
    // DownloadUpdateImpl and InstallUpdateImpl always use the same directory.

    private void scheduleApplicationShutdown() {
        scheduler.schedule(() -> {
            logger.info("Shutting down current application instance for update replacement...");
            int exitCode = SpringApplication.exit(applicationContext, () -> 0);
            scheduler.shutdown();
            System.exit(exitCode);
        }, SHUTDOWN_DELAY_SECONDS, TimeUnit.SECONDS);
    }
}