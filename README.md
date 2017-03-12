# signetie
Signetie is a electronic signature application that relies on the user's Google+ identity as the authority of the signature.

Project is under development.

Updated status 3/12/2017:
 

 * Java code under `cli_client` contains routines that will be needed to support the signing of documents.  Will need to be ported to golang.
 * Golang code under `client_golang` contains a working GTK+/golang application that performs the Google App Engine login and authentication for desktop applications.  The is the first necessary piece of the application.  The rest of the functionality will rely on this.


-----

Signetie is an application that allows you to digitally sign documents like text files and source-code files, linking them to your Google profile.  These documents can be shown to be signed by you, and the signature itself will be linked to the Signetie application.  The recepient of the document will be able to verify the document as yours by using Signetie to verify that the document's signature is linked to your Google account.

Traditionally, documents are digitally signed and verified using public-key infrastructure (PKI).  Documents are verified by using a certificate authority signed certificate containing the author's public key.  The certificate authority is the entity that assures that the signer is who the certificate says they are.  Signetie exsists because certificate authorities charge a significate amount of money for a certificate, and they only remain valid for relatively short period of time, usually one year.  This may not be practical for an open-source project, and is certainly a barier to the use of code-signing.

Signetie uses a similar approach, but instead of the certificate authority providing the assurance of identity, the signer's Signetie's account, linked to their Google account, provides some level of assurance of identity.  

As a Signetie user, Signetie's primary function is to sign your public key and link it to your Google account.  You create a public and private key-pair locally using any approprate means, including the SignetieDesktop application.  As a result, the Signetie application does not know or store your private key.  You can use the SignetieDesktop application to manage your key pair and sign documents.  Signetie provides the Signetie-signed public key and identity, known as a "signed-signetie-identity-record," or SSIR.

The receipient of the document can use this SSIR to validate the document as coming from you, or specifically, coming from the owner of the Google account the signature is associated with.

Signetie is alpha software and no assurance is provided to its suitability for use.  However, the project is open-source, so you are encouraged to examine the software to determine its suitability.  If you identify issues with Signetie, please send this as feedback [provide link].

Signetie is built on Google App Engine, and it relies on this platform for its functionality, ie. it is not useful outside of this environment.  There are costs associated with running the application, and as it is still alpha software, the mechinisms for making the project self-supporting have not been worked out.  If the service is unavailable, please make a support request [provide link].  Any assistance in this regard would be welcomed.

The intended purpose of Signetie is only that of signing/verifying documents and not for encrypting/decrypting documents.

Thank you for your interest in Signetie.

