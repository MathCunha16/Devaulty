package com.devaulty.backend.adapter.out.crypto;

import com.devaulty.backend.application.exception.CryptoException;
import com.devaulty.backend.application.exception.InvalidMasterPasswordException;
import com.devaulty.backend.application.port.out.security.CryptoPort;
import com.devaulty.backend.application.port.out.security.dto.CryptoResultDto;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.springframework.stereotype.Component;

import javax.crypto.AEADBadTagException;
import javax.crypto.Cipher;
import javax.crypto.SecretKey;
import javax.crypto.spec.GCMParameterSpec;
import java.security.SecureRandom;
import java.security.Security;
import java.util.Arrays;

@Component
public class BouncyCastleCryptoAdapter implements CryptoPort {

    private static final String ALGORITHM = "AES/GCM/NoPadding";
    private static final int TAG_BIT_LENGTH = 128; // 16 bytes
    private static final int IV_BYTE_LENGTH = 12; // 12 bytes

    private final SecureRandom secureRandom;

    public BouncyCastleCryptoAdapter() {
        // Registers Bouncy Castle as a security provider in the Java Runtime Environment
        Security.addProvider(new BouncyCastleProvider());
        this.secureRandom = new SecureRandom();
    }

    @Override
    public CryptoResultDto encrypt(byte[] plainData, SecretKey secretKey) {
        try {
            // 1. Generates a unique random 12-byte IV for THIS encryption operation
            byte[] iv = new byte[IV_BYTE_LENGTH];
            secureRandom.nextBytes(iv);

            // 2. Initializes the Bouncy Castle Cipher in ENCRYPT mode
            Cipher cipher = Cipher.getInstance(ALGORITHM, "BC");
            GCMParameterSpec parameterSpec = new GCMParameterSpec(TAG_BIT_LENGTH, iv);
            cipher.init(Cipher.ENCRYPT_MODE, secretKey, parameterSpec);

            // 3. Encrypts the plaintext
            byte[] cipherTextWithTag = cipher.doFinal(plainData);

            // 4. Separates the ciphertext and the authentication tag (the last 16 bytes).
            int cipherTextLength = cipherTextWithTag.length - 16;
            byte[] cipherText = Arrays.copyOfRange(cipherTextWithTag, 0, cipherTextLength);
            byte[] authTag = Arrays.copyOfRange(cipherTextWithTag, cipherTextLength, cipherTextWithTag.length);

            return new CryptoResultDto(cipherText, iv, authTag);

        } catch (Exception e) {
            throw new CryptoException("Critical failure when attempting to encrypt payload", e);
        }
    }

    @Override
    public byte[] decrypt(byte[] cipherText, byte[] iv, byte[] authTag, SecretKey secretKey) {
        try {
            // 1. Recombines the ciphertext and authentication tag into a single byte array
            byte[] cipherTextWithTag = new byte[cipherText.length + authTag.length];
            System.arraycopy(cipherText, 0, cipherTextWithTag, 0, cipherText.length);
            System.arraycopy(authTag, 0, cipherTextWithTag, cipherText.length, authTag.length);

            // 2. Initializes the Bouncy Castle Cipher in DECRYPT mode
            Cipher cipher = Cipher.getInstance(ALGORITHM, "BC");
            GCMParameterSpec parameterSpec = new GCMParameterSpec(TAG_BIT_LENGTH, iv);
            cipher.init(Cipher.DECRYPT_MODE, secretKey, parameterSpec);

            // 3. Performs the reverse operation and verifies the authentication tag
            return cipher.doFinal(cipherTextWithTag);

        } catch (AEADBadTagException e) {
            throw new InvalidMasterPasswordException("Incorrect master password or corrupted data (invalid authentication tag)");
        } catch (Exception e) {
            throw new CryptoException("Critical failure when attempting to decrypt payload", e);
        }
    }

}
