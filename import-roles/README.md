import-role
================

This script is to add roles stored in the roles.json file to mongo. 

### How to run the script

* Run `go build`
* Run `./import-roles.go -mongo-url=<url>`

The url should look like the following `localhost:27017`. If a username and password are needed follow this structure `<username>:<password>@<host>:<port>`