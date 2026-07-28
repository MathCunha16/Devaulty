plugins {
    java
    application
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
    id("org.openjfx.javafxplugin") version "0.1.0"
    id("org.beryx.runtime") version "2.0.1"
}

val mapStructVersion = "1.6.3"
val swaggerOpenAPIVersion = "3.0.3"
val bouncyCastleVersion = "1.84"

val yamlVersion = file("src/main/resources/application.yaml")
    .readLines()
    .firstOrNull { it.trim().startsWith("version:") }
    ?.substringAfter("version:")
    ?.trim()
    ?.removeSurrounding("\"")
    ?.removeSurrounding("'")
    ?: "0.1.0-alpha"

group = "com.devaulty"
version = yamlVersion

// Clean version for Linux/Windows (e.g. "0.1.0-alpha" -> "0.1.0")
val packageVersion = yamlVersion.replace(Regex("(?i)-.*$"), "")

// Numeric 3-part version for macOS (e.g. "0.1.0")
val macPackageVersion = if (packageVersion.split(".").size >= 3) packageVersion else "$packageVersion.0"

// macOS jpackage requires the first version component to be >= 1 (CFBundleVersion spec).
// Map 0.x.y → 1.x.y for the jpackageImage task so the image build succeeds on all platforms.
// Each platform installer task already passes its own --app-version independently.
val jpackageImageVersion = run {
    val parts = packageVersion.split(".")
    val major = parts.getOrElse(0) { "1" }.toIntOrNull() ?: 1
    (if (major == 0) listOf("1") + parts.drop(1) else parts).joinToString(".")
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    mainClass.set("com.devaulty.backend.desktop.DevaultyMainLauncher")
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-liquibase")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-webmvc"){
        exclude(group = "org.springframework.boot", module = "spring-boot-starter-tomcat")
    }
    implementation("org.springframework.boot:spring-boot-starter-jetty")
    implementation("org.springframework.security:spring-security-crypto")
    implementation("org.mapstruct:mapstruct:${mapStructVersion}")
    implementation("org.springdoc:springdoc-openapi-starter-webmvc-ui:${swaggerOpenAPIVersion}")
    implementation("org.bouncycastle:bcprov-jdk18on:${bouncyCastleVersion}")

    developmentOnly("org.springframework.boot:spring-boot-devtools")

    runtimeOnly("org.xerial:sqlite-jdbc")
    runtimeOnly("org.hibernate.orm:hibernate-community-dialects")

    // Annotation Processor (mapstruct)
    annotationProcessor("org.mapstruct:mapstruct-processor:${mapStructVersion}")

    testImplementation("org.springframework.boot:spring-boot-starter-data-jpa-test")
    testImplementation("org.springframework.boot:spring-boot-starter-liquibase-test")
    testImplementation("org.springframework.boot:spring-boot-starter-validation-test")
    testImplementation("org.springframework.boot:spring-boot-starter-webmvc-test")

    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

javafx {
    version = "21.0.2"
    modules = listOf("javafx.controls", "javafx.web")
}

tasks.withType<Test> {
    useJUnitPlatform()
}

tasks.named<org.springframework.boot.gradle.tasks.bundling.BootJar>("bootJar") {
    mainClass.set("com.devaulty.backend.BackendApplication")
}

tasks.named<org.springframework.boot.gradle.tasks.run.BootRun>("bootRun") {
    mainClass.set("com.devaulty.backend.BackendApplication")
}

tasks.withType<JavaExec>().configureEach {
    jvmArgs(
        "-Xms64m",
        "-Xmx256m",
        "-XX:MetaspaceSize=96m",
        "-XX:MaxMetaspaceSize=192m",
        "-XX:ParallelGCThreads=2",
        "-XX:ConcGCThreads=1",
        "-XX:+UseG1GC",
        "-XX:MaxGCPauseMillis=100"
    )
}

// task to build frontend (npm run build)
val buildFrontend by tasks.registering(Exec::class) {
    group = "build"
    description = "Executes React/Vite application build"

    workingDir = file("../frontend")

    val isWindows = org.gradle.internal.os.OperatingSystem.current().isWindows
    if (isWindows) {
        commandLine("cmd", "/c", "npm run build")
    } else {
        commandLine("bash", "-l", "-c", "npm run build")
    }

    // Gradle will only execute the task if any of these inputs/outputs have changed
    inputs.dir("../frontend/src")
    inputs.file("../frontend/package.json")
    outputs.dir("../frontend/dist")
}

// Copy frontend build to statics resources
val copyFrontendResources by tasks.registering(Copy::class) {
    group = "build"
    description = "Copies frontend compilated files to statics resources from Spring Boot"

    dependsOn(buildFrontend)

    from("../frontend/dist")
    into(layout.buildDirectory.dir("resources/main/static"))
}

// Ties frontend build to Spring Boot resources
tasks.named("processResources") {
    dependsOn(copyFrontendResources)
}

