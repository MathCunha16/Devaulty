package com.devaulty.backend.adapter.out.crypto;

import com.devaulty.backend.application.port.out.security.MasterKeySessionPort;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.util.Arrays;

@Component
public class MasterKeySessionHolder implements MasterKeySessionPort {

    private SecretKey masterKey = null;
    private byte[] rawKeyBytes = null;

    public synchronized void setKey(SecretKey key) {
        this.masterKey = key;
    }

    public synchronized SecretKey getKey() {
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
    }

}
