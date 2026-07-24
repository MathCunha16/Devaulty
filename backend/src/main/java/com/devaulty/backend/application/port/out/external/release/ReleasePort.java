package com.devaulty.backend.application.port.out.external.release;

import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;

public interface ReleasePort {

    LatestReleaseInfo getLatestRelease();
}
