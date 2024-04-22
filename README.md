# SQLC Demo

This Project contains a demo of [sqlc](https://docs.sqlc.dev/en/latest/) with [golang-migrate](https://github.com/golang-migrate/migrate). Also demo is created using POSTGRESQL Database with [pgx](https://github.com/jackc/pgx) with version 5 that is pgx/v5 instead of native library [lib/pq](https://github.com/lib/pq) as it has long undergone maintenance mode and pgx library is recommended by them for further use.

### ENV Variables
_________________
Below are the env variables needed to run the example.
```dotenv
DRIVER_NAME=pgx
DB_HOST=
DB_NAME=
DB_PASSWORD=
DB_PORT=
DB_USER=
```

### Gist of Example
_________________
- Logs are added to display and understand the working of the example.
- Down functionality is added to drop the migration created.
- DB created for migration and for performing sql queries are different as migrator library needs DB of database/sql library and query needs pgx/v5 DB for performing any operations.
- Both DB And Migrator are closed for cleanup process and as of good practice to do so.
- An Example to create and retrieve example is implemented to display the simple working.
