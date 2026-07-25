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
import java.nio.channels.AsynchronousFileChannel;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardOpenOption;
import java.time.Duration;
import java.util.concurrent.atomic.AtomicLong;

public class DownloadUpdateImpl implements DownloadUpdateUseCase {

    private static final Duration PROGRESS_SAMPLE_INTERVAL = Duration.ofMillis(200);

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

        File tempDir = createTempFolder();
        String fileName = downloadUrl.substring(downloadUrl.lastIndexOf('/') + 1);
        Path targetPath = new File(tempDir, fileName).toPath();

        return startDownloadProcess(downloadUrl, totalBytes, targetPath);
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

    private File createTempFolder() {
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
        return tempDir;
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