package com.devaulty.backend.application.port.out.security.dto;

public record CryptoResultDto(
        byte[] cipherText,
        byte[] iv,
        byte[] authTag
) {
}
