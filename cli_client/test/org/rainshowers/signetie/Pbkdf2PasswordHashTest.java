/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie;

import org.junit.Test;
import static org.junit.Assert.*;

/**
 *
 * @author Ken Galle <ken@rainshowers.org>
 */
public class Pbkdf2PasswordHashTest {
    
    public Pbkdf2PasswordHashTest() {
    }

    @Test
    public void testGenerateHash_password() throws Exception {
        System.out.println("generateHash");
        
        String password = "Super big secret";
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        // Have to push the salt value in to get repeatable results
        
        byte[] testSalt = new byte[32];  // zero filled
        byte[] testHash = new byte[32];  // zero filled
        int iterationCount = 100000;

        PasswordHash expResult = 
                new PasswordHash(testHash, testSalt, iterationCount);
        
        PasswordHash result = instance.generateHash(password);
        /* 
         * The expected salt and hash can't be known. There is an infinetly 
         * small chance that the generated hash and salt will actually be the
         * expected value (all zeros), so check that it doesn't match.
         */
        assertNotEquals(expResult, result);
        verifyDuration(instance);
        System.out.println(String.format(
                "Duration (msec): %d", instance.getDuration()));
    }

    private static final byte[] VALID_TEST_SALT = new byte[] {113, 5, -97, -55, -24, 42, 113, 
        -2, 21, 45, -107, -15, 112, -26, 52, 64, 46, -122, -98, 56, 74, -14, 25, 113, -104, 
        50, 41, 127, 101, -93, 53, -51};
    
    private static final byte[] VALID_TEST_HASH = new byte[] {115, 28, -34, -64, -29, -82, 
        108, -81, 15, -6, -28, -57, 11, -90, 16, -49, -57, 77, 76, 59, -24, -37, 92, -122, 
        72, 70, -9, -106, -42, -52, -98, -61};
    
    @Test
    public void testGenerateHash_password_salt() throws Exception {
        System.out.println("generateHash");
        
        String password = "Super big secret";
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        int iterationCount = 100000;

        // Have to push the salt value in to get repeatable results
        PasswordHash expResult = new PasswordHash(
                VALID_TEST_HASH, VALID_TEST_SALT, iterationCount);
        
        PasswordHash result = instance.generateHash(password, VALID_TEST_SALT);
        /*
         * Here I specify the salt value, so I can determine ahead what the 
         * resulting hash should be.
         */
        assertEquals(expResult, result);
        verifyDuration(instance);
        System.out.println(String.format(
                "Duration (msec): %d", instance.getDuration()));
    }

    @Test
    public void testVerifyHash() throws Exception {
        System.out.println("verifyHash");
        String password = "Super big secret";

        int iterationCount = 100000;
        PasswordHash passwordHash = new PasswordHash(
                VALID_TEST_HASH, VALID_TEST_SALT, iterationCount);

        boolean expResult = true;
        boolean result = Pbkdf2PasswordHasher.verifyHash(password, passwordHash);
        assertEquals(expResult, result);
    }

    private void verifyDuration(Pbkdf2PasswordHasher instance) {
        /* 
         * Verify that the time taken is somewhere between 0.1 seconds 
         * and 1 seconds.
         */
        assertTrue((instance.getDuration() > 100) && 
                   (instance.getDuration() < 1000));
    }
   
}
