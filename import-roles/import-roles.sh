#!/usr/bin/env bash

if [[ -z "$1" ]]
  then
    echo "Please supply the mongo connection string as the first parameter, e.g mongodb://localhost:27017"
    exit 1
fi

mongo $1 <<EOF

 var file = cat('./roles.json');
 use permissions
 var roles = JSON.parse(file);
 db.permissions.insert(roles)

EOF