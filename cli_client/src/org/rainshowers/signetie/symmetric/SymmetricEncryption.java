/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.symmetric;

import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.Key;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.NoSuchAlgorithmException;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.SecretKey;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import org.rainshowers.signetie.SignetieException;

/**
 * Routines that involve use of symmetric encryption (even if the target of what is being
 * encrypted in a private key, for example).
 * 
 * @author ken
 */
public class SymmetricEncryption {

    private static final String ENCRYPT_PRIVATE_KEY_ERROR = "Error encrypting private key.";
    /**
     * Encrypt the private key in the supplied keyPair and return the encrypted private key.
     */
    // TODO: this has to return the salt as well (in a object that hopefully already exists
    // in the crypto library, so that the salt can be supplied later on when decrypting
    // the private key.
    public static SymmetricResult encryptPrivateKey(KeyPair keyPair, SecretKey pbeKey) throws SignetieException {
        
        // encrypt the private key
        //     create a key in the specific form needed for AES, based on the key value in the pbeKey
        //     http://stackoverflow.com/a/13770749/3728147

        byte[] encryptedPrivateKey = null;
        IvParameterSpec initializationVector = null;
        try {
            SecretKeySpec aesKeySpec = new SecretKeySpec(pbeKey.getEncoded(), "AES");
            Cipher aesPbeEncryptCipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            // changing the IV size will result in an exception
            byte[] ivBytes = new byte[aesPbeEncryptCipher.getBlockSize()];
            initializationVector = new IvParameterSpec(ivBytes);
            aesPbeEncryptCipher.init(Cipher.ENCRYPT_MODE, aesKeySpec, initializationVector);  // pbeKey has to match DESede - somehow
            encryptedPrivateKey = aesPbeEncryptCipher.doFinal(keyPair.getPrivate().getEncoded());
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException | 
                InvalidAlgorithmParameterException | IllegalBlockSizeException | BadPaddingException ex) {
            throw new SignetieException(ENCRYPT_PRIVATE_KEY_ERROR, ex);
        }
        SymmetricResult symmetricResult = new SymmetricResult(encryptedPrivateKey, initializationVector);
        return symmetricResult;
    }

    private static final String DECRYPT_PRIVATE_KEY_ERROR = "Error decrypting private key.";
    
    public static Key decryptPrivateKey(KeyPair keyPair, SymmetricResult symmetricResult, SecretKey pbeKey) throws SignetieException {
        Cipher aesPbeDecryptCipher;
        Key decryptedPrivateKey = null;
        try {
            SecretKeySpec aesKeySpec = new SecretKeySpec(pbeKey.getEncoded(), "AES");
            aesPbeDecryptCipher = Cipher.getInstance("AES/CBC/PKCS5Padding"); // was RSA, but that doesn't make sense - why use an asymetric key here?
            // https://community.oracle.com/thread/1528052?start=0
            aesPbeDecryptCipher.init(Cipher.DECRYPT_MODE, aesKeySpec, symmetricResult.getInitializationVector());
            byte[] decryptedPrivateKeyData = aesPbeDecryptCipher.doFinal(symmetricResult.getEncryptedPrivateKey());
//            decryptedPrivateKey = new PrivateKey(decryptedPrivateKeyData, "RSA");

            // http://stackoverflow.com/a/8455164/3728147
            KeyFactory keyFactory = KeyFactory.getInstance("RSA");
            KeySpec privateKeySpec = new PKCS8EncodedKeySpec(decryptedPrivateKeyData);
            decryptedPrivateKey = keyFactory.generatePrivate(privateKeySpec);

        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException | 
                InvalidAlgorithmParameterException | IllegalBlockSizeException | 
                BadPaddingException | InvalidKeySpecException ex) {
            throw new SignetieException(DECRYPT_PRIVATE_KEY_ERROR, ex);
        }
        return decryptedPrivateKey;
    }
    
}
