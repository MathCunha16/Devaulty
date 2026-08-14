plugins {
    java
    application
    id("org.springframework.boot") version "4.1.0"
    id("io.spring.dependency-management") version "1.1.7"
}

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

val mapStructVersion = "1.6.3"
val swaggerOpenAPIVersion = "3.0.3"
val bouncyCastleVersion = "1.84"


java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

application {
    mainClass.set("com.devaulty.backend.BackendApplication")
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

