package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.adapter.in.web.common.BackgroundTaskRunner;
import com.devaulty.backend.adapter.in.web.release.dto.AppUpdateInfoResponse;
import com.devaulty.backend.adapter.in.web.release.dto.CurrentVersionResponse;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.in.release.DownloadUpdateUseCase;
import com.devaulty.backend.application.port.in.release.GetCurrentVersionUseCase;
import com.devaulty.backend.application.port.in.release.UpdateProgressInfo;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

import java.io.IOException;

@RestController
@RequestMapping("/api/v1/releases")
public class ReleaseController implements ReleaseApi {

    private final CheckForUpdatesUseCase checkForUpdatesUseCase;
    private final DownloadUpdateUseCase downloadUpdateUseCase;
    private final GetCurrentVersionUseCase getCurrentVersionUseCase;
    private final ReleaseWebMapper releaseWebMapper;
    private final BackgroundTaskRunner backgroundTaskRunner;

    private static final long SSE_TIMEOUT_MS = 20 * 60 * 1000L;

    public ReleaseController(CheckForUpdatesUseCase checkForUpdatesUseCase, DownloadUpdateUseCase downloadUpdateUseCase, GetCurrentVersionUseCase getCurrentVersionUseCase, ReleaseWebMapper releaseWebMapper, BackgroundTaskRunner backgroundTaskRunner) {
        this.checkForUpdatesUseCase = checkForUpdatesUseCase;
        this.downloadUpdateUseCase = downloadUpdateUseCase;
        this.getCurrentVersionUseCase = getCurrentVersionUseCase;
        this.releaseWebMapper = releaseWebMapper;
        this.backgroundTaskRunner = backgroundTaskRunner;
    }

    @Override
    @GetMapping("/check")
    public ResponseEntity<AppUpdateInfoResponse> checkUpdates() {
        return ResponseEntity.ok(releaseWebMapper.toAppUpdateInfoResponse(checkForUpdatesUseCase.execute()));
    }

    @Override
    @PostMapping("/download-and-install")
    public SseEmitter downloadUpdate() {
        SseEmitter emitter = new SseEmitter(SSE_TIMEOUT_MS);
        backgroundTaskRunner.run(() -> runDownload(emitter));
        return emitter;
    }

    @Override
    @GetMapping("/current-app-version")
    public ResponseEntity<CurrentVersionResponse> getAppInfo() {
        return ResponseEntity.ok(new CurrentVersionResponse(getCurrentVersionUseCase.execute()));
    }

    private void runDownload(SseEmitter emitter) {
        try {
            downloadUpdateUseCase.execute(progress -> sendProgress(emitter, progress));
            emitter.complete();
        } catch (Exception ex) {
            emitter.completeWithError(ex);
        }
    }

    private void sendProgress(SseEmitter emitter, UpdateProgressInfo progress) {
        try {
            emitter.send(releaseWebMapper.toProgressResponse(progress));
        } catch (IOException e) {
            emitter.completeWithError(e);
        }
    }
}
