package com.github.onsdigital.dp.permissions.api.sdk;

import java.io.Closeable;
import java.io.IOException;

import com.github.onsdigital.dp.permissions.api.sdk.exception.BadRequestException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PermissionsAPIException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PolicyNotFoundException;

public interface PermissionsClient extends Closeable {
    /**
     * Deletes a policy by sending a DELETE request to the /policies/{id}
     * endpoint.
     *
     * A {@code 204 No Content} status indicates successful deletion.
     * A {@code 404 Not Found} status indicates the policy does not exist.
     *
     * @param policyID the policy ID
     * @throws IOException             if an I/O error occurs during the request
     * @throws PermissionsAPIException if the permissions API returns an error
     *                                 response
     */
    void deletePolicy(String policyID) throws IOException, BadRequestException,
            PolicyNotFoundException, PermissionsAPIException;

    //TODO: add methods for the other API endpoints (e.g. createPolicy,
    //putPolicy etc.)
}
