dp-permissions-api
================
API for managing access control permissions for Digital Publishing API resources

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable              | Default   | Description
| ----------------------------      | --------- | -----------
| BIND_ADDR                         | :25400    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT         | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL              | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT      | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| MOGODB_BIND_ADDR                  | localhost:27017 | The MongoDB bind address
| MONGODB_PERMISSIONS_DATABASE      | permissions     | The MongoDB permissions database
| MONGODB_PERMISSIONS_COLLECTION    | permissions     | The MongoDB permissions collection

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

