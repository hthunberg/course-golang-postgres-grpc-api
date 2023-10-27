# Build

## sqlc to generate database code for PostgreSQL

We use [sqlc](https://docs.sqlc.dev/en/stable/index.html) for database queries and mapping. This library has support for PostgreSQL, MySQL and SQLite, when i did this course.

## Docker compose landscape

### Start landscape

Makefile contains targets wrapping some commands for managing the docker compose landscape.

~~~
$ make up
~~~

### PostgreSQL Adminer UI

Login with [Postgres Adminer](https://www.adminer.org) using

[http://localhost:9080](http://localhost:9080)

* System: PostgreSQL
* Server: db
* Username: postgres
* Password: docker
* Database: bankdb

### psql cli

Using make target

~~~
$ make db-shell
docker compose -f docker-compose.yml exec db psql -U postgres -d postgres
psql (15.4 (Debian 15.4-1.pgdg110+1))
Type "help" for help.

postgres=#
~~~

### Migrate DB

We use [golang-migrate](https://github.com/golang-migrate/migrate) to manage db migrations.

Migrate to next version

~~~
$ make migrate-db-up
20230907110931/u bank-db-migration (5.762833ms)
~~~

Migrate back to previous version

~~~
$ make migrate-db-down
20230907110931/d bank-db-migration (4.526208ms)
~~~

Create next migration path

~~~
$ make migrate-db-create-next-migration-path
/migrations/20230907110931_bank-db-migration.up.sql
/migrations/20230907110931_bank-db-migration.down.sql
~~~

## Test

Example uses [httpie](https://httpie.io)

~~~json
$ http POST localhost:8080/users username=johndoe password=qwerty full_name=John email=john.doe@testbank.qwerty
HTTP/1.1 200 OK
Content-Length: 164
Content-Type: application/json; charset=utf-8
Date: Fri, 27 Oct 2023 18:41:58 GMT

{
    "created_at": "2023-10-27T18:41:58.248839Z",
    "email": "john.doe@testbank.qwerty",
    "full_name": "John",
    "password_changed_at": "0001-01-01T00:00:00Z",
    "username": "johndoe"
}
~~~