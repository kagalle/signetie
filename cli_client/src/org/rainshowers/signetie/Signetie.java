/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie;

import java.io.IOException;
import java.security.AlgorithmParameters;
import java.security.InvalidAlgorithmParameterException;
import java.security.InvalidKeyException;
import java.security.Key;
import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.InvalidParameterSpecException;
import java.security.spec.KeySpec;
import java.security.spec.PKCS8EncodedKeySpec;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.EncryptedPrivateKeyInfo;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.PBEKeySpec;
import javax.crypto.spec.PBEParameterSpec;
import javax.crypto.spec.SecretKeySpec;

/**
 *
 * @author Ken Galle <ken@rainshowers.org>
 */
public class Signetie {

    private class SignetieKeyPair {

        private Key privateKey;
        public Key publicKey;

    }

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws SignetieException {
//      KeyPair keyPair = generateKeyPair();

        try {
            Signetie.createPkcs8();
        } catch (IOException | IllegalBlockSizeException |
                BadPaddingException | NoSuchPaddingException | InvalidKeySpecException |
                InvalidKeyException | InvalidAlgorithmParameterException |
                InvalidParameterSpecException ex) {
            Logger.getLogger(Signetie.class.getName()).log(Level.SEVERE, null, ex);
        }

    }

    // based on http://stackoverflow.com/a/6164414/3728147
    // http://stackoverflow.com/a/16586921/3728147
    // http://www.programcreek.com/java-api-examples/index.php?api=java.security.spec.PKCS8EncodedKeySpec
    // http://snipplr.com/view/18368/
    // http://stackoverflow.com/questions/19640735/load-public-key-data-from-file
    // http://anandsekar.github.io/exporting-the-private-key-from-a-jks-keystore/
    // http://www.javamex.com/tutorials/cryptography/rsa_key_length.shtml
    // http://www.javamex.com/tutorials/cryptography/rsa_encryption.shtml
    private static void createPkcs8() throws IOException, IllegalBlockSizeException, BadPaddingException, NoSuchPaddingException, InvalidKeySpecException, InvalidKeyException, InvalidAlgorithmParameterException, InvalidParameterSpecException, SignetieException {
        // generate key pair
        KeyPair keyPair = Signetie.generateKeyPair();

        // Encrypt something with the public key.
        String message = "An important message to be sent PKI.";
        byte[] encryptedMessage = RsaPkiEncryption.encryptWithPublicKey(keyPair, message);

        // Decrypt and verify
        String firstDecryptedMessage = RsaPkiEncryption.decryptWithPrivateKey(keyPair, encryptedMessage);
        System.out.println( (message.equals(firstDecryptedMessage) ? "yes" : "no") );

        // Encrypt the private key
        
        String password = "super_secret";
        byte[] encryptedPrivateKey = AsymmetricEncryption.encryptPrivateKey(keyPair);
        
        // Create key from the password needed to encrypt the private key
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        PasswordHash passwordHash = instance.generateHash(password);
        SecretKey pbeKey = passwordHash.getKey();

        // encrypt the private key
//        Cipher aesPbeEncryptCipher = Cipher.getInstance("AES"); // was RSA, but that doesn't make sense - why use an asymetric key here?
        //     create a key in the specific form needed for AES, based on the key value in the pbeKey
        //     http://stackoverflow.com/a/13770749/3728147
        SecretKeySpec aesKeySpec = new SecretKeySpec(pbeKey.getEncoded(), "AES");

        Cipher aesPbeEncryptCipher;
        byte[] encryptedPrivateKey = {0x00};
        IvParameterSpec iv = null;
        try {
            aesPbeEncryptCipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            byte[] ivBytes = new byte[aesPbeEncryptCipher.getBlockSize()];
            iv = new IvParameterSpec(ivBytes);
//        aesPbeEncryptCipher.init(Cipher.ENCRYPT_MODE, pbeKey, iv);        
            aesPbeEncryptCipher.init(Cipher.ENCRYPT_MODE, aesKeySpec, iv);  // pbeKey has to match DESede - somehow
            encryptedPrivateKey = aesPbeEncryptCipher.doFinal(keyPair.getPrivate().getEncoded());
        } catch (NoSuchAlgorithmException ex) {
            Logger.getLogger(Signetie.class.getName()).log(Level.SEVERE, null, ex);
        }
        // changing the IV size will result in an exception

        // Extract the private key
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

        // Decrypt message again with en/decrypted version key and verify
        Cipher rsaDecryptCipher2;
        try {
            rsaDecryptCipher2 = Cipher.getInstance("RSA");
            rsaDecryptCipher2.init(Cipher.DECRYPT_MODE, decryptedPrivateKey);
            String secondDecryptedMessage = new String(rsaDecryptCipher2.doFinal(encryptedMessage));
            if (message.equals(secondDecryptedMessage)) {
                System.out.println("yes");
            }
        } catch (NoSuchAlgorithmException ex) {
            Logger.getLogger(Signetie.class.getName()).log(Level.SEVERE, null, ex);
        }

    }
    private static final String GENERATE_KEYPAIR_ERROR = "Error generating key pair.";

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

    // http://stackoverflow.com/a/3985508/3728147
//    private void readOutPkcs8EncryptedKey() throws SignetieException, IOException {
//      KeyPair keyPair = Signetie.generateKeyPair();
//      PrivateKey privateKey = keyPair.getPrivate();
//EncryptedPrivateKeyInfo encryptPKInfo = new EncryptedPrivateKeyInfo(algParams, encryptedData);
//
//Cipher cipher = Cipher.getInstance(encryptPKInfo.getAlgName());
//PBEKeySpec pbeKeySpec = new PBEKeySpec(passwd.toCharArray());
//SecretKeyFactory secFac = SecretKeyFactory.getInstance(encryptPKInfo.getAlgName());
//Key pbeKey = secFac.generateSecret(pbeKeySpec);
//
//AlgorithmParameters algParams = encryptPKInfo.getAlgParameters();
//cipher.init(Cipher.DECRYPT_MODE, pbeKey, algParams);
//KeySpec pkcs8KeySpec = encryptPKInfo.getKeySpec(cipher);
//KeyFactory kf = KeyFactory.getInstance("RSA");
//return kf.generatePrivate(pkcs8KeySpec);      
//    }
}
