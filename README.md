# aggrego
An RSS feed aggregator.

## Setup

Aggrego uses [PostgreSQL]() as its database engine and uses [goose](https://github.com/pressly/goose) and [sqlc](https://github.com/sqlc-dev/sqlc) for migrations and SQL query generation.

A `.env` file at the top level is used for both database migrations and server initialization. It uses the following variables:

|Variable| Description |
|--------|------------ |
| `PORT` | The port the server will listen and serve to. |
| `CONN` | The URL for connecting to the postgres services. Expected format is: `postgres://profile:password@address:port/databasename?sslmode=disable`.
| `TESTFEEDS` | Used by `main_test` to test the database and api endpoints. Expected to be a comma-seprated list of URLs for RSS feeds. |

Once the `.env` file has been created, the following commands will setup and run the server.

```bash
sqlc generate
bash migrate.sh up
go install
go run main.go
```