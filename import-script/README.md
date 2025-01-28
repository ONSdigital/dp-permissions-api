# Import Script

This utility adds predefined data to the MongoDB permissions API database. The utility uses the same config type as the permissions API service, so any custom configuration can be added via environment variables.

- The roles collection is populated with data in the [roles.json](roles.json) file.
- The policies collection is populated with data in the [policies.json](policies.json)

## How to run the utility against a local MongoDB

In a terminal, ensure you are in the import-script directory:

```sh
cd import-script
```

Run the import script with the default configuration:

```sh
go run import.go
```

## How to run the utility against an environment (DocumentDB)

In a terminal, ensure you are in the import-script directory:

```sh
cd import-script
```

Open an SSH tunnel to the environment (replace `{cluster address}`):

```sh
dp ssh develop publishing 1 -p 27017:{cluster address}:27017
```

Run the import script, setting the required configuration values:

```sh
MONGODB_IS_SSL=true MONGODB_USERNAME=... MONGODB_PASSWORD=... go run import.go
```
