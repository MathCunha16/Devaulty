package com.devaulty.backend.application.port.out.security;

import javax.crypto.SecretKey;
import java.time.Duration;

public interface MasterKeySessionPort {
    void setKey(SecretKey key);
    SecretKey getKey();
    boolean hasKey();
    void clear();
    void touch();
    Long getSecondsRemaining(Duration timeout);
    SecretKey getOrClearIfExpired(Duration timeout);
}
