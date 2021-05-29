# teleport_interview
A repository for an implementation of [Teleport's fullstack interview challenge](https://github.com/gravitational/careers/blob/main/challenges/fullstack/dashboard.pdf).

## Requirements
Requires Docker/Docker-Compose (a recent enough version to support compose spec 3.9).

## Commands
All commands can be found in the `Makefile` in the root of the repository.
- `make build` - Build the Compose service for the app.
- `make up` - Serve the app over HTTPS at localhost on port 8080. Ctrl+C to stop.
- `make test_backend` - Run backend tests.
- `make test_frontend` - Run frontend tests.
- `make test` - Run both backend and frontend tests.
