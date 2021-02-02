import-role
================

This utility adds the roles stored in the (roles.json)[roles.json] file to mongodb. 

### How to run the utility

* Run `go run import-roles.go -mongo-url=<url>`

The url should look like the following `localhost:27017`. If a username and password are needed follow this structure `<username>:<password>@<host>:<port>`