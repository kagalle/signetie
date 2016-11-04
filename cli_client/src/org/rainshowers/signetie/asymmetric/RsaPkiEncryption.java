/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie.asymmetric;

import java.security.InvalidKeyException;
import java.security.Key;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import org.rainshowers.signetie.SignetieException;

// based on http://stackoverflow.com/a/6164414/3728147
// http://stackoverflow.com/a/16586921/3728147
// http://www.programcreek.com/java-api-examples/index.php?api=java.security.spec.PKCS8EncodedKeySpec
// http://snipplr.com/view/18368/
// http://stackoverflow.com/questions/19640735/load-public-key-data-from-file
// http://anandsekar.github.io/exporting-the-private-key-from-a-jks-keystore/
// http://www.javamex.com/tutorials/cryptography/rsa_key_length.shtml
// http://www.javamex.com/tutorials/cryptography/rsa_encryption.shtml
// http://stackoverflow.com/a/3985508/3728147

/**
 *
 * @author ken
 */
public class RsaPkiEncryption {

    private static final String GENERATE_KEYPAIR_ERROR = "Error generating key pair.";
    /**
     * Generate a new RSA keypair.
     */
    public static KeyPair generateKeyPair() throws SignetieException {
        try {
            // generate key pair
            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(4096);
            KeyPair keyPair = keyPairGenerator.genKeyPair();
            return keyPair;
        } catch (NoSuchAlgorithmException ex) {
            throw new SignetieException(GENERATE_KEYPAIR_ERROR, ex);
        }
    }

    private static final String RSA_ALGORITHM_NAME = "RSA";
    private static final String ENCRYPT_WITH_PUBLIC_KEY_ERROR
            = "Unable to encrypt message with RSA public key.";
    /**
     * Used the public key in the supplied RSA keypair to encrypt the given message.
     */
    public static byte[] encryptWithPkiKey(Key key, String messageToEncrypt)
            throws SignetieException {
        byte[] encryptedMessage = null;
        try {
            Cipher rsaEncryptCipher;
            rsaEncryptCipher = Cipher.getInstance(RSA_ALGORITHM_NAME); // max 501 bytes can be encrypted
            rsaEncryptCipher.init(Cipher.ENCRYPT_MODE, key);
            encryptedMessage = rsaEncryptCipher.doFinal(messageToEncrypt.getBytes());
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException |
                IllegalBlockSizeException | BadPaddingException ex) {
            throw new SignetieException(ENCRYPT_WITH_PUBLIC_KEY_ERROR, ex);
        }
        return encryptedMessage;
    }

    public static String decryptWithPkiKey(Key key, byte[] encryptedMessage) {
        String decryptedMessage = null;
        try {
            Cipher rsaDecryptCipher;
            rsaDecryptCipher = Cipher.getInstance("RSA");
            rsaDecryptCipher.init(Cipher.DECRYPT_MODE, key);
            decryptedMessage = new String(rsaDecryptCipher.doFinal(encryptedMessage));
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException | IllegalBlockSizeException | BadPaddingException ex) {
            Logger.getLogger(RsaPkiEncryption.class.getName()).log(Level.SEVERE, null, ex);
        }
        return decryptedMessage;

    }

}
