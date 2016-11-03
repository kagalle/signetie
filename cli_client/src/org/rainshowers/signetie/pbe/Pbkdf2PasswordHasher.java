/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.pbe;

import java.security.NoSuchAlgorithmException;
import java.security.spec.InvalidKeySpecException;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import org.rainshowers.signetie.SignetieException;



// http://stackoverflow.com/q/24652602/3728147
// https://adambard.com/blog/3-wrong-ways-to-store-a-password/
/**
 * Generates and stores a hash and salt value based on a supplied password using the
 * PBDF2 standard.
 * 
 * @author Ken Galle <ken@rainshowers.org>
 */
public class Pbkdf2PasswordHasher implements PasswordHasher {
    private static final String SECRET_KEY_FACTORY_ALGORITHM = "PBKDF2WithHmacSHA1";
    //timing: http://stackoverflow.com/a/180191/3728147
    private int iterationCount = 100000;
    private int hashLength = 32;
    private long duration;

    public Pbkdf2PasswordHasher() {
    }
    
    // https://en.wikipedia.org/wiki/PBKDF2
    /**
     * Generate a hash value based on the password supplied and a random
     * salt value.
     */
    @Override
    public PasswordHash generateHash(String password) throws SignetieException {
        // create default hash parameters to use
        PasswordHashParams passwordHashParams = new PasswordHashParams();
        // create the hash
        PasswordHash passwordHash = generateHash(password, passwordHashParams);
        return passwordHash;
    }

    private static final String GENERATE_HASH_ERROR = "Error generating password hash.";
    /**
     * Generate a hash value based on the password and salt values supplied.
     */
    @Override
    public PasswordHash generateHash(String password, PasswordHashParams passwordHashParams) 
            throws SignetieException {
        
        long startTime = System.nanoTime();
        try {
            SecretKey secretKey = generateKey(password, 
                    passwordHashParams.getSalt(), 
                    passwordHashParams.getInterationCount(), 
                    passwordHashParams.getHashLength());
            
            // track the time it took to compute
            long endTime = System.nanoTime();
            // divide by 1e6 to convert from nanoseconds to milliseconds.
            duration = (endTime - startTime) / 1000000;
            // bundle up the result to return
            PasswordHash passwordHash = new PasswordHash(secretKey, passwordHashParams);
            return passwordHash;
        } catch (NoSuchAlgorithmException | InvalidKeySpecException ex) {
            throw new SignetieException(GENERATE_HASH_ERROR, ex);
        }
    }
    private static final String VERIFY_HASH_ERROR = "Error verifying password hash.";
    /**
     * Verify that password resolves to the same hash in passwordHash.
     */
    public static boolean verifyHash(String password, PasswordHash passwordHash) 
            throws SignetieException {
        try {
            // generate a hash based on the supplied password
            SecretKey secretKey = generateKey(password, passwordHash.getPasswordHashParams().getSalt(),
                    passwordHash.getPasswordHashParams().getInterationCount(), passwordHash.getHashLength());
            // convert both old and new hashes to strings
            String oldHashString = PasswordHash.byteArrayToString(passwordHash.getKey().getEncoded());
            String newHashString = PasswordHash.byteArrayToString(secretKey.getEncoded());
            // return comparison
            return oldHashString.equals(newHashString);
        } catch (NoSuchAlgorithmException | InvalidKeySpecException ex) {
            throw new SignetieException(VERIFY_HASH_ERROR, ex);
        }
    }

    /**
     * The JCA/JCE internals of actually creating the hash.
     * @param password The password as a String.
     * @param salt
     * @param iterationCount
     * @param hashLength The hash length to generate, in number of bytes.
     * @return
     * @throws NoSuchAlgorithmException
     * @throws InvalidKeySpecException 
     */
    private static SecretKey generateKey(String password, byte[] salt, int iterationCount, 
            int hashLength) throws NoSuchAlgorithmException, InvalidKeySpecException {
        
        // PBEKeySpec requires the password to be in a byte array
        char[] passwordAsCharArray = password.toCharArray();
        // define what I want the key to end up looking like
        PBEKeySpec pbeKeySpec = new PBEKeySpec(passwordAsCharArray, salt, 
                iterationCount, (hashLength * 8 /* bits per byte */));
        // create the object that will make the key
        SecretKeyFactory secretKeyFactory = SecretKeyFactory.getInstance(SECRET_KEY_FACTORY_ALGORITHM);
        // generate the key
        SecretKey secretKey = secretKeyFactory.generateSecret(pbeKeySpec);

        // System.out.println(new String(secretKey.getEncoded()));
        System.out.println(secretKey.getFormat());
        System.out.println(secretKey.getAlgorithm());
        return secretKey;
    }
    
    /**
     * Get the iteration count that would be used if generateHash() were called.
     */
    @Override
    public int getInterationCount() {
        return iterationCount;
    }

    /**
     * Modify the default iteration count for the process.  The iteration count
     * effects the time required to generate the key, and should be reasonably
     * long as to prohibit brute force attacks.  Default value is reasonable
     * for native JVM execution.  See also #getDuration().
     */
    @Override
    public void setInterationCount(int interationCount) {
        this.iterationCount = interationCount;
    }

    /**
     * Get the current setting for the number of bytes to be returned in the 
     * output hash.
     */
    @Override
    public int getHashLength() {
        return hashLength;
    }

    /**
     * Modify the default length of the output hash, in number of bytes.
     */
    @Override
    public void setHashLength(int hashLength) {
        this.hashLength = hashLength;
    }

    /**
     * Retrieve the approximate time it took to generate the key.  The iteration
     * count should be adjusted so that any form of brute force attack on the 
     * key is infeasible, say 500 ms.
     */
    @Override
    public long getDuration() {
        return duration;
    }
    
}
