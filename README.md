dp-permissions-api
================
API for managing access control permissions for Digital Publishing API resources

### Getting started

* Load inital roles into your local mongodb using the `import-script` utility in the import-script folder. Follow the
  steps in the [README](./import-script/README.md).
* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable           | Default                                             | Description                                                                                                         |
|--------------------------------|-----------------------------------------------------|---------------------------------------------------------------------------------------------------------------------|
| BIND_ADDR                      | :25400                                              | The host and port to bind to                                                                                        |
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                                                  | The graceful shutdown timeout in seconds (`time.Duration` format)                                                   |
| HEALTHCHECK_INTERVAL           | 30s                                                 | Time between self-healthchecks (`time.Duration` format)                                                             |
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                                                 | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)  |
| MONGODB_BIND_ADDR              | localhost:27017                                     | The MongoDB bind address                                                                                            |
| MONGODB_USERNAME               |                                                     | The MongoDB Username                                                                                                |
| MONGODB_PASSWORD               |                                                     | The MongoDB Password                                                                                                |
| MONGODB_DATABASE               | permissions                                         | The MongoDB database                                                                                                |
| MONGODB_COLLECTIONS            | RolesCollection:roles, PoliciesCollection:policies  | The MongoDB collections                                                                                             |
| MONGODB_REPLICA_SET            |                                                     | The name of the MongoDB replica set                                                                                 |
| MONGODB_ENABLE_READ_CONCERN    | false                                               | Switch to use (or not) majority read concern                                                                        |
| MONGODB_ENABLE_WRITE_CONCERN   | true                                                | Switch to use (or not) majority write concern                                                                       |
| MONGODB_CONNECT_TIMEOUT        | 5s                                                  | The timeout when connecting to MongoDB (`time.Duration` format)                                                     |
| MONGODB_QUERY_TIMEOUT          | 15s                                                 | The timeout for querying MongoDB (`time.Duration` format)                                                           |
| MONGODB_IS_SSL                 | false                                               | Switch to use (or not) TLS when connecting to mongodb                                                               |
| DEFAULT_LIMIT                  | 20                                                  | Default limit for pagination                                                                                        |
| DEFAULT_OFFSET                 | 0                                                   | Default offset for pagination                                                                                       |
| DEFAULT_MAXIMUM_LIMIT          | 1000                                                | Default maximum limit for pagination                                                                                |

### SDK Package

[An SDK for the API is available as a subpackage in `/sdk`](sdk/README.md)

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

