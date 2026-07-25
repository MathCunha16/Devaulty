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

        int hyphenCurrent = currentVersion.indexOf('-');
        int hyphenLatest = latestVersion.indexOf('-');

        String currentCore = hyphenCurrent >= 0 ? currentVersion.substring(0, hyphenCurrent) : currentVersion;
        String latestCore  = hyphenLatest  >= 0 ? latestVersion.substring(0, hyphenLatest)   : latestVersion;

        String currentPre = hyphenCurrent >= 0 ? currentVersion.substring(hyphenCurrent + 1) : null;
        String latestPre  = hyphenLatest  >= 0 ? latestVersion.substring(hyphenLatest + 1)   : null;

        // 1. Compare numeric core (major.minor.patch)
        int coreComparison = compareCore(currentCore, latestCore);
        if (coreComparison != 0) {
            return coreComparison < 0;
        }

        // 2. Same core: stable (no pre-release) > any pre-release (SemVer spec)
        if (currentPre == null && latestPre != null) {
            return false; // current is stable, latest is pre-release → not newer
        }
        if (currentPre != null && latestPre == null) {
            return true;  // current is pre-release, latest is stable → newer
        }
        if (currentPre == null) {
            return false; // both stable, same version
        }

        // 3. Both have pre-release identifiers: compare dot-separated identifiers
        return comparePreRelease(currentPre, latestPre) < 0;
    }

    private int compareCore(String a, String b) {
        String[] aParts = a.split("\\.");
        String[] bParts = b.split("\\.");
        int length = Math.max(aParts.length, bParts.length);
        for (int i = 0; i < length; i++) {
            int va = i < aParts.length ? parseVersionPart(aParts[i]) : 0;
            int vb = i < bParts.length ? parseVersionPart(bParts[i]) : 0;
            if (va != vb) return Integer.compare(va, vb);
        }
        return 0;
    }

    /**
     * Compares two pre-release strings per SemVer 2.0 spec:
     * - Numeric identifiers are compared numerically.
     * - Alphanumeric identifiers are compared lexically (ASCII).
     * - A numeric identifier always has lower precedence than alphanumeric.
     * Returns negative if a < b, positive if a > b, 0 if equal.
     */
    private int comparePreRelease(String a, String b) {
        String[] aIds = a.split("\\.");
        String[] bIds = b.split("\\.");
        int length = Math.max(aIds.length, bIds.length);
        for (int i = 0; i < length; i++) {
            if (i >= aIds.length) return -1; // a has fewer identifiers → a is smaller
            if (i >= bIds.length) return  1;

            String ai = aIds[i];
            String bi = bIds[i];

            boolean aNum = ai.chars().allMatch(Character::isDigit);
            boolean bNum = bi.chars().allMatch(Character::isDigit);

            int cmp;
            if (aNum && bNum) {
                cmp = Integer.compare(Integer.parseInt(ai), Integer.parseInt(bi));
            } else if (aNum) {
                cmp = -1; // numeric < alphanumeric
            } else if (bNum) {
                cmp = 1;
            } else {
                cmp = ai.compareTo(bi);
            }
            if (cmp != 0) return cmp;
        }
        return 0;
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
