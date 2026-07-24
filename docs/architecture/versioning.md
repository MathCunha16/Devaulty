# Application Versioning Architecture

This document defines how the application version is configured, managed, and propagated across the Devaulty codebase, Spring Boot backend, and native `jpackage` installers.

## The Single Source of Truth

To prevent version desynchronization across the build system, backend, and installer packagers, the application version is declared in **one single place**:

> **`src/main/resources/application.yaml`** -> `app.version`

```yaml
app:
  version: "0.1.0-alpha"
```

No other file in the repository (including `build.gradle.kts` or Java source files) hardcodes the version number.

## How Version Information Propagates

### 1. Spring Boot Backend (`DevaultyProperties`)

The backend reads the version dynamically at application startup using a strongly typed, validated configuration properties bean:

- **Class:** `com.devaulty.backend.infrastructure.properties.DevaultyProperties`
- **Prefix:** `app`
- **Validation:** `@Validated` with `@NotBlank` on `appVersion`.

```java
@Component
@Validated
@ConfigurationProperties(prefix = "app")
public class DevaultyProperties {

    @NotBlank(message = "app.version cannot be blank")
    private String appVersion;

    public String getAppVersion() {
        return appVersion;
    }

    public void setAppVersion(String appVersion) {
        this.appVersion = appVersion;
    }
}
```

Use cases like `CheckForUpdatesImpl` inject `DevaultyProperties` to obtain `devaultyProperties.getAppVersion()`. This eliminates the need for Spring Boot `BuildProperties` or `build-info.properties` generation, preventing false-positive hot-reloads during local development.

### 2. Gradle Build System & Installers (`build.gradle.kts`)

`build.gradle.kts` reads `application.yaml` dynamically at configuration time:

```kotlin
// Reads app.version dynamically from application.yaml as Single Source of Truth
val yamlVersion = file("src/main/resources/application.yaml")
    .readLines()
    .firstOrNull { it.trim().startsWith("version:") }
    ?.substringAfter("version:")
    ?.trim()
    ?.removeSurrounding("\"")
    ?.removeSurrounding("'")
    ?: "0.1.0-alpha"

group = "com.devaulty.backend"
version = yamlVersion

// Clean version for Linux/Windows packaging (e.g. "0.1.0-alpha" -> "0.1.0")
val packageVersion = yamlVersion.replace(Regex("(?i)-.*$"), "")

// Numeric 3-part version for macOS packaging (e.g. "0.1.0")
val macPackageVersion = if (packageVersion.split(".").size >= 3) packageVersion else "$packageVersion.0"
```

From this single string:
- **`version`**: Used for standard Gradle task metadata (`0.1.0-alpha`).
- **`packageVersion`**: Strips prerelease suffixes (`-alpha`, `-beta`) to comply with Linux (`.deb`/`.rpm`) and Windows (`.msi`) package managers (`0.1.0`).
- **`macPackageVersion`**: Formats a strict 3-part numeric SemVer required by Apple (`0.1.0`).

## Checklist: Bumping a Version

To release a new version of Devaulty (e.g. from `0.1.0-alpha` to `0.2.0`):

- [ ] Update `app.version` in `src/main/resources/application.yaml`.
- [ ] Run `./gradlew compileJava` to verify property loading.
- [ ] Run `./gradlew packageInstallers` (or `./gradlew packageDeb`) to verify installer generation.
