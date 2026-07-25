package com.devaulty.backend.application.impl.release;

import com.devaulty.backend.application.port.in.release.AppUpdateInfo;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.application.port.out.external.release.dto.LatestReleaseInfo;
import com.devaulty.backend.application.port.out.external.release.dto.ReleaseAssetInfo;
import com.devaulty.backend.infrastructure.properties.DevaultyProperties;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.util.List;
import java.util.Optional;

public class CheckForUpdatesImpl implements CheckForUpdatesUseCase {

    private final ReleasePort releasePort;
    private final DevaultyProperties devaultyProperties;

    private static final String OS_RELEASE_FILE = "/etc/os-release";
    private static final Logger logger = LoggerFactory.getLogger(CheckForUpdatesImpl.class);

    public CheckForUpdatesImpl(ReleasePort releasePort, DevaultyProperties devaultyProperties) {
        this.releasePort = releasePort;
        this.devaultyProperties = devaultyProperties;
    }

    @Override
    public AppUpdateInfo execute() {

        LatestReleaseInfo latestRelease = releasePort.getLatestRelease();

        if (latestRelease == null) {
            return new AppUpdateInfo(
                    false,
                    devaultyProperties != null ? devaultyProperties.getVersion() : null,
                    null,
                    null,
                    null,
                    null,
                    0L,
                    null
            );
        }

        String cleanLatestVersion = latestRelease.tagName().replaceAll("(?i)^v", "").trim();
        String currentVersion = devaultyProperties != null ? devaultyProperties.getVersion() : null;

        boolean updateAvailable = isNewerVersion(currentVersion, cleanLatestVersion);

        ReleaseAssetInfo targetAsset = findAssetForCurrentOs(latestRelease.assets());
        String downloadUrl = targetAsset != null ? targetAsset.downloadUrl() : null;
        long downloadSize = targetAsset != null ? targetAsset.sizeInBytes() : 0L;

        return new AppUpdateInfo(
                updateAvailable,
                currentVersion,
                latestRelease.tagName(),
                latestRelease.name(),
                latestRelease.body(),
                downloadUrl,
                downloadSize,
                latestRelease.publishedAt()
        );
    }

    private boolean isNewerVersion(String currentVersion, String latestVersion) {
        if (currentVersion == null || latestVersion == null) {
            return false;
        }

        String currentClean = currentVersion.split("-")[0];
        String latestClean = latestVersion.split("-")[0];

        String[] currentParts = currentClean.split("\\.");
        String[] latestParts = latestClean.split("\\.");

        int length = Math.max(currentParts.length, latestParts.length);
        for (int i = 0; i < length; i++) {
            int vCurrent = i < currentParts.length ? parseVersionPart(currentParts[i]) : 0;
            int vLatest = i < latestParts.length ? parseVersionPart(latestParts[i]) : 0;

            if (vLatest > vCurrent) {
                return true;
            }
            if (vLatest < vCurrent) {
                return false;
            }
        }

        return currentVersion.contains("-") && !latestVersion.contains("-");
    }

    private int parseVersionPart(String part) {
        try {
            return Integer.parseInt(part);
        } catch (NumberFormatException e) {
            return 0;
        }
    }

    private ReleaseAssetInfo findAssetForCurrentOs(List<ReleaseAssetInfo> assets) {
        if (assets == null) {
            return null;
        }

        String os = System.getProperty("os.name").toLowerCase();

        String extension;
        if (os.contains("win")) {
            extension = ".msi";
        } else if (os.contains("mac")) {
            extension = ".dmg";
        } else {
            extension = detectLinuxExtension();
            if (extension == null) {
                return null;
            }
        }

        return assets.stream()
                .filter(asset -> asset.fileName().endsWith(extension))
                .findFirst()
                .orElse(null);
    }

    private String detectLinuxExtension() {

        Optional<String> osRelease = readOsRelease();
        if (osRelease.isEmpty()) {
            return null;
        }

        String content = osRelease.get().toLowerCase();

        boolean isRpmBased = content.contains("rhel")
                || content.contains("fedora")
                || content.contains("suse");

        boolean isDebBased = content.contains("debian")
                || content.contains("ubuntu");

        if (isRpmBased) {
            return ".rpm";
        }
        if (isDebBased) {
            return ".deb";
        }

        logger.warn("Could not detect Linux distribution from /etc/os-release");
        return null;
    }

    private Optional<String> readOsRelease() {
        File file = new File(OS_RELEASE_FILE);
        if (!file.exists()) {
            return Optional.empty();
        }
        try {
            return Optional.of(Files.readString(file.toPath()));
        } catch (IOException e) {
            return Optional.empty();
        }
    }
}
