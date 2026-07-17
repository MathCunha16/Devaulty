package com.devaulty.backend.adapter.out.crypto;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import javax.crypto.SecretKey;
import java.security.SecureRandom;

import static org.junit.jupiter.api.Assertions.*;

class Argon2KeyDeriverAdapterTest {

    private Argon2KeyDeriverAdapter keyDeriver;

    @BeforeEach
    void setUp() {
        keyDeriver = new Argon2KeyDeriverAdapter();
    }

    @Test
    void shouldDeriveValidKeyFromPasswordAndSalt() {
        // Arrange
        char[] password = "mySecretMasterPassword".toCharArray();
        byte[] salt = new byte[16];
        new SecureRandom().nextBytes(salt);

        // Act
        SecretKey key = keyDeriver.deriveKey(password, salt);

        // Assert
        assertNotNull(key);
        assertEquals("AES", key.getAlgorithm());
        assertNotNull(key.getEncoded());
        assertEquals(32, key.getEncoded().length); // 256 bits

        // Check original password input was NOT mutated by the keyDeriver call itself
        // (Wiping standard in Devaulty happens in the use cases finally block, but the keyDeriver method 
        // doesn't touch the original password reference, it copies/encodes it to bytes and clears the bytes copy).
        assertEquals('m', password[0]);
    }

    @Test
    void shouldDeriveSameKey_forSameInputs() {
        // Arrange
        char[] password = "mySecretMasterPassword".toCharArray();
        byte[] salt = new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16};

        // Act
        SecretKey key1 = keyDeriver.deriveKey(password, salt);
        SecretKey key2 = keyDeriver.deriveKey(password, salt);

        // Assert
        assertArrayEquals(key1.getEncoded(), key2.getEncoded());
    }

    @Test
    void shouldDeriveDifferentKey_whenPasswordsAreDifferent() {
        // Arrange
        char[] password1 = "mySecretMasterPassword".toCharArray();
        char[] password2 = "differentMasterPassword".toCharArray();
        byte[] salt = new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16};

        // Act
        SecretKey key1 = keyDeriver.deriveKey(password1, salt);
        SecretKey key2 = keyDeriver.deriveKey(password2, salt);

        // Assert
        assertFalse(java.util.Arrays.equals(key1.getEncoded(), key2.getEncoded()));
    }

    @Test
    void shouldDeriveDifferentKey_whenSaltsAreDifferent() {
        // Arrange
        char[] password = "mySecretMasterPassword".toCharArray();
        byte[] salt1 = new byte[]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16};
        byte[] salt2 = new byte[]{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 10, 11, 12, 13, 14, 15};

        // Act
        SecretKey key1 = keyDeriver.deriveKey(password, salt1);
        SecretKey key2 = keyDeriver.deriveKey(password, salt2);

        // Assert
        assertFalse(java.util.Arrays.equals(key1.getEncoded(), key2.getEncoded()));
    }
}
