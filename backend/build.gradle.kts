plugins {
	java
	id("org.springframework.boot") version "4.1.0"
	id("io.spring.dependency-management") version "1.1.7"
}

val mapStructVersion = "1.6.3"
val swaggerOpenAPIVersion = "3.0.3"
val dotEnvVersion = "5.1.0"

group = "com.devaulty"
version = "0.0.1-SNAPSHOT"

java {
	toolchain {
		languageVersion = JavaLanguageVersion.of(21)
	}
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
	implementation(platform("me.paulschwarz:spring-dotenv-bom:${dotEnvVersion}"))

	developmentOnly("org.springframework.boot:spring-boot-devtools")
	developmentOnly("me.paulschwarz:springboot4-dotenv:${dotEnvVersion}")

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
