import script
================

This utility adds predefined data to the MongoDB permissions API database.
- The roles collection is populated with data in the [roles.json](roles.json) file. 
- The policies collection is populated with data in the [policies.json](policies.json)

### How to run the utility

* Run `go run import.go -mongo-url=<url>`

The url should look like the following `localhost:27017`. If a username and password are needed follow this structure `<username>:<password>@<host>:<port>`