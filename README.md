# goquiz

## Quick start

    git clone git@github.com:alexbyk/goquiz.git
    cd goquiz
    docker-compose down && docker-compose up

This will generate mock csv file and run a workflow

## Testing

    docker run --rm -p 5432:5432 -e POSTGRES_USER=test  postgres:alpine postgres -c log_statement=all
    TEST_DSN="postgres://test@localhost/test?sslmode=disable" go test -count=1 -v ./...  -timeout 2s

## Improvements

I've made some improvements to the original task.
  - The logic doesn't depend on `postgres/CSV/http`. We can use any backends that satisfy required interfaces
  - Services may be located on different machines, as you can see in the `docker-compose.yml` file
  - The project tolerates different faults on `source`, `database` and `3d party api` as well. For example, if it gets broken `csv` file, it tries to handle as many correct records as possible.

## Description

![Project diagram](https://github.com/alexbyk/goquiz/blob/master/diagram.png?raw=true)

Service `parser` loads CSV records in chunks (can be configured according to available resources), stores it into a database and notifies service `integrator`. Each customer record should have a unique `ID`, otherwise, it will be silently ignored.

Service `integrator` listens for notifications and act as a daemon, sending records one by one to `httpapi` ASAP. If `httpapi` responds with `200`, the record marked as successfully sent. If code isn't `200`, `integrator` considers it as `API is down` and will continue after some interval until success.

Notifications are used to send records ASAP in realtime

## Code structure

This project can use different backends for source, storage and api integration.

As requested, there are backend implementations: `postgres` for storage and notification, `CSV` as a source and `http` as 3d party api. Mock backends are used for testing.

  - **impl** - contains `CSV`, `Postgres` and `http` backends, wich satisfy common interfaces
  - **common** - contains business logic and interfaces for it.
  - **cmd** - contains binaries

There can be many `parsers` with different sources. `parsers` communicate with a single `integrator`, which is responsible for data consistency. `integrator` sends data to dummy `httpapi`, which emulates `3d party api`

