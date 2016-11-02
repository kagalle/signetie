/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package org.rainshowers.signetie;

import java.security.InvalidKeyException;
import java.security.KeyPair;
import java.security.NoSuchAlgorithmException;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;

/**
 *
 * @author ken
 */
public class RsaPkiEncryption {

    private static final String RSA_ALGORITHM_NAME = "RSA";
    private static final String ENCRYPT_WITH_PUBLIC_KEY_ERROR
            = "Unable to encrypt message with RSA public key.";

    /**
     * Used the public key in the supplied RSA keypair to encrypt the given message.
     */
    public static byte[] encryptWithPublicKey(KeyPair keyPair, String messageToEncrypt)
            throws SignetieException {
        byte[] encryptedMessage = null;
        try {
            Cipher rsaEncryptCipher;
            rsaEncryptCipher = Cipher.getInstance(RSA_ALGORITHM_NAME); // max 501 bytes can be encrypted
            rsaEncryptCipher.init(Cipher.ENCRYPT_MODE, keyPair.getPublic());
            encryptedMessage = rsaEncryptCipher.doFinal(messageToEncrypt.getBytes());
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException |
                IllegalBlockSizeException | BadPaddingException ex) {
            throw new SignetieException(ENCRYPT_WITH_PUBLIC_KEY_ERROR, ex);
        }
        return encryptedMessage;
    }

    static String decryptWithPrivateKey(KeyPair keyPair, byte[] encryptedMessage) {
        String decryptedMessage = null;
        try {
            Cipher rsaDecryptCipher;
            rsaDecryptCipher = Cipher.getInstance("RSA");
            rsaDecryptCipher.init(Cipher.DECRYPT_MODE, keyPair.getPrivate());
            decryptedMessage = new String(rsaDecryptCipher.doFinal(encryptedMessage));
        } catch (NoSuchAlgorithmException | NoSuchPaddingException | InvalidKeyException | IllegalBlockSizeException | BadPaddingException ex) {
            Logger.getLogger(RsaPkiEncryption.class.getName()).log(Level.SEVERE, null, ex);
        }
        return decryptedMessage;

    }

}
