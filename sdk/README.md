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

    permissionsBundle, err := c.cache.GetPermissionsBundle(context.Backgroud)
}
```

## Alternative Client instantiation

In the unlikely event that there is a need to use non-default initialisation, it is possible to obtain a new client with
customised options and/or underlying http client.

### With Options

NB. There are currently no defined options available, this is included for future expansion.

```go
apiClient := sdk.NewClientWithOptions("http://localhost:25400", sdk.Options{})
```

### With Options and a custom http client

```go
import dphttp "github.com/ONSdigital/dp-net/v2/http"

apiClient := sdk.NewClientWithClienter("http://localhost:25400", dphttp.NewClient(), sdk.Options{})
```
