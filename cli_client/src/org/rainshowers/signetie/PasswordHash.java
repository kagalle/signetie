/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie;

import java.util.Arrays;
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
    private byte[] salt;
    private int interationCount;

    public PasswordHash() {
    }

    public PasswordHash(SecretKey key, byte[] salt, int interationCount) {
        this.key = key;
        this.salt = salt;
        this.interationCount = interationCount;
    }

    public SecretKey getKey() {
        return key;
    }

    public void setKey(SecretKey key) {
        this.key = key;
    }

    public byte[] getSalt() {
        return salt;
    }

    public void setSalt(byte[] salt) {
        this.salt = salt;
    }

    public int getInterationCount() {
        return interationCount;
    }

    public void setInterationCount(int interationCount) {
        this.interationCount = interationCount;
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

    // Netbeans generated
    @Override
    public int hashCode() {
        int hash = 3;
        hash = 71 * hash + Arrays.hashCode(this.key.getEncoded());
        hash = 71 * hash + Arrays.hashCode(this.salt);
        hash = 71 * hash + this.interationCount;
        return hash;
    }

    // Netbeans generated
    @Override
    public boolean equals(Object obj) {
        if (this == obj) {
            return true;
        }
        if (obj == null) {
            return false;
        }
        if (getClass() != obj.getClass()) {
            return false;
        }
        final PasswordHash other = (PasswordHash) obj;
        if (this.interationCount != other.interationCount) {
            return false;
        }
        if (!Arrays.equals(this.key.getEncoded(), other.key.getEncoded())) {
            return false;
        }
        if (!Arrays.equals(this.salt, other.salt)) {
            return false;
        }
        return true;
    }
    


}
