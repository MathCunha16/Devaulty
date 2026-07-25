package com.devaulty.backend.infrastructure.configuration;

import com.devaulty.backend.application.impl.release.CheckForUpdatesImpl;
import com.devaulty.backend.application.impl.release.DownloadUpdateImpl;
import com.devaulty.backend.application.impl.release.InstallUpdateImpl;
import com.devaulty.backend.application.port.in.release.CheckForUpdatesUseCase;
import com.devaulty.backend.application.port.in.release.DownloadUpdateUseCase;
import com.devaulty.backend.application.port.in.release.InstallUpdateUseCase;
import com.devaulty.backend.application.port.out.external.release.ReleasePort;
import com.devaulty.backend.infrastructure.properties.DevaultyProperties;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ReleaseConfig {

    @Bean
    public CheckForUpdatesUseCase checkForUpdatesUseCase(
        ReleasePort releasePort,
        DevaultyProperties devaultyProperties
    ){
        return new CheckForUpdatesImpl(releasePort,
                devaultyProperties
        );
    }

    @Bean
    public DownloadUpdateUseCase downloadUpdateUseCase(
            ReleasePort releasePort,
            CheckForUpdatesUseCase checkForUpdatesUseCase,
            InstallUpdateUseCase installUpdateUseCase
    ){
        return new DownloadUpdateImpl(
                releasePort,
                checkForUpdatesUseCase,
                installUpdateUseCase
        );
    }

    @Bean
    public InstallUpdateUseCase installUpdateUseCase(
            ConfigurableApplicationContext configurableApplicationContext
    ){
        return new InstallUpdateImpl(
                configurableApplicationContext
        );
    }
}
