package com.devaulty.backend.application.port.out.security;

import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;

import javax.crypto.SecretKey;

public interface CryptoPort {

    CryptoResultDto encrypt(byte[] plainData, SecretKey secretKey);

    byte[] decrypt(byte[] cipherText, byte[] iv, byte[] authTag, SecretKey secretKey);
}
