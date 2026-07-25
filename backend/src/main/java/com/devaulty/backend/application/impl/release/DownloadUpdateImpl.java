package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.exception.UpdateNotAvailableException;
import com.devaulty.backend.application.port.in.release.*;
import com.devaulty.backend.application.port.in.release.enums.UpdateStatus;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import org.springframework.core.io.buffer.DataBufferUtils;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.scheduler.Schedulers;

import java.io.File;
import java.io.IOException;
import java.nio.channels.AsynchronousFileChannel;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.time.Duration;
import java.util.concurrent.atomic.AtomicLong;

public class DownloadUpdateImpl implements DownloadUpdateUseCase {

    private static final Duration PROGRESS_SAMPLE_INTERVAL = Duration.ofMillis(200);
    private static final Duration DOWNLOAD_TIMEOUT = Duration.ofMinutes(15);

    private final ReleasePort releasePort;
    private final CheckForUpdatesUseCase checkForUpdatesUseCase;
    private final InstallUpdateUseCase installUpdateUseCase;

    public DownloadUpdateImpl(ReleasePort releasePort, CheckForUpdatesUseCase checkForUpdatesUseCase, InstallUpdateUseCase installUpdateUseCase) {
        this.releasePort = releasePort;
        this.checkForUpdatesUseCase = checkForUpdatesUseCase;
        this.installUpdateUseCase = installUpdateUseCase;
    }

    @Override
    public Flux<UpdateProgressInfo> execute() {

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

        return startDownloadProcess(downloadUrl, totalBytes, targetPath);
    }

    private String sanitizeFileName(String rawFileName) {
        // Strip any directory separators or null bytes to prevent path traversal
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

    private Flux<UpdateProgressInfo> startDownloadProcess(String downloadUrl, long totalBytes, Path targetPath) {
        AtomicLong downloadedBytes = new AtomicLong(0);
        boolean knownSize = totalBytes > 0;

        UpdateProgressInfo initialProgress = new UpdateProgressInfo(
                UpdateStatus.DOWNLOADING, 0, 0L, totalBytes, null
        );

        return Flux.using(
                () -> AsynchronousFileChannel.open(
                        targetPath,
                        StandardOpenOption.CREATE,
                        StandardOpenOption.WRITE,
                        StandardOpenOption.TRUNCATE_EXISTING
                ),

                fileChannel -> {
                    Flux<UpdateProgressInfo> downloadProgress = releasePort.downloadAsset(downloadUrl)
                            .timeout(DOWNLOAD_TIMEOUT)
                            .concatMap(dataBuffer -> {
                                int chunkSize = dataBuffer.readableByteCount();
                                long currentTotal = downloadedBytes.addAndGet(chunkSize);
                                int percentage = knownSize ? (int) ((currentTotal * 100) / totalBytes) : 0;

                                UpdateProgressInfo progressInfo = new UpdateProgressInfo(
                                        UpdateStatus.DOWNLOADING, percentage, currentTotal, totalBytes, null
                                );

                                return DataBufferUtils.write(Flux.just(dataBuffer), fileChannel)
                                        .then(Mono.just(progressInfo));
                            })
                            .distinctUntilChanged(p -> knownSize ? p.percentage() : p.downloadedBytes())
                            .sample(PROGRESS_SAMPLE_INTERVAL)
                            .concatWith(Mono.defer(() -> Mono.just(new UpdateProgressInfo(
                                    UpdateStatus.DOWNLOADING,
                                    knownSize ? 100 : 0,
                                    downloadedBytes.get(),
                                    totalBytes,
                                    null
                            ))));


                    Mono<UpdateProgressInfo> installAndComplete = Mono
                            .fromRunnable(() -> installUpdateUseCase.execute(targetPath))
                            .subscribeOn(Schedulers.boundedElastic())
                            .thenReturn(new UpdateProgressInfo(
                                    UpdateStatus.COMPLETED, 100, totalBytes, totalBytes, null
                            ))
                            .onErrorResume(installError -> Mono.just(new UpdateProgressInfo(
                                    UpdateStatus.FAILED,
                                    100,
                                    totalBytes,
                                    totalBytes,
                                    "Download Succeed, but the application failed to Install : " + installError.getMessage()
                            )));

                    UpdateProgressInfo installingProgress = new UpdateProgressInfo(
                            UpdateStatus.INSTALLING, 100, totalBytes, totalBytes, null
                    );

                    return Mono.just(initialProgress)
                            .concatWith(downloadProgress)
                            .concatWith(Mono.just(installingProgress))
                            .concatWith(installAndComplete);
                },

                this::closeChannel
        ).onErrorResume(ex -> deleteQuietly(targetPath)
                .then(Mono.just(new UpdateProgressInfo(
                        UpdateStatus.FAILED, 0, downloadedBytes.get(), totalBytes, ex.getMessage()
                )))
        );
    }

    private Mono<Void> deleteQuietly(Path path) {
        return Mono.fromRunnable(() -> {
            try {
                Files.deleteIfExists(path);
            } catch (Exception ignored) {
                // Best-effort cleanup
            }
        });
    }

    private void closeChannel(AsynchronousFileChannel channel) {
        if (channel != null && channel.isOpen()) {
            try {
                channel.close();
            } catch (Exception ignored) {
                // Defensive channel closure
            }
        }
    }
}