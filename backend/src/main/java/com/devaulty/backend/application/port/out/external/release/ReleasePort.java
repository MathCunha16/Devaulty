package com.devaulty.backend.application.port.out.external.release;

import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;

import java.io.InputStream;

public interface ReleasePort {
    LatestReleaseInfo getLatestRelease();

    InputStream downloadAsset(String downloadUrl);
}
