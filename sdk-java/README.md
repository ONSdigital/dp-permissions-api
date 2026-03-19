# dp-permissions-api Java SDK

## Overview

This module provides a Java client for interacting with the dp-permissions-api.
At the moment, the SDK supports deleting policies via the permissions API.
Refer to the [swagger specification](../swagger.yaml) for endpoint details.

## Add to your pom.xml

```xml
<dependency>
  <groupId>com.github.onsdigital.dp-permissions-api</groupId>
  <artifactId>dp-permissions-api-sdk-java</artifactId>
  <version>${permissionsSDK.version}</version>
</dependency>
```

## Initialise a client

```java
package com.github.onsdigital.dis.my.application;

import com.github.onsdigital.dp.permissions.api.sdk.PermissionsAPIClient;
import com.github.onsdigital.dp.permissions.api.sdk.PermissionsClient;

public class MyApplicationClass {

    private static final String PERMISSIONS_API_URL = "http://localhost:29900";
    private static final String SERVICE_AUTH_TOKEN = "xyz1234";

    public static void main(String[] args) throws Exception {
        try (PermissionsClient client = new PermissionsAPIClient(
                PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN)) {
            // use client
        }
    }
}
```

## Delete a policy

```java
import com.github.onsdigital.dp.permissions.api.sdk.PermissionsAPIClient;
import com.github.onsdigital.dp.permissions.api.sdk.PermissionsClient;
import com.github.onsdigital.dp.permissions.api.sdk.exception.BadRequestException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PolicyNotFoundException;
import com.github.onsdigital.dp.permissions.api.sdk.exception.PermissionsAPIException;

try (PermissionsClient client = new PermissionsAPIClient(
        PERMISSIONS_API_URL, SERVICE_AUTH_TOKEN)) {
    client.deletePolicy("policy-id");
} catch (BadRequestException ex) {
    // invalid policy id or request
} catch (PolicyNotFoundException ex) {
    // policy does not exist
} catch (PermissionsAPIException ex) {
    // other API error
}
```

## Notes

- `deletePolicy` throws `BadRequestException`, `PolicyNotFoundException`, or
  `PermissionsAPIException` for non-2xx responses.
- Additional endpoints will be added in the future.
