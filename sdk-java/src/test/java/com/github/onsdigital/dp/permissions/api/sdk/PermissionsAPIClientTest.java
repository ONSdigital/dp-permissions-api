package com.github.onsdigital.dp.permissions.api.sdk;

import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.classic.methods.HttpDelete;
import com.github.onsdigital.dp.permissions.api.sdk.exception.BadRequestException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PolicyNotFoundException;
import org.apache.hc.core5.http.io.HttpClientResponseHandler;
import org.junit.jupiter.api.Test;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

class PermissionsAPIClientTest {

    private static final String SERVICE_AUTH_TOKEN = "test-token";

    /**
     * Permissions API URL for testing.
     */
    private static final String PERMISSIONS_API_URL = "http://permissions-api:1234";

    @Test
    void testInvalidUrlThrowsURISyntaxException() {
        String invalidUrl = "ht!tp://bad_url";
        CloseableHttpClient mockClient = mock(CloseableHttpClient.class);
		assertThrows(java.net.URISyntaxException.class, () -> {
			try (PermissionsAPIClient client = new PermissionsAPIClient(invalidUrl, SERVICE_AUTH_TOKEN, mockClient)) {
				// The exception is expected to be thrown during construction
			}
		});
    }

	@Test
	void testDeletePolicySuccess() throws Exception {
		CloseableHttpClient mockClient = mock(CloseableHttpClient.class);
		PermissionsAPIClient client = new PermissionsAPIClient(PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN, mockClient);
		CloseableHttpResponse mockResponse = MockHttp.response(200);
		stubExecuteWithHandler(mockClient, mockResponse);
		assertDoesNotThrow(() -> client.deletePolicy("policy123"));
		client.close();
	}

	@Test
	void testDeletePolicyNotFound() throws Exception {
		CloseableHttpClient mockClient = mock(CloseableHttpClient.class);
		PermissionsAPIClient client = new PermissionsAPIClient(PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN, mockClient);
		CloseableHttpResponse mockResponse = MockHttp.response(404);
		stubExecuteWithHandler(mockClient, mockResponse);
		assertThrows(PolicyNotFoundException.class, () -> client.deletePolicy("policy123"));
		client.close();
	}

	@Test
	void testDeletePolicyBadRequest() throws Exception {
		CloseableHttpClient mockClient = mock(CloseableHttpClient.class);
		PermissionsAPIClient client = new PermissionsAPIClient(PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN, mockClient);
		CloseableHttpResponse mockResponse = MockHttp.response(400);
		stubExecuteWithHandler(mockClient, mockResponse);
		assertThrows(BadRequestException.class, () -> client.deletePolicy("policy123"));
		client.close();
	}

	@Test
	void testDeletePolicyThrowsOnNullOrEmpty() throws Exception {
		CloseableHttpClient mockClient = mock(CloseableHttpClient.class);
		PermissionsAPIClient client = new PermissionsAPIClient(PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN, mockClient);
		assertThrows(IllegalArgumentException.class, () -> client.deletePolicy(null));
		assertThrows(IllegalArgumentException.class, () -> client.deletePolicy(""));
		client.close();
	}

	@SuppressWarnings({"rawtypes", "unchecked"})
	private void stubExecuteWithHandler(final CloseableHttpClient mockClient,
			final CloseableHttpResponse mockResponse) throws Exception {
		when(mockClient.execute(any(HttpDelete.class), any(HttpClientResponseHandler.class)))
				.thenAnswer(invocation -> {
					HttpClientResponseHandler handler = invocation.getArgument(1);
					return handler.handleResponse(mockResponse);
				});
	}


}
