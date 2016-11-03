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
    private  byte[] encryptedPrivateKey;
    private IvParameterSpec initializationVector;

    public SymmetricResult() {
    }

    public SymmetricResult(byte[] encryptedPrivateKey, IvParameterSpec initializationVector) {
        this.encryptedPrivateKey = encryptedPrivateKey;
        this.initializationVector = initializationVector;
    }

    public byte[] getEncryptedPrivateKey() {
        return encryptedPrivateKey;
    }

    public void setEncryptedPrivateKey(byte[] encryptedPrivateKey) {
        this.encryptedPrivateKey = encryptedPrivateKey;
    }

    public IvParameterSpec getInitializationVector() {
        return initializationVector;
    }

    public void setInitializationVector(IvParameterSpec initializationVector) {
        this.initializationVector = initializationVector;
    }
    
}
