package com.devaulty.backend.adapter.in.scheduler;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Component;

import java.time.Duration;

@Component
public class VaultAutoLockScheduler {

    private final Duration inactivityTimeout = Duration.ofMinutes(15);

    private static final Logger logger = LoggerFactory.getLogger(VaultAutoLockScheduler.class);

    private final MasterKeySessionPort masterKeySessionPort;

    public VaultAutoLockScheduler(MasterKeySessionPort masterKeySessionPort) {
        this.masterKeySessionPort = masterKeySessionPort;
    }

    @Scheduled(fixedDelay = 10000) // Runs every 10 seconds
    public void purgeExpiredSession(){

        if(!masterKeySessionPort.hasKey()){ // if no master key is set, do nothing
            return;
        }

        if (masterKeySessionPort.getSecondsRemaining(inactivityTimeout) <= 0) {
            masterKeySessionPort.clear();
            logger.info("Inactivity Timeout reached. Vault locked.");
        }
    }
}
