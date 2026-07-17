package com.devaulty.backend.adapter.out.crypto;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.time.Duration;
import java.time.Instant;
import java.util.Arrays;

@Component
public class MasterKeySessionHolder implements MasterKeySessionPort {

    private SecretKey masterKey = null;
    private byte[] rawKeyBytes = null;
    private volatile Instant lastActivityAt;

    private static final Duration DEFAULT_TIMEOUT = Duration.ofMinutes(15);

    public synchronized void setKey(SecretKey key) {
        if (this.rawKeyBytes != null) {
            Arrays.fill(this.rawKeyBytes, (byte) 0);
            this.rawKeyBytes = null;
        }
        this.masterKey = key;
        if (key != null && key.getEncoded() != null) {
            this.rawKeyBytes = key.getEncoded();
        }
        if (key != null) {
            this.lastActivityAt = Instant.now();
        } else {
            this.lastActivityAt = null;
        }
    }

    /**
     * Retrieves the active master key and automatically resets the inactivity timer.
     *
     * <p><b>⚠ SIDE EFFECT:</b> Calling this method triggers a {@link #touch()},
     * sliding the session expiration window forward.</p>
     */

    public synchronized SecretKey getKey() {
        if (hasKey()) {
            if (isExpired(DEFAULT_TIMEOUT)) {
                clear();
                return null;
            }
            touch();
        }
        return this.masterKey;
    }

    public synchronized boolean hasKey() {
        return this.masterKey != null;
    }

    public synchronized void clear() {
        if (rawKeyBytes != null) {
            // Best-effort reference clearing; does not clear internal immutable copies in SecretKey
            Arrays.fill(rawKeyBytes, (byte) 0);
        }
        this.masterKey = null;
        this.rawKeyBytes = null;
        this.lastActivityAt = null;
    }

    public synchronized void touch() {
        if (hasKey()) {
            if (isExpired(DEFAULT_TIMEOUT)) {
                clear();
                return;
            }
            this.lastActivityAt = Instant.now();
        }
    }

    public synchronized SecretKey getOrClearIfExpired(Duration timeout) {
        if (hasKey() && isExpired(timeout)) {
                clear();
                return null;
            }

        return this.masterKey;
    }

    public synchronized Long getSecondsRemaining(Duration timeout) {
        if (this.masterKey == null || lastActivityAt == null) {
            return 0L;
        }
        long elapsedSeconds = Duration.between(lastActivityAt, Instant.now()).toSeconds();
        long remaining = timeout.toSeconds() - elapsedSeconds;

        // Return 0 if time is expired
        return Math.max(0L, remaining);
    }

    private boolean isExpired(Duration timeout) {
        if (lastActivityAt == null) {
            return true;
        }
        long elapsedSeconds = Duration.between(lastActivityAt, Instant.now()).toSeconds();
        return elapsedSeconds >= timeout.toSeconds();
    }

}
