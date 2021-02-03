# go-scrt-events

Golang rewrite of get-scrt-events. Currently a work in progress. 

## Work in Progress

## How go?

### Configuration 

Connection to a Secret Node or Figment Datahub can be configured via the 
`node`-> `host`and/or `path` variable.

Connection string for postgres can be configured easiest via the `conn` variable. 

```yaml
node:
  host: "secret-2--rpc--full.datahub.figment.io"
  path: "APIKEY or Path if necessary"

database:
  conn: "postgres://postgres:postgrespassword@postgres:5432/postgres"
```

### Quickstart - Docker compose


### Data Model
