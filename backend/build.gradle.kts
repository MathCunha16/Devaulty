plugins {
    java
    application
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
    id("org.openjfx.javafxplugin") version "0.1.0"
}

val mapStructVersion = "1.6.3"
val swaggerOpenAPIVersion = "3.0.3"
val bouncyCastleVersion = "1.84"

group = "com.devaulty"
version = "0.0.1-SNAPSHOT"

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    mainClass.set("com.devaulty.backend.desktop.DevaultyDesktop")
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-data-jpa")
    implementation("org.springframework.boot:spring-boot-starter-liquibase")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("org.springframework.boot:spring-boot-starter-webmvc")
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

// task to build frontend (npm run build)
val buildFrontend by tasks.registering(Exec::class) {
    group = "build"
    description = "Executes React/Vite application build"

    workingDir = file("../frontend")

    val isWindows = org.gradle.internal.os.OperatingSystem.current().isWindows()
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
