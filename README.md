# ms-users (WIP) - boilerplate for microservice with clean architecture

## Features
- Totally support Clean Architecture;
- Auto-generate gRPC-server, gRPC-client, HTTP/Web server, swagger documentation from .proto files;
- Graceful shutdown;
- Accept interface, return struct pattern;
- CQRS pattern for usecases;
- PASETO token;

## Prometeus metrics

http://0.0.0.0/metric/

## TODO
[ ] Healthcheck for all dependencies
    [ ] PostgreSQL
    [ ] NATS
[ ] Opentelemetry
[x] Prometheus