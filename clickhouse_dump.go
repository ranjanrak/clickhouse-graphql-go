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
	connect    *sql.DB
	err        error
	token_list []uint32
)

func setDB() {

	// Use DSN as your clickhouse DB setup.
	// visit https://github.com/ClickHouse/clickhouse-go#dsn to know more
	connect, err = sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")

	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}

	_, err = connect.Exec(`
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
		log.Fatal(err)
	}
}

var (
	ticker *kiteticker.Ticker
)

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("Error: ", err)
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	defer connect.Close()
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

// Triggered when tick is recevived
func onTick(tick kiteticker.Tick) {

	fmt.Printf("%+v\n", tick)

	tx, err := connect.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO tickstore (instrument_token, timestamp, last_price,
		average_traded_price, volume_traded, oi) VALUES (?, ?, ?, ?, ?, ?)`)

	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
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

func ClickhouseDump(tokens []uint32) {
	apiKey := "your_api_key"
	accessToken := "your_access_token"
	token_list = tokens

	// Perform DB related part
	setDB()

	// Create new Kite ticker instance
	ticker = kiteticker.New(apiKey, accessToken)

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
