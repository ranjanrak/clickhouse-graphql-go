package clickhousegraphqlgo

import (
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGraphqlQueryList(t *testing.T) {
	// Timestamp in time.Time object
	timestamp := time.Date(2022, 6, 8, 14, 04, 0, 0, time.Local)
	// Add mock row for test
	mockedRow := sqlmock.NewRows([]string{"instrument_token", "timestamp", "lastprice", "volumetraded", "oi"}).
		AddRow(60192519, timestamp, 1000, 12000, 12123).
		AddRow(10000000, timestamp, 9182, 15935, 12000)

	// Add expected query
	query := `select instrument_token, 	timestamp, last_price, volume_traded, oi from tickstore where (instrument_token = ?)`

	dbMock := setupMock(mockedRow, query)

	reqQuery := `query {
		Tick(instrument_token:61461767) {
		  instrument_token
		  timestamp
		  lastprice
		  volumetraded
		  oi
		}
	  }`

	listQueryStruct, err := dbMock.GraphqlQueryList(reqQuery)
	if err != nil {
		log.Fatalf("failed to execute object list graphql operation, errors: %+v", err)
	}

	assert.Equal(t, listQueryStruct.Output[1].InstrumentToken, 10000000, "Instrument token not matching of output struct")
	assert.Equal(t, listQueryStruct.Output[1].VolumeTraded, 15935, "Volume traded not matching of output struct")
	assert.Equal(t, listQueryStruct.Output[1].OI, 12000, "OI not matching of output struct")
	assert.Equal(t, listQueryStruct.Output[1].LastPrice, float64(9182), "LastPrice not matching of output struct")
}
