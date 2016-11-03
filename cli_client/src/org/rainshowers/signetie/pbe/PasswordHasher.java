/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.pbe;

import org.rainshowers.signetie.SignetieException;

/**
 * Use some hashing method to generate a hash value of the specified length.
 * Implementations should provide a duration value in milliseconds for the time the
 * computation took so that the user can adjust the iteration count to obtain a balance
 * between security and usability.
 * 
 * @author Ken Galle <ken@rainshowers.org>
 */
public interface PasswordHasher {

    /**
     * Generate a hash value based on the password supplied and a random
     * salt value.
     */
    PasswordHash generateHash(String password) throws SignetieException;

    /**
     * Generate a hash value based on the password and salt values supplied.
     */
    PasswordHash generateHash(String password, PasswordHashParams passwordHashParams) throws SignetieException;

    /**
     * Retrieve the approximate time it took to generate the key.  The iteration
     * count should be adjusted so that any form of brute force attack on the
     * key is infeasible, say 500 ms.
     */
    long getDuration();

    /**
     * Get the current setting for the number of bytes to be returned in the
     * output hash.
     */
    int getHashLength();

    /**
     * Get the iteration count that would be used if generateHash() were called.
     */
    int getInterationCount();

    /**
     * Modify the default length of the output hash, in number of bytes.
     */
    void setHashLength(int hashLength);

    /**
     * Modify the default iteration count for the process.  The iteration count
     * effects the time required to generate the key, and should be reasonably
     * long as to prohibit brute force attacks.  Default value is reasonable
     * for native JVM execution.  See also #getDuration().
     */
    void setInterationCount(int interationCount);
    
}
