/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.symmetric;

import javax.crypto.spec.IvParameterSpec;

/**
 *
 * @author Ken Galle <ken@rainshowers.org>
 */
public class SymmetricResult {
    private  byte[] encryptedKey;
    private IvParameterSpec initializationVector;

    public SymmetricResult() {
    }

    public SymmetricResult(byte[] encryptedKey, IvParameterSpec initializationVector) {
        this.encryptedKey = encryptedKey;
        this.initializationVector = initializationVector;
    }

    public byte[] getEncryptedKey() {
        return encryptedKey;
    }

    public void setEncryptedKey(byte[] encryptedKey) {
        this.encryptedKey = encryptedKey;
    }

    public IvParameterSpec getInitializationVector() {
        return initializationVector;
    }

    public void setInitializationVector(IvParameterSpec initializationVector) {
        this.initializationVector = initializationVector;
    }
    
}
