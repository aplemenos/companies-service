### Companies microservice in Golang [Clean Architecture] - Rest API example

#### 👨 Full list what has been used:
* [gin](github.com/gin-gonic/gin) - Web framework
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [jaeger](https://github.com/uber/jaeger) - Jaeger
* [tracer](https://github.com/opentracing/opentracing-go) - OpenTracing
* [kafka](https://github.com/segmentio/kafka-go) - Kafka
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [jwt-go](https://github.com/dgrijalva/jwt-go) - JSON Web Tokens (JWT)
* [uuid](https://github.com/google/uuid) - UUID
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [bluemonday](https://github.com/microcosm-cc/bluemonday) - HTML sanitizer
* [swag](https://github.com/swaggo/swag) - Swagger
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker

#### Recomendation for docker development:
    make prod // run all containers

#### Docker-compose files:
    docker-compose.yml - production build

#### Run linter
    make run-linter

#### Run tests
   make test

#### Migrate the DB schema
    make migrate_up - you can easily migrate the DB schema in the postgres DB of docker image

#### Clean all docker containers
    make down-local

### Swagger UI:

http://localhost:8080/swagger/index.html

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000