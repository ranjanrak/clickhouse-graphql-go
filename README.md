# clickhouse-graphql-go
GraphQL implementation for clickhouse in Go. 
This package stores real time [streaming websocket data](https://kite.trade/docs/connect/v3/websocket/) in [clickhouse](https://clickhouse.tech/) and uses [GraphQL](https://graphql.org/) to consume the same.   

## Installation
```
go get github.com/ranjanrak/clickhouse-graphql-go
```

## Usage
```
import (
    clickhousegraphqlgo "github.com/ranjanrak/clickhouse-graphql-go"
)

// Dump tick websocket data to clickhouse
// Pass list of instrument token for subscription to websocket feeds
clickhousegraphqlgo.ClickhouseDump([]uint32{779521, 256265, 1893123, 13209858})

// Run graphql server on clickhouse
clickhousegraphqlgo.GraphqlServer()

// Run graphql server to fetch GraphQL List
clickhousegraphqlgo.GraphqlServerList()

```
#### GraphQL query
1> `GraphqlServer()`
```
query {
  Tick(instrument_token:1893123) {
    instrument_token
    timestamp
    lastprice
    volumetraded
    oi
  }
}
```
2> `GraphqlServerList()`
```
{
  Tick(instrument_token: 779521) {
    instrument_token
    timestamp
    lastprice
    volumetraded
    oi
  }
}
```

## Response
1> `GraphqlServer()`
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
2> `GraphqlServerList()`
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
      {
        "instrument_token": 779521,
        "lastprice": 412.65,
        "oi": 0,
        "timestamp": "2021-08-26T12:19:09+05:30",
        "volumetraded": 7619425
      },
      ...
      ...
      {
        "instrument_token": 779521,
        "lastprice": 412.65,
        "oi": 0,
        "timestamp": "2021-08-26T12:19:09+05:30",
        "volumetraded": 7619425
      },
      ......
    ]
  }
}
```
#### Sample query on graphiQL UI
1> `GraphqlServer()`

![graphQL_dash](https://user-images.githubusercontent.com/29432131/130611805-cb60ba36-4e3e-4a24-8b56-722f0b8ef238.png)

2> `GraphqlServerList()`

![graphQL_dash_list](https://user-images.githubusercontent.com/29432131/137927877-ccac9786-9695-447a-92fe-8c4744ea240c.png)