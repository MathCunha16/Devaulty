package com.devaulty.backend.application.impl.security;

import com.devaulty.backend.application.port.in.security.GetSessionStatusUseCase;
import com.devaulty.backend.application.port.in.security.SessionStatus;
import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;

import java.time.Duration;

public class GetSessionStatusImpl implements GetSessionStatusUseCase {

    private static final Duration SESSION_TIMEOUT = Duration.ofMinutes(15);

    private final MasterKeySessionPort sessionHolder;

    public GetSessionStatusImpl(MasterKeySessionPort sessionHolder) {
        this.sessionHolder = sessionHolder;
    }

    @Override
    public SessionStatus execute() {
        boolean active = sessionHolder.hasKey();
        long secondsLeft = sessionHolder.getSecondsRemaining(SESSION_TIMEOUT);
        return new SessionStatus(active, secondsLeft);
    }
}
