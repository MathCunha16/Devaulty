package com.devaulty.backend.application.port.out.security;

import javax.crypto.SecretKey;

public interface MasterKeySessionPort {
    void setKey(SecretKey key);
    SecretKey getKey();
    boolean hasKey();
    void clear();
}
