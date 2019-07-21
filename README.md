# golang-backend-postgres

This project illustrates a backend microservice built with golang which serves data from a postgresql database.

One thing to note is I am using some of the patterns from https://github.com/sourcegraph/sourcegraph to power this service, so big shout out to a great product and project.

# Goals

The goal of the project is to cover a few key pillars:

1. Perform the core function of providing a REST service for data stored in a PostgreSQL database.
2. Be observable, through the delivery of metrics, structured logs and open tracing data to a monitoring system.
3. Be operations friendly, configuration via env variables, and simple deployment model using docker or binary.
4. Illustrate good test coverage across the entire service, both unit and integration.
5. Be secure, provide standard robust authentication, authorisation and auditing facilities.

# Dependencies

- [golangci-lint](https://github.com/golangci/golangci-lint) is used for linting.
- [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) is used for migrations, the CLI can be installed using `brew install golang-migrate`.

# License

This code is Copyright Mark Wolfe and licensed under Apache License 2.0