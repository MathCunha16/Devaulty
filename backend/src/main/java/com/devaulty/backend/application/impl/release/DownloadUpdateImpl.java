package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import com.devaulty.backend.application.port.in.release.*;
import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.util.function.Consumer;

public class DownloadUpdateImpl implements DownloadUpdateUseCase {

    private static final long PROGRESS_EMIT_THRESHOLD_BYTES = 512 * 1024; // 512KB
    private static final int COPY_BUFFER_SIZE = 8 * 1024; // 8KB for read
    private final ReleasePort releasePort;
    private final CheckForUpdatesUseCase checkForUpdatesUseCase;
    private final InstallUpdateUseCase installUpdateUseCase;

    public DownloadUpdateImpl(ReleasePort releasePort, CheckForUpdatesUseCase checkForUpdatesUseCase, InstallUpdateUseCase installUpdateUseCase) {
        this.releasePort = releasePort;
        this.checkForUpdatesUseCase = checkForUpdatesUseCase;
        this.installUpdateUseCase = installUpdateUseCase;
    }

    @Override
    public void execute(Consumer<UpdateProgressInfo> onProgress) {

        AppUpdateInfo updateInfo = checkForUpdatesUseCase.execute();

        if (!updateInfo.updateAvailable() || updateInfo.downloadUrl() == null) {
            throw new UpdateNotAvailableException(
                    "There is no update available for this application"
            );
        }

        String downloadUrl = updateInfo.downloadUrl();
        long totalBytes = updateInfo.downloadSizeInBytes();

        File tempDir = ReleaseTempFolder.resolve();
        String rawFileName = downloadUrl.substring(downloadUrl.lastIndexOf('/') + 1);
        String fileName = sanitizeFileName(rawFileName);

        Path targetPath = new File(tempDir, fileName).toPath();
        verifyWithinTempDir(targetPath, tempDir);

        runDownloadProcess(downloadUrl, totalBytes, targetPath, onProgress);
    }

    private String sanitizeFileName(String rawFileName) {
        String sanitized = rawFileName
                .replace("/", "")
                .replace("\\", "")
                .replace("\0", "")
                .trim();

        if (sanitized.isEmpty()) {
            throw new IllegalArgumentException("Derived installer file name is empty after sanitization.");
        }

        return sanitized;
    }

    private void verifyWithinTempDir(Path targetPath, File tempDir) {
        try {
            Path canonicalTarget = targetPath.toFile().getCanonicalFile().toPath();
            Path canonicalTempDir = tempDir.getCanonicalFile().toPath();

            if (!canonicalTarget.startsWith(canonicalTempDir)) {
                throw new IllegalStateException("Installer target path resolved outside of the expected temp directory.");
            }
        } catch (IOException e) {
            throw new IllegalStateException("Failed to resolve installer target path.", e);
        }
    }

    private void runDownloadProcess(String downloadUrl, long totalBytes, Path targetPath, Consumer<UpdateProgressInfo> onProgress) {
        boolean knownSize = totalBytes > 0;
        long downloadedBytes = 0;
        long bytesSinceLastEmit = 0;

        onProgress.accept(new UpdateProgressInfo(UpdateStatus.DOWNLOADING, 0, 0L, totalBytes, null));

        try (InputStream in = releasePort.downloadAsset(downloadUrl);
             OutputStream out = Files.newOutputStream(
                     targetPath,
                     StandardOpenOption.CREATE,
                     StandardOpenOption.WRITE,
                     StandardOpenOption.TRUNCATE_EXISTING)) {

            byte[] buffer = new byte[COPY_BUFFER_SIZE];
            int read;

            while ((read = in.read(buffer)) != -1) {
                out.write(buffer, 0, read);

                downloadedBytes += read;
                bytesSinceLastEmit += read;

                // Só emite progresso de tempos em tempos, não a cada 8KB lido
                if (bytesSinceLastEmit >= PROGRESS_EMIT_THRESHOLD_BYTES) {
                    int percentage = knownSize ? (int) ((downloadedBytes * 100) / totalBytes) : 0;
                    onProgress.accept(new UpdateProgressInfo(
                            UpdateStatus.DOWNLOADING, percentage, downloadedBytes, totalBytes, null
                    ));
                    bytesSinceLastEmit = 0;
                }
            }

            // Garante um evento final de 100% mesmo que o último lote tenha sido pequeno
            onProgress.accept(new UpdateProgressInfo(
                    UpdateStatus.DOWNLOADING, knownSize ? 100 : 0, downloadedBytes, totalBytes, null
            ));

            onProgress.accept(new UpdateProgressInfo(
                    UpdateStatus.INSTALLING, 100, totalBytes, totalBytes, null
            ));

            try {
                installUpdateUseCase.execute(targetPath);
                onProgress.accept(new UpdateProgressInfo(
                        UpdateStatus.COMPLETED, 100, totalBytes, totalBytes, null
                ));
            } catch (Exception installError) {
                onProgress.accept(new UpdateProgressInfo(
                        UpdateStatus.FAILED, 100, totalBytes, totalBytes,
                        "Download Succeed, but the application failed to Install : " + installError.getMessage()
                ));
            }

        } catch (IOException ex) {
            deleteQuietly(targetPath);
            onProgress.accept(new UpdateProgressInfo(
                    UpdateStatus.FAILED, 0, downloadedBytes, totalBytes, ex.getMessage()
            ));
        }
    }

    private void deleteQuietly(Path path) {
        try {
            Files.deleteIfExists(path);
        } catch (Exception ignored) {
            // Best-effort cleanup
        }
    }
}