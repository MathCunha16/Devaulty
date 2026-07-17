package com.devaulty.backend.application.port.out.security;

import javax.crypto.SecretKey;

public interface KeyDerivationPort {
    SecretKey deriveKey(char[] password, byte[] salt);
}
