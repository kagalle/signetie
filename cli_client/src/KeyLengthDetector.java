import javax.crypto.Cipher;
import java.security.NoSuchAlgorithmException;

// http://derjan.io/blog/2013/03/15/nevermind-jce-unlimited-strength-use-openjdk/
// https://gist.github.com/jehrhardt/5167854
public class KeyLengthDetector {
  public static void main(String[] args) {
    int allowedKeyLength = 0;

    try {
      allowedKeyLength = Cipher.getMaxAllowedKeyLength("AES");
    } catch (NoSuchAlgorithmException e) {
      e.printStackTrace();
    }

    System.out.println("The allowed key length for AES is: " + allowedKeyLength);
  }
}