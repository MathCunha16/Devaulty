package com.devaulty.backend.application.port.out.external.release;

import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import org.springframework.core.io.buffer.DataBuffer;
import reactor.core.publisher.Flux;

public interface ReleasePort {
    LatestReleaseInfo getLatestRelease();

    Flux<DataBuffer> downloadAsset(String downloadUrl);
}
