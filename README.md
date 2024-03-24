# exitus

exitus is a very simple issue tracking API, it was built to illustrate a backend [microservice](https://martinfowler.com/articles/microservices.html) written in [Go](https://golang.org/) which serves data from a [PostgreSQL](https://www.postgresql.org/) database.

The name is derived from the origin of the word "issue".

> Middle English (in the sense ‘outflowing’): from Old French, based on Latin exitus, past participle of exire ‘go out’.

One thing to note is I am using some of the patterns from https://github.com/sourcegraph/sourcegraph to power this service, so big shout out to a great product and project.

# Goals

The goal of the project is to cover a few key pillars:

1. Perform the core function of providing a REST service for data stored in a PostgreSQL database.
2. Be observable, through the delivery of metrics, structured logs and open tracing data to a monitoring system.
3. Be operations friendly, configuration via env variables, and simple deployment model using docker or binary.
4. Illustrate good test coverage across the entire service, both unit and integration.
5. Be secure, provide standard robust authentication, authorisation and auditing facilities.

# Overview

## Authentication

Authentication for this service is provided by an external OpenID provider such as [AWS Cognito](https://aws.amazon.com/cognito/) or [Keycloak](https://www.keycloak.org). Clients authenticate with one of these services and then provide their JWT token, which is [validated](pkg/jwt/jwt.go) by the exitus service using [go-oidc](github.com/wolfeidau/go-oidc) developed by CoreOS.


### Secrets

With the help of CDK the RDS credentials are available from [AWS secrets manager](https://aws.amazon.com/secrets-manager/) using [Specifying Sensitive Data](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/specifying-sensitive-data.html) feature of ECS. This is neatly wrapped by CDK, and the JSON value containing host, port, database, username and password is [unmarshalled by my service](pkg/conf/conf.go). This eleminates display of these values in the ECS console environment variables section.

# Development

To run the tests, and run locally you will need two environment variables, being:

```
# Used to connect to AWS for deployments and monitoring
export AWS_PROFILE=whatever
export AWS_REGION=ap-southeast-2

# Nice random passwords go here to seed the passwords used by docker compose because security
export POSTGRES_PASSWORD=xxx
export POSTGRES_ROOT_PASSWORD=xxx

# Used for development locally
export OAUTH_CLIENT_ID=xxx
export OAUTH_CLIENT_SECRET=xxx
export OPENID_PROVIDER_URL=https://cognito-idp.ap-southeast-2.amazonaws.com/ap-southeast-2_XXXXXXXXXS

# Used by CDK for deployment
export STAGE=dev
export BRANCH=master
export DOMAIN_NAME=whatever.cloud.
export HOSTED_ZONE_ID=XXXXXXXXXXXXX
export ACM_CERTIFICATE_ARN=arn:aws:acm:ap-southeast-2:123456789012:certificate/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

These vars are used in docker compose when configuring the PostgreSQL database.

**Note:** I recommend using [direnv](https://direnv.net/) to load environment variables when you navigate into the project, there is an `.envrc.example` which can be renamed to `.envrc` to set these values.

Once these variables are configured run:

```
make docker-compose-test
```

# Dependencies

* [golangci-lint](https://github.com/golangci/golangci-lint) is used for linting.
* [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) is used for migrations, the CLI can be installed using `brew install golang-migrate`.

# References

* [F1: A Distributed SQL Database That Scales](http://static.googleusercontent.com/media/research.google.com/en//pubs/archive/41344.pdf)
* [Using PostgreSQL Arrays with Golang](https://www.opsdash.com/blog/postgres-arrays-golang.html)
* [PostgreSQL Foreign Key](http://www.postgresqltutorial.com/postgresql-foreign-key/)

# License

This code is Copyright Mark Wolfe and licensed under Apache License 2.0