runtime {
    options.set(listOf("--strip-debug", "--compress", "zip-6", "--no-header-files", "--no-man-pages"))
    modules.set(
        listOf(
            "java.xml", "java.sql", "java.naming", "java.desktop", "java.management",
            "java.instrument", "java.scripting", "java.security.jgss", "jdk.unsupported",
            "java.compiler", "java.net.http", "java.logging", "java.prefs", "java.rmi",
            "java.transaction.xa", "jdk.crypto.ec", "jdk.zipfs"
        )
    )

    jpackage {
        imageName = "devaulty"
        appVersion = jpackageImageVersion
        imageOptions = listOf(
            // 1. Memory do Heap Limits (Obj dynamic RAM)
            "--java-options", "-Xms64m",
            "--java-options", "-Xmx256m",

            // 2. Metaspace limits (Spring/Hibernate Class Metadata Memory)
            "--java-options", "-XX:MetaspaceSize=96m",
            "--java-options", "-XX:MaxMetaspaceSize=192m",

            // 3. Garbage Collector (GC) Thread Control
            "--java-options", "-XX:ParallelGCThreads=2",
            "--java-options", "-XX:ConcGCThreads=1",
            "--java-options", "-XX:+UseG1GC",
            "--java-options", "-XX:MaxGCPauseMillis=100",

            // 4. Spring Production Profile
            "--java-options", "-Dspring.profiles.active=prod"
        )

    }
}

val isMac = org.gradle.internal.os.OperatingSystem.current().isMacOsX
val appImageDir = if (isMac) {
    layout.buildDirectory.dir("jpackage/devaulty.app")
} else {
    layout.buildDirectory.dir("jpackage/devaulty")
}

fun getResourceDirArgs(dirPath: String): List<String> {
    val dir = file(dirPath)
    return if (dir.exists() && (dir.listFiles()?.isNotEmpty() == true)) {
        listOf("--resource-dir", dir.absolutePath)
    } else {
        emptyList()
    }
}

val packageDeb by tasks.registering(Exec::class) {
    group = "distribution"
    dependsOn("jpackageImage")
    onlyIf { org.gradle.internal.os.OperatingSystem.current().isLinux }
    commandLine(
        buildList {
            addAll(listOf(
                "jpackage",
                "--type", "deb",
                "--app-image", appImageDir.get().asFile.path,
                "--name", "devaulty",
                "--app-version", packageVersion,
                "--vendor", "Devaulty",
                "--icon", file("src/main/resources/static/icon/devaulty-icon.png").absolutePath
            ))
            addAll(getResourceDirArgs("src/main/resources/jpackage/linux"))
            addAll(listOf(
                "--linux-shortcut",
                "--linux-menu-group", "Utility",
                "--dest", layout.buildDirectory.dir("jpackage/deb").get().asFile.path
            ))
        }
    )
}

val packageRpm by tasks.registering(Exec::class) {
    group = "distribution"
    dependsOn("jpackageImage")
    onlyIf { org.gradle.internal.os.OperatingSystem.current().isLinux }
    commandLine(
        buildList {
            addAll(listOf(
                "jpackage",
                "--type", "rpm",
                "--app-image", appImageDir.get().asFile.path,
                "--name", "devaulty",
                "--app-version", packageVersion,
                "--vendor", "Devaulty",
                "--icon", file("src/main/resources/static/icon/devaulty-icon.png").absolutePath
            ))
            addAll(getResourceDirArgs("src/main/resources/jpackage/linux"))
            addAll(listOf(
                "--linux-shortcut",
                "--linux-menu-group", "Utility",
                "--dest", layout.buildDirectory.dir("jpackage/rpm").get().asFile.path
            ))
        }
    )
}

val packageMsi by tasks.registering(Exec::class) {
    group = "distribution"
    dependsOn("jpackageImage")
    onlyIf { org.gradle.internal.os.OperatingSystem.current().isWindows }
    commandLine(
        buildList {
            addAll(listOf(
                "jpackage",
                "--type", "msi",
                "--app-image", appImageDir.get().asFile.path,
                "--name", "devaulty",
                "--app-version", packageVersion,
                "--vendor", "Devaulty",
                "--icon", file("src/main/resources/static/icon/devaulty-icon.ico").absolutePath
            ))
            addAll(getResourceDirArgs("src/main/resources/jpackage/windows"))
            addAll(listOf(
                "--win-shortcut",
                "--win-menu",
                "--win-menu-group", "Utility",
                "--dest", layout.buildDirectory.dir("jpackage/msi").get().asFile.path
            ))
        }
    )
}

val packageDmg by tasks.registering(Exec::class) {
    group = "distribution"
    dependsOn("jpackageImage")
    onlyIf { org.gradle.internal.os.OperatingSystem.current().isMacOsX }
    commandLine(
        buildList {
            addAll(listOf(
                "jpackage",
                "--type", "dmg",
                "--app-image", appImageDir.get().asFile.path,
                "--name", "devaulty",
                "--app-version", macPackageVersion,
                "--vendor", "Devaulty",
                "--icon", file("src/main/resources/static/icon/devaulty-icon.icns").absolutePath
            ))
            addAll(getResourceDirArgs("src/main/resources/jpackage/macos"))
            addAll(listOf(
                "--mac-package-name", "Devaulty",
                "--dest", layout.buildDirectory.dir("jpackage/dmg").get().asFile.path
            ))
        }
    )
}

val packageInstallers by tasks.registering {
    group = "distribution"
    description = "Generates all applicable installer for the current OS"
    dependsOn(
        when {
            org.gradle.internal.os.OperatingSystem.current().isLinux -> listOf(packageDeb, packageRpm)
            org.gradle.internal.os.OperatingSystem.current().isWindows -> listOf(packageMsi)
            org.gradle.internal.os.OperatingSystem.current().isMacOsX -> listOf(packageDmg)
            else -> emptyList()
        }
    )
}
