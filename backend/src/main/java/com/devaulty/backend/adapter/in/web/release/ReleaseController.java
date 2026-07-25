package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.adapter.in.web.release.dto.CurrentVersionResponse;
import com.devaulty.backend.adapter.in.web.release.dto.AppUpdateInfoResponse;
import com.devaulty.backend.adapter.in.web.release.dto.UpdateDownloadProgressResponse;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.in.release.DownloadUpdateUseCase;
import com.devaulty.backend.application.port.in.release.GetCurrentVersionUseCase;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import reactor.core.publisher.Flux;

@RestController
@RequestMapping("/api/v1/releases")
public class ReleaseController implements ReleaseApi {

    private final CheckForUpdatesUseCase checkForUpdatesUseCase;
    private final DownloadUpdateUseCase downloadUpdateUseCase;
    private final GetCurrentVersionUseCase getCurrentVersionUseCase;
    private final ReleaseWebMapper releaseWebMapper;

    public ReleaseController(CheckForUpdatesUseCase checkForUpdatesUseCase, DownloadUpdateUseCase downloadUpdateUseCase, GetCurrentVersionUseCase getCurrentVersionUseCase, ReleaseWebMapper releaseWebMapper) {
        this.checkForUpdatesUseCase = checkForUpdatesUseCase;
        this.downloadUpdateUseCase = downloadUpdateUseCase;
        this.getCurrentVersionUseCase = getCurrentVersionUseCase;
        this.releaseWebMapper = releaseWebMapper;
    }

    @Override
    @GetMapping("/check")
    public ResponseEntity<AppUpdateInfoResponse> checkUpdates() {
        return ResponseEntity.ok(releaseWebMapper.toAppUpdateInfoResponse(checkForUpdatesUseCase.execute()));
    }

    @Override
    @GetMapping(value = "/download-and-install", produces = MediaType.TEXT_EVENT_STREAM_VALUE)
    public ResponseEntity<Flux<UpdateDownloadProgressResponse>> downloadUpdate() {
        Flux<UpdateDownloadProgressResponse> stream = downloadUpdateUseCase.execute()
                .map(releaseWebMapper::toProgressResponse);

        return ResponseEntity.ok()
                .contentType(MediaType.TEXT_EVENT_STREAM)
                .body(stream);
    }

    @Override
    @GetMapping("/current-app-version")
    public ResponseEntity<CurrentVersionResponse> getAppInfo() {
        return ResponseEntity.ok(new CurrentVersionResponse(getCurrentVersionUseCase.execute()));
    }
}
