/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie;

import org.rainshowers.signetie.asymmetric.RsaPkiEncryption;
import org.rainshowers.signetie.symmetric.SymmetricEncryption;
import java.security.Key;
import java.security.KeyPair;
import java.util.logging.Level;
import java.util.logging.Logger;
import org.rainshowers.signetie.pbe.PasswordHash;
import org.rainshowers.signetie.pbe.Pbkdf2PasswordHasher;
import org.rainshowers.signetie.symmetric.SymmetricResult;

/**
 *
 * @author Ken Galle <ken@rainshowers.org>
 */
public class Signetie {

    /**
     * @param args the command line arguments
     */
    public static void main(String[] args) throws SignetieException {

        try {
            Signetie.createPkcs8();
        } catch (SignetieException ex) {
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
    private static void createPkcs8() throws SignetieException {
        // generate key pair
        KeyPair keyPair = RsaPkiEncryption.generateKeyPair();

        // Encrypt something with the public key.
        String message = "An important message to be sent PKI.";
        byte[] encryptedMessage = RsaPkiEncryption.encryptWithPkiKey(keyPair.getPublic(), message);

        // Decrypt with the matching private key and verify
        String firstDecryptedMessage = RsaPkiEncryption.decryptWithPkiKey(keyPair.getPrivate(), encryptedMessage);
        System.out.println( (message.equals(firstDecryptedMessage) ? "yes" : "no") );

        // generate hash from password
        String password = "super_secret";
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        PasswordHash passwordHash = instance.generateHash(password);
        
        // Encrypt the private key
        SymmetricResult symmetricResult = SymmetricEncryption.encryptKey(keyPair.getPrivate(), passwordHash.getKey());

        // Extract the private key
        Key decryptedPrivateKey = SymmetricEncryption.decryptKey(symmetricResult, passwordHash.getKey());
        
        // Decrypt message again with en/decrypted version key and verify
        String secondDecryptedMessage = RsaPkiEncryption.decryptWithPkiKey(decryptedPrivateKey, encryptedMessage);
        System.out.println( (message.equals(secondDecryptedMessage) ? "yes" : "no") );
        
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
