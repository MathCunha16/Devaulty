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

    public synchronized void setKey(SecretKey key) {
        this.masterKey = key;
        if (key != null && key.getEncoded() != null) {
            this.rawKeyBytes = key.getEncoded();
        }
        touch();
    }


    /**
     * Retrieves the active master key and automatically resets the inactivity timer.
     *
     * <p><b>⚠ SIDE EFFECT:</b> Calling this method triggers a {@link #touch()}, sliding
     * the session expiration window forward.</p>
     *
     * <p>To simply check if the vault is unlocked without resetting the timer,
     * use {@link #hasKey()} instead.</p>
     *
     * @return the active {@link SecretKey}, or {@code null} if the vault is locked.
     */
    public synchronized SecretKey getKey() {
        touch();
        return this.masterKey;
    }

    public synchronized boolean hasKey() {
        return this.masterKey != null;
    }

    public synchronized void clear() {
        if (rawKeyBytes != null) {
            // Fills the array with zeros, securely destroying the key bits in RAM!
            Arrays.fill(rawKeyBytes, (byte) 0);
        }
        this.masterKey = null;
        this.rawKeyBytes = null;
        this.lastActivityAt = null;
    }

    public synchronized void touch(){
        this.lastActivityAt = Instant.now();
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

}
