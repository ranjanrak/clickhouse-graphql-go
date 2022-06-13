package clickhousegraphqlgo

import (
	"database/sql"
	"time"

	"github.com/graphql-go/graphql"
)

// Client represents clickhouse DB client connection
type Client struct {
	dbClient     *sql.DB
	apiKey       string
	accessToken  string
	schemaSingle *graphql.Schema
	schemaList   *graphql.Schema
}

// ClientParam represents interface to connect clickhouse and kite ticker stream
type ClientParam struct {
	DBSource    string
	ApiKey      string
	AccessToken string
}

// tickData represents fields of streaming ticks
type tickData struct {
	InstrumentToken int
	Timestamp       time.Time
	LastPrice       float64
	VolumeTraded    int
	OI              int
}

// tickOutput represents graphql output json fields
type tickOutput struct {
	InstrumentToken int       `json:"instrument_token"`
	LastPrice       float64   `json:"lastprice"`
	OI              int       `json:"oi"`
	Timestamp       time.Time `json:"timestamp"`
	VolumeTraded    int       `json:"volumetraded"`
}

// Result represents graphql output struct
type Result struct {
	Output tickOutput `json:"Tick"`
}

// ResultList represents graphql tick list schema output struct
type ResultList struct {
	Output []tickOutput `json:"Tick"`
}
