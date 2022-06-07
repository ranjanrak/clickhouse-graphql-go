package clickhousegraphqlgo

import (
	"database/sql"
	"time"
)

// Client represents clickhouse DB client connection
type Client struct {
	dbClient    *sql.DB
	apiKey      string
	accessToken string
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
