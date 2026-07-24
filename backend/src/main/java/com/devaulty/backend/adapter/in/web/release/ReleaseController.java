package com.devaulty.backend.adapter.in.web.release;

import com.devaulty.backend.adapter.in.web.release.dto.AppUpdateInfoResponse;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/release")
public class ReleaseController implements ReleaseApi{

    private final CheckForUpdatesUseCase checkForUpdatesUseCase;
    private final ReleaseWebMapper releaseWebMapper;

    public ReleaseController(CheckForUpdatesUseCase checkForUpdatesUseCase, ReleaseWebMapper releaseWebMapper) {
        this.checkForUpdatesUseCase = checkForUpdatesUseCase;
        this.releaseWebMapper = releaseWebMapper;
    }

    @GetMapping("/check")
    public ResponseEntity<AppUpdateInfoResponse> checkUpdates(){
        return ResponseEntity.ok(releaseWebMapper.toAppUpdateInfoResponse(checkForUpdatesUseCase.execute()));
    }


}
