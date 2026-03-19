package com.github.onsdigital.dp.permissions.api.sdk.exception;

import org.apache.hc.core5.http.HttpStatus;

import lombok.Getter;

public class PermissionsAPIException extends Exception {
    /**
     * Status code of the error.
     */
    @Getter
    private final int code;

    /**
     * Create a new instance of a PermissionsAPIException.
     *
     * @param message    A string detailing the reason for the exception
     * @param statusCode The http status code that caused the API exception
     */
    public PermissionsAPIException(final String message, final int statusCode) {
        super(message);
        this.code = statusCode;
    }

    /**
     * New default constructor.
     */
    public PermissionsAPIException() {
        this.code = HttpStatus.SC_INTERNAL_SERVER_ERROR;
    }
}
