This directory contains database migrations for the backend [PostgreSQL](https://www.postgresql.org/) database.

This project uses [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) to manage database migrations.

#  Add a new migration

**IMPORTANT:** All migrations must be backward-compatible, meaning that the existing version of the backend command must be able to run against the new (post-migration) version of the schema.

Run the following:

```
./dev/add_migration.sh MIGRATION_NAME
```

After adding SQL statements to those files, embed them into the Go code:

```
make generate
```

To only run the DB generate scripts (subset of the command above):

```
go generate ./migrations/
```