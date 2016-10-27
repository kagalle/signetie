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
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;
import java.security.PrivateKey;
import java.security.SecureRandom;
import java.security.spec.InvalidKeySpecException;
import java.security.spec.InvalidParameterSpecException;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.crypto.BadPaddingException;
import javax.crypto.Cipher;
import javax.crypto.EncryptedPrivateKeyInfo;
import javax.crypto.IllegalBlockSizeException;
import javax.crypto.NoSuchPaddingException;
import javax.crypto.SecretKey;
import javax.crypto.SecretKeyFactory;
import javax.crypto.spec.PBEKeySpec;
import javax.crypto.spec.PBEParameterSpec;

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
      KeyPair keyPair = generateKeyPair();
    }

    // based on http://stackoverflow.com/a/6164414/3728147
    // http://stackoverflow.com/a/16586921/3728147
    // http://www.programcreek.com/java-api-examples/index.php?api=java.security.spec.PKCS8EncodedKeySpec
    // http://snipplr.com/view/18368/
    // http://stackoverflow.com/questions/19640735/load-public-key-data-from-file
    // http://anandsekar.github.io/exporting-the-private-key-from-a-jks-keystore/
    // http://www.javamex.com/tutorials/cryptography/rsa_key_length.shtml
    // http://www.javamex.com/tutorials/cryptography/rsa_encryption.shtml
    private static void createPkcs8() throws NoSuchAlgorithmException, IOException, IllegalBlockSizeException, BadPaddingException, NoSuchPaddingException, InvalidKeySpecException, InvalidKeyException, InvalidAlgorithmParameterException, InvalidParameterSpecException, SignetieException {
        // generate key pair
        KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
        keyPairGenerator.initialize(4096);
        KeyPair keyPair = keyPairGenerator.genKeyPair();

        // extract the encoded private key, this is an unencrypted PKCS#8 private key
//        byte[] encodedprivkey = keyPair.getPrivate().getEncoded();

        // We must use a PasswordBasedEncryption algorithm in order to encrypt the private key, you may use any common algorithm supported by openssl, you can check them in the openssl documentation http://www.openssl.org/docs/apps/pkcs8.html
//        String MYPBEALG = "PBEWithSHA1AndDESede";
        String password = "super_secret";

//        int count = 20;// hash iteration count
//        SecureRandom random = new SecureRandom();
//        byte[] salt = new byte[8];
//        random.nextBytes(salt);

        // Create key from the password needed to encrypt the private key
        Pbkdf2PasswordHasher instance = new Pbkdf2PasswordHasher();
        SecretKey pbeKey = instance.generateHash(password).getKey();

        // encrypt the private key
        Cipher c = Cipher.getInstance("RSA");
        c.init(Cipher.WRAP_MODE, pbeKey);
        byte[] encryptedPrivateKey = c.wrap(keyPair.getPrivate());
        
        
        
        
        

        // Now construct  PKCS #8 EncryptedPrivateKeyInfo object
//        AlgorithmParameters algparms = AlgorithmParameters.getInstance(MYPBEALG);
//        algparms.init(pbeParamSpec);
//        EncryptedPrivateKeyInfo encinfo = new EncryptedPrivateKeyInfo(algparms, ciphertext);

        // and here we have it! a DER encoded PKCS#8 encrypted key!
//        byte[] encryptedPkcs8 = encinfo.getEncoded();
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
    private void x() throws SignetieException, IOException {
      KeyPair keyPair = Signetie.generateKeyPair();
      PrivateKey privateKey = keyPair.getPrivate();
EncryptedPrivateKeyInfo encryptPKInfo = new EncryptedPrivateKeyInfo(privateKey.getEncoded());

Cipher cipher = Cipher.getInstance(encryptPKInfo.getAlgName());
PBEKeySpec pbeKeySpec = new PBEKeySpec(passwd.toCharArray());
SecretKeyFactory secFac = SecretKeyFactory.getInstance(encryptPKInfo.getAlgName());
Key pbeKey = secFac.generateSecret(pbeKeySpec);

AlgorithmParameters algParams = encryptPKInfo.getAlgParameters();
cipher.init(Cipher.DECRYPT_MODE, pbeKey, algParams);
KeySpec pkcs8KeySpec = encryptPKInfo.getKeySpec(cipher);
KeyFactory kf = KeyFactory.getInstance("RSA");
return kf.generatePrivate(pkcs8KeySpec);      
    }
}
