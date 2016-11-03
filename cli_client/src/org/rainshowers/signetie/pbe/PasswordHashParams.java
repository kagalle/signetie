/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.pbe;

import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import javax.xml.bind.DatatypeConverter;
import org.rainshowers.signetie.SignetieException;

/**
 * Stores all of the necessary pieces of a password hash such that it can be
 * later verified against a password.
 * 
 * @author Ken Galle <ken@rainshowers.org>
 */
public class PasswordHashParams {
    private byte[] salt;
    private int iterationCount;
    private int hashLength;
    private String SecretKeyFactoryAlgorithm;

    private static final String RANDOM_ALGORITHM = "SHA1PRNG";
    private static final String GENERATE_SALT_ERROR = 
            "Error creating random salt for password hash.";
    
    public PasswordHashParams() throws SignetieException {
        SecretKeyFactoryAlgorithm = "PBKDF2WithHmacSHA1";
        iterationCount = 100000;
        hashLength = 32;
        
        // Generate a new salt value;
        salt = new byte[32];
        SecureRandom sr;
        try {
            sr = SecureRandom.getInstance(RANDOM_ALGORITHM);
        } catch (NoSuchAlgorithmException ex) {
            throw new SignetieException(GENERATE_SALT_ERROR, ex);
        }
        sr.nextBytes(salt);
    }

    public PasswordHashParams(byte[] salt) {
        SecretKeyFactoryAlgorithm = "PBKDF2WithHmacSHA1";
        iterationCount = 100000;
        hashLength = 32;
        this.salt = salt;
    }

    public PasswordHashParams(byte[] salt, int interationCount, int hashLength, String SecretKeyFactoryAlgorithm) {
        this.salt = salt;
        this.iterationCount = interationCount;
        this.hashLength = hashLength;
        this.SecretKeyFactoryAlgorithm = SecretKeyFactoryAlgorithm;
    }

    public byte[] getSalt() {
        return salt;
    }

    public void setSalt(byte[] salt) {
        this.salt = salt;
    }

    public int getInterationCount() {
        return iterationCount;
    }

    public void setInterationCount(int interationCount) {
        this.iterationCount = interationCount;
    }

    public int getHashLength() {
        return hashLength;
    }

    public void setHashLength(int hashLength) {
        this.hashLength = hashLength;
    }

    public String getSecretKeyFactoryAlgorithm() {
        return SecretKeyFactoryAlgorithm;
    }

    public void setSecretKeyFactoryAlgorithm(String SecretKeyFactoryAlgorithm) {
        this.SecretKeyFactoryAlgorithm = SecretKeyFactoryAlgorithm;
    }

    /**
     * Retrieve the generated salt value as a (printable) String.
     */
    public static String byteArrayToString(byte[] byteArray) {
        String result = "";  // default value if input is null
        // check for null, as IllegalArgumentException could be thrown
        if (byteArray != null) {
            result = DatatypeConverter.printHexBinary(byteArray);
        }
        return result;
    }

}
