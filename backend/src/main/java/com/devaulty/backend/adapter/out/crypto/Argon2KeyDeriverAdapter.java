package com.devaulty.backend.adapter.out.crypto;

import com.devaulty.backend.application.port.out.security.KeyDerivationPort;
import org.bouncycastle.crypto.generators.Argon2BytesGenerator;
import org.bouncycastle.crypto.params.Argon2Parameters;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.nio.CharBuffer;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;

@Component
public class Argon2KeyDeriverAdapter implements KeyDerivationPort {

    /**
     * Derives a secure 256-bit AES key from the master password using the memory-hard Argon2id KDF.
     * <p>
     * This method configures Argon2id with 3 iterations, 64MB of memory, and 2 parallel threads
     * to significantly elevate the computational cost of brute-force and hardware-accelerated (GPU/ASIC) attacks.
     * The resulting key bytes are safely encapsulated into a {@link SecretKey} and dispatched to the
     * {@link MasterKeySessionHolder} to back subsequent cryptographic operations without repetitive key re-derivation.
     * </p>
     *
     * @param password the plain-text master password provided by the user
     * @param salt     the cryptographically secure global application salt
     */

    @Override
    public SecretKey deriveKey(char[] password, byte[] salt) {
        Argon2BytesGenerator gen = new Argon2BytesGenerator();
        Argon2Parameters.Builder builder = new Argon2Parameters.Builder(Argon2Parameters.ARGON2_id)
                .withVersion(Argon2Parameters.ARGON2_VERSION_13)
                .withIterations(3)
                .withMemoryAsKB(65536)
                .withParallelism(2)
                .withSalt(salt);

        gen.init(builder.build());
        byte[] keyBytes = new byte[32];

        byte[] passwordBytes = StandardCharsets.UTF_8.encode(CharBuffer.wrap(password)).array();

        try {
            gen.generateBytes(passwordBytes, keyBytes, 0, keyBytes.length);
            return new SecretKeySpec(keyBytes, "AES");
        } finally {
            Arrays.fill(passwordBytes, (byte) 0);
        }
    }
}