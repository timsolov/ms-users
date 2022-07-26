# ms-users (WIP) - boilerplate for microservice with clean architecture

## General features
- Totally support Clean Architecture;
- Auto-generate gRPC-server, gRPC-client, HTTP/Web server, swagger documentation from .proto files;
- Graceful shutdown;
- Accept interface, return struct pattern;
- CQRS pattern for usecases;
- PASETO token;

## Service specific features
- [x] Create profile with email-password identity
- [x] Confirmation of email-password identity
- [x] Repeat email confirmation for email-password identity
- [x] User's profile
- [x] Update user's profile
- [x] Authentication by email-password identity
    - [x] Switch JWT token to PASETO token
- [x] Reset password process for email-password identity
    - [x] Init reset password process end-point
    - [x] Confirm reset password process and set new password end-point
- [x] Timeout for http handlers
- [x] JSONSchema configurable profile info
- [ ] Healthcheck for all dependencies
    - [x] PostgreSQL
    - [x] gRPC Server
    - [ ] ms-emails
- [x] Status of service
- [ ] Opentelemetry
- [x] Prometheus
- [ ] Authentication by Google OAuth 2.0
- [ ] Authentication by phone-password identity

## Dependancies

- PostgreSQL - OLTP database for storing data;
- PgQ - PostgreSQL native queue plugin for handling `outbox` pattern;
- ms-emails - service for sending emails;

## Prometeus metrics

Prometeus metrics available on:
http://\$HTTP_HOST:\$HTTP_PORT/metric/

Default:
http://0.0.0.0:11000/metric/

## Healthcheck

http://\$HTTP_HOST:\$HTTP_PORT/health/

Default:
http://0.0.0.0:11000/health/

After making the request you will receive a response:
```jsonc
{
  "app": {
    "buildtime": "2022-07-26T19:19:51Z",
    "version": "108f44c"
  },
  "grpc_server": {
    "duration": "649.655µs",
    "status": "UP"
  },
  "postgres": {
    "duration": "1.039368ms",
    "stats": {
      "idle": 0,
      "in_use": 1,
      "max_idle_closed": 0,
      "max_idle_time_closed": 0,
      "max_life_time_closed": 0,
      "max_open_connections": 5,
      "open_connections": 1,
      "wait_count": 0,
      "wait_duration": 0
    },
    "status": "UP"
  },
  "status": "UP"
}
```

## Status

http://\$HTTP_HOST:\$HTTP_PORT/status/

Default:
http://0.0.0.0:11000/status/

Response values:
1. `NOT_READY` - the service not ready yet it's in preparing dependancies status;
2. `UP` - the service in battle mode you can use it;
3. `DOWN` - the service is shuting down.