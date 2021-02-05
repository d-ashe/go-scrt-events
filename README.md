# go-scrt-events

Golang rewrite of get-scrt-events. Currently a work in progress. 

## Work in Progress


### Description

The application first queries Postgresql to see the available block-heights for the input chain-id. All communication with Secret node is via websockets.

The /status? endpoint is then queried for the latest block height. The latest block height is used to determine which blocks remain to bring the database up to date. All remaining block-heights  are then requested via the /block_results?height=_ endpoint. 

As a base, the application parses all Blocks, Txs, and Events into separate tables. Within the configuration, it is possible to specify further enrichments.

### How go?

### Configuration 

Connection to a Secret Node or Figment Datahub can be configured via the 
`node`-> `host`and/or `path` variable.

Connection string for postgres can be configured easiest via the `conn` variable. 

Setting the `enrichment`-> `run`variable list allows specification of enrichments to be processed at runtime. 

```yaml
node:
  host: "secret-2--rpc--full.datahub.figment.io"
  path: "APIKEY or Path if necessary"

database:
  conn: "postgres://postgres:postgrespassword@postgres:5432/postgres"

enrichment:
  run: ["all"]
```

### Quickstart - Docker compose


### Data Model
