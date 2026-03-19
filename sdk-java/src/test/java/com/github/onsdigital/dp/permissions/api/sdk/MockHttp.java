package com.github.onsdigital.dp.permissions.api.sdk;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.core5.http.io.entity.StringEntity;

import java.io.UnsupportedEncodingException;

import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

class MockHttp {

    protected MockHttp() {
        // prevents calls from subclass
        throw new UnsupportedOperationException();
      }

    /**
     * JSON mapper.
     */
    private static final ObjectMapper JSON = new ObjectMapper();

    public static CloseableHttpResponse response(final int httpStatus) {

        CloseableHttpResponse mockHttpResponse = mock(
                CloseableHttpResponse.class);

        when(mockHttpResponse.getCode()).thenReturn(httpStatus);

        return mockHttpResponse;
    }

    public static void responseBody(
            final CloseableHttpResponse mockHttpResponse,
            final Object responseBody)
            throws JsonProcessingException, UnsupportedEncodingException {
        String responseJSON = JSON.writeValueAsString(responseBody);
        when(mockHttpResponse.getEntity()).thenReturn(
                new StringEntity(responseJSON));
    }
}
