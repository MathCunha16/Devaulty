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

    private final Logger logger = LoggerFactory.getLogger(CheckForUpdatesImpl.class);

    private final ReleasePort releasePort;
    private final DevaultyProperties devaultyProperties;

    private static final String OS_RELEASE_FILE = "/etc/os-release";

    public CheckForUpdatesImpl(ReleasePort releasePort, DevaultyProperties devaultyProperties) {
        this.releasePort = releasePort;
        this.devaultyProperties = devaultyProperties;
    }


    @Override
    public AppUpdateInfo execute() {

        LatestReleaseInfo latestRelease = releasePort.getLatestRelease();

        if(latestRelease == null){
            return new AppUpdateInfo(
                    false,
                    "0.1.0-alpha",
                    null,
                    null,
                    null,
                    null,
                    0L,
                    null
            );
        }

        String cleanLatestVersion = latestRelease.tagName().replaceAll("(?i)^v", "").trim();
        String currentVersion = devaultyProperties.getVersion();

        boolean updateAvailable = false;
        if (currentVersion != null) {
            updateAvailable = !currentVersion.equals(cleanLatestVersion);
        }

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

    private ReleaseAssetInfo findAssetForCurrentOs(List<ReleaseAssetInfo> assets) {
        String os = System.getProperty("os.name").toLowerCase();

        String extension;
        if (os.contains("win")) {
            extension = ".msi";
        } else if (os.contains("mac")) {
            extension = ".dmg";
        } else {
            extension = detectLinuxExtension();
            if (extension == null) {
                return null; // No suitable asset found
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
