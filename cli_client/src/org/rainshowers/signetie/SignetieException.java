/*
 * Copyright 2016 Kenneth Galle.  All rights reserved.
 * Use of this software is governed by a GPL v2 license
 * that can be found in the LICENSE file.
 */
package org.rainshowers.signetie;

/**
 * Exception specific to the running of Signetie.
 * The is the default Netbeans generated class.
 * 
 * @author Ken Galle <ken@rainshowers.org>
 */
public class SignetieException extends Exception {

    /**
     * Creates a new instance of <code>SignetieException</code> without detail
     * message.
     */
    public SignetieException() {
    }

    /**
     * Constructs an instance of <code>SignetieException</code> with the
     * specified detail message.
     *
     * @param msg the detail message.
     */
    public SignetieException(String msg) {
        super(msg);
    }

    /**
     * Constructs an instance of <code>SignetieException</code>, wrapping a
     * causing Exception, with the specified detail message.
     */
    public SignetieException(String string, Throwable thrwbl) {
        super(string, thrwbl);
    }
    
}
