/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.asymmetric;

import java.security.Key;
import java.security.KeyPair;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import static org.junit.Assert.*;
import org.rainshowers.signetie.SignetieException;
import org.rainshowers.signetie.pbe.PasswordHash;
import org.rainshowers.signetie.pbe.Pbkdf2PasswordHasher;
import org.rainshowers.signetie.symmetric.SymmetricEncryption;
import org.rainshowers.signetie.symmetric.SymmetricResult;

/**
 *
 * @author Ken Galle <ken@rainshowers.org>
 */
public class RsaPkiEncryptionTest {
    
    private KeyPair keyPair;
    
    public RsaPkiEncryptionTest() {
    }
    
    @Before
    public void setUp() {
        try {
            // generate key pair
            // This is done in setup because it is a very slow process.
            keyPair = RsaPkiEncryption.generateKeyPair();
        } catch (SignetieException ex) {
            Logger.getLogger(RsaPkiEncryptionTest.class.getName()).log(Level.SEVERE, "setUp() failed", ex);
        }
    }
    
    @After
    public void tearDown() {
    }

    @Test
    public void testEncryptDecryptWithPkiKeyPairOnly() throws Exception {

        // Encrypt something with the public key.
        String message = "An important message to be sent PKI.";
        byte[] encryptedMessage = RsaPkiEncryption.encryptWithPkiKey(keyPair.getPublic(), message);
        // Decrypt with the matching private key and verify
        String decryptedMessage = RsaPkiEncryption.decryptWithPkiKey(keyPair.getPrivate(), encryptedMessage);
        // assert
        assertEquals(message, decryptedMessage);
    }

    @Test
    public void testEncryptUsingPublicKeyThenDecryptUsingEncryptedDecryptedPrivateKey() throws Exception {
        
        // Encrypt something with the public key.
        String message = "An important message to be sent PKI.";
        byte[] encryptedMessage = RsaPkiEncryption.encryptWithPkiKey(keyPair.getPublic(), message);

        // generate hash from password
        String password = "super_secret";
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        PasswordHash passwordHash = instance.generateHash(password);
        
        // Encrypt the private key
        SymmetricResult symmetricResult = SymmetricEncryption.encryptKey(keyPair.getPrivate(), passwordHash.getKey());

        // Extract the private key
        Key decryptedPrivateKey = SymmetricEncryption.decryptKey(symmetricResult, passwordHash.getKey());
        
        // Decrypt message again with en/decrypted version key and verify
        String decryptedMessage = RsaPkiEncryption.decryptWithPkiKey(decryptedPrivateKey, encryptedMessage);
        assertEquals(message, decryptedMessage);
    }
   
}
