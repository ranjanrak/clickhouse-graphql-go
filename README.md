# clickhouse-graphql-go

[![Run Tests](https://github.com/ranjanrak/clickhouse-graphql-go/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/ranjanrak/clickhouse-graphql-go/actions/workflows/go-test.yml)

GraphQL implementation for clickhouse in Go.
This package stores real time [streaming websocket data](https://kite.trade/docs/connect/v3/websocket/) in [clickhouse](https://clickhouse.tech/) and uses [GraphQL](https://graphql.org/) to consume the same.

## Installation

```
go get github.com/ranjanrak/clickhouse-graphql-go
```

## Usage

```go
import (
    clickhousegraphqlgo "github.com/ranjanrak/clickhouse-graphql-go"
)

// Create new graphql instance
client := clickhousegraphqlgo.New(clickhousegraphqlgo.ClientParam{
		DBSource:    "tcp://127.0.0.1:9000?debug=true",
		ApiKey:      "your_api_key",
		AccessToken: "your_access_token",
}))

// Dump tick websocket data to clickhouse
// Pass list of instrument token for subscription to websocket feeds
// Nothing will run after this
client.ClickhouseDump([]uint32{779521, 256265, 1893123, 13209858})

// Query
reqQuery := `query {
		Tick(instrument_token:779521) {
		  instrument_token
		  timestamp
		  lastprice
		  volumetraded
		  oi
		}
	}`

// Make single object schema graphql Query
singleQueryStruct, err := client.GraphqlQuery(reqQuery)
if err != nil {
  log.Fatalf("failed to execute single object graphql query, errors: %+v", err)
}
fmt.Printf("%+v\n", singleQueryStruct)

// Make list of object schema graphqlQuery
listQueryStruct, err := client.GraphqlQueryList(reqQuery)
if err != nil {
  log.Fatalf("failed to execute object list graphql query, errors: %+v", err)
}
fmt.Printf("%+v\n", listQueryStruct)

// Run graphql server on clickhouse with single schema
client.GraphqlServer("")

// Run graphql server to fetch list of object schema GraphQL
client.GraphqlServer("List")

```

#### GraphQL request query

1> `GraphqlQuery(reqQuery)`

```
reqQuery := `query {
		Tick(instrument_token:779521) {
		  instrument_token
		  timestamp
		  lastprice
		  volumetraded
		  oi
		}
	}`
```

2> `GraphqlQueryList(reqQuery)`

```
reqQuery := `query {
		Tick(instrument_token:779521) {
		  instrument_token
		  timestamp
		  lastprice
		  volumetraded
		  oi
		}
	}`
```

3> `GraphqlServer("")`

```

query {
  Tick(instrument_token:779521) {
    instrument_token
    timestamp
    lastprice
    volumetraded
    oi
  }
}

```

4> `GraphqlServer("List")`

```

query {
  Tick(instrument_token:779521) {
    instrument_token
    timestamp
    lastprice
    volumetraded
    oi
  }
}

```

## Response

1> `GraphqlQuery(reqQuery)`

```
{Output:{InstrumentToken:779521 LastPrice:463.4 OI:0 Timestamp:2022-06-07 17:12:48 +0530 IST
VolumeTraded:7672515}}
```

2> `GraphqlQueryList(reqQuery)`

```
{Output:[{InstrumentToken:779521 LastPrice:463.4 OI:0 Timestamp:2022-06-07 17:12:48 +0530 IST
VolumeTraded:7672515}
{InstrumentToken:779521 LastPrice:463.4 OI:0 Timestamp:2022-06-07 17:12:48 +0530 IST
VolumeTraded:7672515}
{InstrumentToken:779521 LastPrice:463.4 OI:0 Timestamp:2022-06-07 17:12:48 +0530 IST
VolumeTraded:7672515}
....
```

3> `GraphqlServer("")`

```

{
  "data": {
    "Tick": {
    "instrument_token": 1893123,
    "lastprice": 74.245,
    "oi": 1990638,
    "timestamp": "2021-08-24T16:38:39+05:30",
    "volumetraded": 1099802
    }
  }
}

```

4> `GraphqlServer("List")`

```
{
  "data": {
    "Tick": [
    {
    "instrument_token": 779521,
    "lastprice": 412.65,
    "oi": 0,
    "timestamp": "2021-08-26T12:19:09+05:30",
    "volumetraded": 7619425
    },
    {
    "instrument_token": 779521,
    "lastprice": 412.65,
    "oi": 0,
    "timestamp": "2021-08-26T12:19:09+05:30",
    "volumetraded": 7619425
    },
    ......]
  }
}

```

#### Sample query on graphiQL UI

1> `GraphqlServer("")`

![graphQL_dash](https://user-images.githubusercontent.com/29432131/130611805-cb60ba36-4e3e-4a24-8b56-722f0b8ef238.png)

2> `GraphqlServer("List")`

![graphQL_dash_list](https://user-images.githubusercontent.com/29432131/137927877-ccac9786-9695-447a-92fe-8c4744ea240c.png)

```

```
