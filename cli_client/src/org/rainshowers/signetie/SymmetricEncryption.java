/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.rainshowers.signetie;

import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.Key;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.NoSuchAlgorithmException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.SecretKey;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

/**
 * Routines that involve use of symmetric encryption (even if the target of what is being
 * encrypted in a private key, for example).
 * 
 * @author ken
 */
class SymmetricEncryption {

    private static final String ENCRYPT_PRIVATE_KEY_ERROR = "Error encrypting private key.";
    /**
     * Encrypt the private key in the supplied keyPair and return the encrypted private key.
     */
    // TODO: this has to return the salt as well (in a object that hopefully already exists
    // in the crypto library, so that the salt can be supplied later on when decrypting
    // the private key.
    static byte[] encryptPrivateKey(KeyPair keyPair, String password) throws SignetieException {
        
        // Create key from the password needed to encrypt the private key
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        PasswordHash passwordHash = instance.generateHash(password);
        SecretKey pbeKey = passwordHash.getKey();

        // encrypt the private key
        //     create a key in the specific form needed for AES, based on the key value in the pbeKey
        //     http://stackoverflow.com/a/13770749/3728147

        byte[] encryptedPrivateKey = null;
        try {
            SecretKeySpec aesKeySpec = new SecretKeySpec(pbeKey.getEncoded(), "AES");
            Cipher aesPbeEncryptCipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            // changing the IV size will result in an exception
            byte[] ivBytes = new byte[aesPbeEncryptCipher.getBlockSize()];
            IvParameterSpec iv = new IvParameterSpec(ivBytes);
            aesPbeEncryptCipher.init(Cipher.ENCRYPT_MODE, aesKeySpec, iv);  // pbeKey has to match DESede - somehow
            encryptedPrivateKey = aesPbeEncryptCipher.doFinal(keyPair.getPrivate().getEncoded());
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException | InvalidAlgorithmParameterException | 
                IllegalBlockSizeException | BadPaddingException ex) {
            throw new SignetieException(ENCRYPT_PRIVATE_KEY_ERROR, ex);
        }
        return encryptedPrivateKey;
    }

    static Key decryptPrivateKey(KeyPair keyPair, String password) {
        Cipher aesPbeDecryptCipher;
        Key decryptedPrivateKey = null;
        try {
            aesPbeDecryptCipher = Cipher.getInstance("AES/CBC/PKCS5Padding"); // was RSA, but that doesn't make sense - why use an asymetric key here?
            // https://community.oracle.com/thread/1528052?start=0
            aesPbeDecryptCipher.init(Cipher.DECRYPT_MODE, aesKeySpec, iv);
            byte[] decryptedPrivateKeyData = aesPbeDecryptCipher.doFinal(encryptedPrivateKey);
//            decryptedPrivateKey = new PrivateKey(decryptedPrivateKeyData, "RSA");

            // http://stackoverflow.com/a/8455164/3728147
            KeyFactory keyFactory = KeyFactory.getInstance("RSA");
            KeySpec privateKeySpec = new PKCS8EncodedKeySpec(decryptedPrivateKeyData);
            decryptedPrivateKey = keyFactory.generatePrivate(privateKeySpec);

        } catch (NoSuchAlgorithmException ex) {
            Logger.getLogger(Signetie.class.getName()).log(Level.SEVERE, null, ex);
        }
    }
    
    private static SecretKeySpec createAesKeySpec() {
        return new SecretKeySpec(pbeKey.getEncoded(), "AES");
    }
    
}
