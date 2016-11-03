/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.pbe;

import javax.crypto.SecretKey;
import javax.xml.bind.DatatypeConverter;

/**
 * Stores all of the necessary pieces of a password hash such that it can be
 * later verified against a password.
 * 
 * @author Ken Galle <ken@rainshowers.org>
 */
public class PasswordHash {
    private SecretKey key;
    private PasswordHashParams passwordHashParams;

    public PasswordHash() {
    }

    public PasswordHash(SecretKey key, PasswordHashParams passwordHashParams) {
        this.key = key;
        this.passwordHashParams = passwordHashParams;
    }

    public SecretKey getKey() {
        return key;
    }

    public void setKey(SecretKey key) {
        this.key = key;
    }

    public PasswordHashParams getPasswordHashParams() {
        return passwordHashParams;
    }

    public void setPasswordHashParams(PasswordHashParams passwordHashParams) {
        this.passwordHashParams = passwordHashParams;
    }

    /**
     * Get the current setting for the number of bytes to be returned in the 
     * output hash.
     */
    public int getHashLength() {
        return key.getEncoded().length;
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
