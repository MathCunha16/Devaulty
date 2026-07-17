package com.devaulty.backend.adapter.out.crypto;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.test.util.ReflectionTestUtils;

import javax.crypto.SecretKey;
import javax.crypto.spec.SecretKeySpec;
import java.time.Duration;
import java.time.Instant;

import static org.junit.jupiter.api.Assertions.*;

class MasterKeySessionHolderTest {

    private MasterKeySessionHolder sessionHolder;

    @BeforeEach
    void setUp() {
        sessionHolder = new MasterKeySessionHolder();
    }

    @Test
    void shouldInitializeWithNoKey() {
        assertNull(sessionHolder.getKey());
        assertFalse(sessionHolder.hasKey());
    }

    @Test
    void shouldSetAndGetKey() {
        // Arrange
        byte[] keyBytes = new byte[32];
        keyBytes[0] = 5;
        SecretKey secretKey = new SecretKeySpec(keyBytes, "AES");

        // Act
        sessionHolder.setKey(secretKey);

        // Assert
        assertTrue(sessionHolder.hasKey());
        assertSame(secretKey, sessionHolder.getKey());

        Instant lastActivity = (Instant) ReflectionTestUtils.getField(sessionHolder, "lastActivityAt");
        assertNotNull(lastActivity);
    }

    @Test
    void shouldClearSessionSecurely() {
        // Arrange
        byte[] keyBytes = new byte[]{1, 2, 3, 4};
        SecretKey secretKey = new SecretKeySpec(keyBytes, "AES");
        sessionHolder.setKey(secretKey);

        byte[] rawKeyBytesInHolder = (byte[]) ReflectionTestUtils.getField(sessionHolder, "rawKeyBytes");
        assertNotNull(rawKeyBytesInHolder);
        assertEquals(1, rawKeyBytesInHolder[0]);

        // Act
        sessionHolder.clear();

        // Assert
        assertFalse(sessionHolder.hasKey());
        assertNull(sessionHolder.getKey());
        assertNull(ReflectionTestUtils.getField(sessionHolder, "rawKeyBytes"));
        assertNull(ReflectionTestUtils.getField(sessionHolder, "lastActivityAt"));

        // Verify the original raw bytes inside the memory array were securely zeroed out!
        for (byte b : rawKeyBytesInHolder) {
            assertEquals(0, b);
        }
    }

    @Test
    void getKey_shouldNotTouch_whenNoKeyIsPresent() {
        assertNull(ReflectionTestUtils.getField(sessionHolder, "lastActivityAt"));
        sessionHolder.getKey();
        assertNull(ReflectionTestUtils.getField(sessionHolder, "lastActivityAt"));
    }

    @Test
    void getSecondsRemaining_shouldReturnZero_whenNoKeyIsPresent() {
        Long remaining = sessionHolder.getSecondsRemaining(Duration.ofMinutes(15));
        assertEquals(0L, remaining);
    }

    @Test
    void getSecondsRemaining_shouldReturnZero_whenLastActivityIsNull() {
        SecretKey secretKey = new SecretKeySpec(new byte[32], "AES");
        sessionHolder.setKey(secretKey);
        ReflectionTestUtils.setField(sessionHolder, "lastActivityAt", null);

        Long remaining = sessionHolder.getSecondsRemaining(Duration.ofMinutes(15));
        assertEquals(0L, remaining);
    }

    @Test
    void getSecondsRemaining_shouldReturnCorrectValue_whenActive() {
        SecretKey secretKey = new SecretKeySpec(new byte[32], "AES");
        sessionHolder.setKey(secretKey);

        // Set last activity to 5 minutes ago
        Instant fiveMinutesAgo = Instant.now().minus(Duration.ofMinutes(5));
        ReflectionTestUtils.setField(sessionHolder, "lastActivityAt", fiveMinutesAgo);

        // 15 minutes timeout - 5 minutes elapsed = 10 minutes (600 seconds) remaining
        Long remaining = sessionHolder.getSecondsRemaining(Duration.ofMinutes(15));
        
        // Assert remaining is roughly 600 seconds (allowing 2 seconds of tolerance)
        assertTrue(remaining >= 598 && remaining <= 600);
    }

    @Test
    void getSecondsRemaining_shouldReturnZero_whenTimeIsExpired() {
        SecretKey secretKey = new SecretKeySpec(new byte[32], "AES");
        sessionHolder.setKey(secretKey);

        // Set last activity to 20 minutes ago
        Instant twentyMinutesAgo = Instant.now().minus(Duration.ofMinutes(20));
        ReflectionTestUtils.setField(sessionHolder, "lastActivityAt", twentyMinutesAgo);

        Long remaining = sessionHolder.getSecondsRemaining(Duration.ofMinutes(15));
        assertEquals(0L, remaining);
    }

    @Test
    void shouldUpdateActivityTimeOnTouch() {
        SecretKey secretKey = new SecretKeySpec(new byte[32], "AES");
        sessionHolder.setKey(secretKey);

        Instant initialActivity = (Instant) ReflectionTestUtils.getField(sessionHolder, "lastActivityAt");
        assertNotNull(initialActivity);

        // Set back to make touch noticeable
        ReflectionTestUtils.setField(sessionHolder, "lastActivityAt", initialActivity.minus(Duration.ofMinutes(10)));

        // Act
        sessionHolder.touch();

        // Assert
        Instant newActivity = (Instant) ReflectionTestUtils.getField(sessionHolder, "lastActivityAt");
        assertNotNull(newActivity);
        assertTrue(newActivity.isAfter(initialActivity.minus(Duration.ofMinutes(10))));
    }
}
