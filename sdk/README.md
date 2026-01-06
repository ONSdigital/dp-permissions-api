# dp-permissions-api SDK

## Overview

This API contains a client with convenience functions for interacting with the permissions api from other applications.
It also contains reusable structs representing models used as payloads in API endpoints.

## Example use of the API SDK

```go
package main

import (
    "context"
    "github.com/ONSdigital/dp-permissions-api/sdk"
)

func main() {
    apiClient := sdk.NewClient("http://localhost:25400")

    permissionsBundle, err := c.cache.GetPermissionsBundle(context.Background())
}
```

## Alternative Client instantiation

In the unlikely event that there is a need to use non-default initialisation, it is possible to obtain a new client with an underlying http client.

### With a custom http client

```go
import dphttp "github.com/ONSdigital/dp-net/v3/http"

apiClient := sdk.NewClientWithClienter("http://localhost:25400", dphttp.NewClient())
```

## Additional Information

### Headers

The [`Headers`](headers.go) struct allows the user to provide an Authorization header if required.
This must be set without the `"Bearer "` prefix as the SDK will automatically add this.
