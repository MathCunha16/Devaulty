package com.devaulty.backend.adapter.out.crypto;

import com.devaulty.backend.application.exception.CryptoException;
import com.devaulty.backend.application.exception.InvalidMasterPasswordException;
import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;

class BouncyCastleCryptoAdapterTest {

    private BouncyCastleCryptoAdapter cryptoAdapter;
    private SecretKey secretKey;
    private byte[] aad;

    @BeforeEach
    void setUp() {
        cryptoAdapter = new BouncyCastleCryptoAdapter();
        // Generate a 256-bit AES key (32 bytes)
        secretKey = new SecretKeySpec(new byte[32], "AES");
        aad = new byte[]{1, 2, 3, 4};
    }

    @Test
    void shouldEncryptAndDecryptSuccessfully() {
        // Arrange
        byte[] originalData = "My sensitive secret data here!".getBytes(StandardCharsets.UTF_8);

        // Act - Encrypt
        CryptoResultDto encryptResult = cryptoAdapter.encrypt(originalData, secretKey, aad);

        // Assert Encryption metadata
        assertNotNull(encryptResult);
        assertNotNull(encryptResult.cipherText());
        assertNotNull(encryptResult.iv());
        assertEquals(12, encryptResult.iv().length); // GCM standard IV is 12 bytes
        assertNotNull(encryptResult.authTag());
        assertEquals(16, encryptResult.authTag().length); // GCM standard Auth Tag is 16 bytes

        // Act - Decrypt
        byte[] decryptedData = cryptoAdapter.decrypt(
                encryptResult.cipherText(),
                encryptResult.iv(),
                encryptResult.authTag(),
                secretKey,
                aad
        );

        // Assert decrypted payload match
        assertArrayEquals(originalData, decryptedData);
        assertEquals("My sensitive secret data here!", new String(decryptedData, StandardCharsets.UTF_8));
    }

    @Test
    void shouldThrowInvalidMasterPasswordException_whenAuthTagIsCorrupted() {
        // Arrange
        byte[] originalData = "Sensitive information".getBytes(StandardCharsets.UTF_8);
        CryptoResultDto encryptResult = cryptoAdapter.encrypt(originalData, secretKey, aad);

        // Corrupt the auth tag
        byte[] corruptedAuthTag = encryptResult.authTag().clone();
        corruptedAuthTag[0] ^= 1; // Flip one bit

        // Act & Assert
        byte[] cipherText = encryptResult.cipherText();
        byte[] iv = encryptResult.iv();
        assertThrows(InvalidMasterPasswordException.class, () -> {
            cryptoAdapter.decrypt(
                    cipherText,
                    iv,
                    corruptedAuthTag,
                    secretKey,
                    aad
            );
        });
    }

    @Test
    void shouldThrowInvalidMasterPasswordException_whenCipherTextIsCorrupted() {
        // Arrange
        byte[] originalData = "Sensitive information".getBytes(StandardCharsets.UTF_8);
        CryptoResultDto encryptResult = cryptoAdapter.encrypt(originalData, secretKey, aad);

        // Corrupt the ciphertext
        byte[] corruptedCipherText = encryptResult.cipherText().clone();
        corruptedCipherText[0] ^= 1;

        // Act & Assert
        byte[] iv = encryptResult.iv();
        byte[] authTag = encryptResult.authTag();
        assertThrows(InvalidMasterPasswordException.class, () -> {
            cryptoAdapter.decrypt(
                    corruptedCipherText,
                    iv,
                    authTag,
                    secretKey,
                    aad
            );
        });
    }

    @Test
    void shouldThrowInvalidMasterPasswordException_whenIvLengthIsInvalid() {
        // Arrange
        byte[] originalData = "Sensitive information".getBytes(StandardCharsets.UTF_8);
        CryptoResultDto encryptResult = cryptoAdapter.encrypt(originalData, secretKey, aad);

        byte[] invalidIv = new byte[5]; // Standard GCM IV should be 12 bytes typically

        // Act & Assert
        byte[] cipherText = encryptResult.cipherText();
        byte[] authTag = encryptResult.authTag();
        assertThrows(InvalidMasterPasswordException.class, () -> {
            cryptoAdapter.decrypt(
                    cipherText,
                    invalidIv,
                    authTag,
                    secretKey,
                    aad
            );
        });
    }

    @Test
    void shouldThrowInvalidMasterPasswordException_whenAadIsDifferent() {
        // Arrange
        byte[] originalData = "Sensitive information".getBytes(StandardCharsets.UTF_8);
        CryptoResultDto encryptResult = cryptoAdapter.encrypt(originalData, secretKey, aad);

        byte[] differentAad = new byte[]{9, 9, 9, 9};

        // Act & Assert
        byte[] cipherText = encryptResult.cipherText();
        byte[] iv = encryptResult.iv();
        byte[] authTag = encryptResult.authTag();
        assertThrows(InvalidMasterPasswordException.class, () -> {
            cryptoAdapter.decrypt(
                    cipherText,
                    iv,
                    authTag,
                    secretKey,
                    differentAad
            );
        });
    }

    @Test
    void shouldThrowCryptoException_whenKeyIsNullOnEncryption() {
        byte[] originalData = "data".getBytes();

        assertThrows(CryptoException.class, () -> {
            cryptoAdapter.encrypt(originalData, null, aad);
        });
    }
}
