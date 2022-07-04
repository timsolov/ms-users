# ms-users (WIP) - boilerplate for microservice with clean architecture

## Features
- Totally support Clean Architecture;
- Auto-generate gRPC-server, gRPC-client, HTTP/Web server, swagger documentation from .proto files;
- Graceful shutdown;
- Accept interface, return struct pattern;
- CQRS pattern for usecases;
- PASETO token;

## Dependancies

- PostgreSQL - OLTP database for storing data;
- PgQ - PostgreSQL native queue plugin for handling `outbox` pattern;
- ms-emails - service for sending emails;

## Prometeus metrics

http://0.0.0.0/metric/

## TODO
- [x] Create profile with email-password identity
- [x] Confirmation of email-password identity
- [x] Repeat email confirmation of email-password identity
- [x] User's profile
- [x] Authentication by email-password identity
    - [x] Switch JWT token to PASETO token
- [x] Identity reset password process
    - [x] Init reset password process end-point
    - [x] Confirm reset password process and set new password end-point
- [x] Timeout for http handlers
- [ ] JSONSchema configurable profile info
- [ ] Healthcheck for all dependencies
    - [ ] PostgreSQL
    - [ ] ms-emails
- [ ] Opentelemetry
- [x] Prometheus
- [ ] Authentication by Google OAuth 2.0
- [ ] Authentication by phone-password identity