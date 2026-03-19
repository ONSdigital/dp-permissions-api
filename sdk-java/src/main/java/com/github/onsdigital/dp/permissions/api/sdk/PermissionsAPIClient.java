package com.github.onsdigital.dp.permissions.api.sdk;

import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.BadRequestException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PermissionsAPIException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PolicyNotFoundException;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.HttpStatus;
import org.apache.hc.client5.http.classic.methods.HttpDelete;
import org.apache.hc.client5.http.classic.methods.HttpUriRequestBase;
import org.apache.hc.client5.http.classic.methods.HttpUriRequest;
import org.apache.hc.core5.http.io.HttpClientResponseHandler;
import org.apache.hc.core5.http.io.entity.EntityUtils;

public class PermissionsAPIClient implements PermissionsClient {

    /**
     * uri for the permissionsAPI.
     */
    private final URI permissionsAPIUri;

    /**
     * Auth token to be used on all requests.
     */
    private final String authToken;

    /**
     * HTTP client to be used on all requests.
     */
    private final CloseableHttpClient client;

    /**
     * Header name to apply authToken to.
     */
    private static final String SERVICE_TOKEN_HEADER_NAME = "Authorization";

    /**
     * Create a new instance of PermissionsAPIClient.
     *
     * @param permissionsAPIURL The URL of the permissions API
     * @param serviceAuthToken  The authentication token for the permissions API
     * @param httpClient        The HTTP client to use internally
     */
    public PermissionsAPIClient(final String permissionsAPIURL,
            final String serviceAuthToken, final CloseableHttpClient httpClient)
            throws URISyntaxException {

        this.permissionsAPIUri = new URI(permissionsAPIURL);
        this.client = httpClient;
        this.authToken = serviceAuthToken;
    }

    /**
     * Create a new instance of PermissionsAPIClient with a default Http client.
     *
     * @param permissionsAPIURL The URL of the permissions API
     * @param serviceAuthToken  The authentication token for the permissions API
     * @throws URISyntaxException
     */
    public PermissionsAPIClient(final String permissionsAPIURL,
            final String serviceAuthToken)
            throws URISyntaxException {
        this(permissionsAPIURL, serviceAuthToken, createDefaultHttpClient());
    }

    private static CloseableHttpClient createDefaultHttpClient() {
        return HttpClients.createDefault();
    }

    /**
     * Deletes a policy by sending a DELETE request to /policies/{id}.
     * The {@code policyID} is the policy ID
     *
     * @param policyID the policy ID
     * @throws IOException if the request fails
     */
    @Override
    public void deletePolicy(final String policyID)
            throws IOException, BadRequestException, PolicyNotFoundException,
            PermissionsAPIException {

        if (policyID == null || policyID.isEmpty()) {
            throw new IllegalArgumentException(
                    "'policyID' must not be null or empty");
        }

        URI requestUri = permissionsAPIUri.resolve("/v1/policies/" + policyID);
        HttpDelete req = new HttpDelete(requestUri);

        req.addHeader(SERVICE_TOKEN_HEADER_NAME, "Bearer " + authToken);

        ResponseResult response = executeRequest(req);
        validateResponseCode(req, response.getStatusCode(), response.getBody());
    }

    private String formatErrResponse(final HttpUriRequestBase httpRequest,
            final int responseCode,
            final int expectedStatusCode,
            final String responseBody) {
        String requestURI = httpRequest.getRequestUri();
        String message = String.format(
                "the permissions api returned a %s response for %s "
                        + "(expected %s)",
                responseCode,
                requestURI,
                expectedStatusCode);
        if (responseBody != null && !responseBody.isEmpty()) {
            message = message + ": " + responseBody;
        }
        return message;
    }

    private void validateResponseCode(final HttpUriRequestBase httpRequest,
            final int statusCode,
            final String responseBody)
            throws BadRequestException, PolicyNotFoundException,
            PermissionsAPIException {
        switch (statusCode) {
            case HttpStatus.SC_OK:
                return;
            case HttpStatus.SC_BAD_REQUEST:
                throw new BadRequestException(formatErrResponse(httpRequest,
                        statusCode, HttpStatus.SC_BAD_REQUEST, responseBody),
                        statusCode);
            case HttpStatus.SC_NOT_FOUND:
                throw new PolicyNotFoundException(formatErrResponse(httpRequest,
                        statusCode, HttpStatus.SC_NOT_FOUND, responseBody),
                        statusCode);
            default:
                throw new PermissionsAPIException(formatErrResponse(httpRequest,
                        statusCode, HttpStatus.SC_INTERNAL_SERVER_ERROR,
                        responseBody), statusCode);
        }
    }

    private ResponseResult executeRequest(final HttpUriRequest req)
            throws IOException {
        HttpClientResponseHandler<ResponseResult> handler = response -> {
            int statusCode = response.getCode();
            String body = null;
            if (response.getEntity() != null) {
                body = EntityUtils.toString(response.getEntity());
            }
            return new ResponseResult(statusCode, body);
        };
        return client.execute(req, handler);
    }

    /**
     * Class to hold the result of an API response, including the status
     * code and body.
     * This is used internally to pass the response details from the
     * executeRequest method to the validateResponse method.
     */
    private static final class ResponseResult {
        /** The HTTP status code of the response. */
        private final int statusCode;
        /** The body of the response. */
        private final String body;

        private ResponseResult(final int responseStatusCode,
                final String responseBody) {
            this.statusCode = responseStatusCode;
            this.body = responseBody;
        }

        private int getStatusCode() {
            return statusCode;
        }

        private String getBody() {
            return body;
        }
    }

    /**
     * Close the http client used by the APIClient.
     *
     * @throws IOException
     */
    @Override
    public void close() throws IOException {
        client.close();
    }
}
