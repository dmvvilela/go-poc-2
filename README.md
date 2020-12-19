# Go lang POC 2 - Daniel Vilela

Golang CRUD with gorilla mux, mysql or postgresql and all the tests.

It is an API for contacts with name and email (unique).

Create a .env file like below. Just uncomment if you want to use mysql or postgres and change it for your user and password.

It uses different databases for production and for testing, so create them accordingly.

```
# Postgres Live
DB_HOST=127.0.0.1
DB_DRIVER=postgres
DB_USER=user
DB_PASSWORD=password
DB_NAME=go-poc-2
DB_PORT=5432 #Default postgres port

# Postgres Test
TestDbHost=127.0.0.1
TestDbDriver=postgres
TestDbUser=user
TestDbPassword=password
TestDbName=go-poc-2_api_test
TestDbPort=5432

# Mysql Live
# DB_HOST=127.0.0.1
# DB_DRIVER=mysql
# DB_USER=user
# DB_PASSWORD=password
# DB_NAME=go-poc-2_api
# DB_PORT=3306 #Default mysql port

# Mysql Test
# TestDbHost=127.0.0.1
# TestDbDriver=mysql
# TestDbUser=user
# TestDbPassword=password
# TestDbName=go-poc-2_api_test
# TestDbPort=3306
```
