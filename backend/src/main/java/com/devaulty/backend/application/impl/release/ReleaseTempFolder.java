package com.devaulty.backend.application.impl.release;

import java.io.File;

/**
 * Shared utility for resolving the OS-specific temporary directory used by
 * the Devaulty update pipeline. Both {@link DownloadUpdateImpl} and
 * {@link InstallUpdateImpl} must resolve the same directory so that the
 * installer file written by the download step can be located by the install
 * step without path mismatches.
 * Can be overridden via {@code -Ddevaulty.temp.dir} (e.g. in tests for isolation).
 */
public final class ReleaseTempFolder {

    private ReleaseTempFolder() {}

    public static File resolve() {
        String customDir = System.getProperty("devaulty.temp.dir");
        if (customDir != null && !customDir.isBlank()) {
            File tempDir = new File(customDir);
            if (!tempDir.exists()) {
                tempDir.mkdirs();
            }
            if (!tempDir.isDirectory()) {
                throw new IllegalStateException(
                        "Configured devaulty.temp.dir is not a valid directory: " + customDir);
            }
            return tempDir;
        }

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
}
