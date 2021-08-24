# clickhouse-graphql-go
GraphQL implementation for clickhouse in Go. 
This package stores real time [streaming websocket data](https://kite.trade/docs/connect/v3/websocket/) in [clickhouse](https://clickhouse.tech/) and uses [GraphQL](https://graphql.org/) to consume the same.   

## Installation
```
git clone https://github.com/ranjanrak/clickhouse-graphql-go.git
cd clickhouse-graphql-go
```

## Usage
#### Dump data to clickhouse
```
go run clickhouseQR.go
```
#### Run local graphql server
```
go run graphQR.go
```
#### GraphQL query
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

## Response
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
#### Sample query on graphiQL
![graphQL_dash](https://user-images.githubusercontent.com/29432131/130611805-cb60ba36-4e3e-4a24-8b56-722f0b8ef238.png)