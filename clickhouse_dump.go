package clickhousegraphqlgo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ClickHouse/clickhouse-go"

	kiteticker "github.com/zerodha/gokiteconnect/v3/ticker"
)

var (
	err        error
	token_list []uint32
	ticker     *kiteticker.Ticker
	dbClient   *sql.DB
)

// Create new client instance
func New(clientPar ClientParam) *Client {
	// Use DSN as your clickhouse DB setup.
	// visit https://github.com/ClickHouse/clickhouse-go#dsn to know more
	if clientPar.DBSource == "" {
		clientPar.DBSource = "tcp://127.0.0.1:9000?debug=true"
	}

	connect, err := sql.Open("clickhouse", clientPar.DBSource)
	if err = connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	schemaSingle, err := createSchema(connect, "single")
	if err != nil {
		log.Fatalf("failed to create single object schema, error: %v", err)
	}
	schemaList, err := createSchema(connect, "List")
	if err != nil {
		log.Fatalf("failed to create list of object schema, error: %v", err)
	}

	return &Client{
		dbClient:     connect,
		apiKey:       clientPar.ApiKey,
		accessToken:  clientPar.AccessToken,
		schemaSingle: &schemaSingle,
		schemaList:   &schemaList,
	}
}

// setDB creates tickstore table
func (c *Client) setDB() {
	_, err = c.dbClient.Exec(`
		CREATE TABLE IF NOT EXISTS tickstore (
			instrument_token       UInt32,
			timestamp              DateTime('Asia/Calcutta'),
			last_price             FLOAT(),
			average_traded_price   FLOAT(),
			volume_traded          UInt32,
			oi                     UInt32
		) engine=MergeTree()
		ORDER BY (timestamp)
	`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	// Close DB client once
	defer dbClient.Close()
	fmt.Println("Close: ", code, reason)
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {
	err := ticker.Subscribe(token_list)
	if err != nil {
		fmt.Println("err: ", err)
	}
	modeerr := ticker.SetMode("full", token_list)
	if modeerr != nil {
		fmt.Println("err: ", modeerr)
	}
}

// Triggered when tick is received
func onTick(tick kiteticker.Tick) {

	tx, err := dbClient.Begin()
	if err != nil {
		log.Fatalf("Error starting DB transaction: %v", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO tickstore (instrument_token, timestamp, last_price,
		average_traded_price, volume_traded, oi) VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatalf("Error creating sql statement: %v", err)
	}
	// Load tick data to DB
	if _, err := stmt.Exec(
		tick.InstrumentToken,
		tick.Timestamp.Time,
		tick.LastPrice,
		tick.AverageTradePrice,
		tick.VolumeTraded,
		tick.OI,
	); err != nil {
		log.Fatalf("Error executing a query: %v", err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Error committing the sql transaction: %v", err)
	}
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

// ClickhouseDump starts ticker and dumps tickdata to clickhouse
func (c *Client) ClickhouseDump(tokens []uint32) {
	token_list = tokens
	dbClient = c.dbClient

	// Perform DB setup
	c.setDB()

	// Create new Kite ticker instance
	ticker = kiteticker.New(c.apiKey, c.accessToken)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)

	// Start the connection
	ticker.Serve()
}
