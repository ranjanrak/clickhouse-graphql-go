package clickhousegraphqlgo

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// Setup mockclient
func setupMock(mockRow *sqlmock.Rows, query string) *Client {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mock.ExpectQuery(query).
		WillReturnRows(mockRow)

	schemaSingle, err := createSchema(db, "single")
	if err != nil {
		log.Fatalf("failed to create new single object schema, error: %v", err)
	}

	schemaList, err := createSchema(db, "List")
	if err != nil {
		log.Fatalf("failed to create new list of object schema, error: %v", err)
	}

	cli := &Client{
		dbClient:     db,
		apiKey:       "your_api_key",
		accessToken:  "your_access_token",
		schemaSingle: &schemaSingle,
		schemaList:   &schemaList,
	}
	return cli
}

func TestGraphqlQuery(t *testing.T) {
	// Timestamp in time.Time object
	timestamp := time.Date(2022, 6, 8, 14, 04, 0, 0, time.Local)
	// Add mock row for test
	mockedRow := sqlmock.NewRows([]string{"instrument_token", "timestamp", "lastprice", "volumetraded", "oi"}).
		AddRow(60192519, timestamp, 9000, 15000, 12000)

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

	singleQueryStruct, err := dbMock.GraphqlQuery(reqQuery)
	if err != nil {
		log.Fatalf("failed to execute single object graphql operation, errors: %+v", err)
	}

	assert.Equal(t, singleQueryStruct.Output.InstrumentToken, 60192519, "Instrument token not matching of output struct")
	assert.Equal(t, singleQueryStruct.Output.VolumeTraded, 15000, "Volume traded not matching of output struct")
	assert.Equal(t, singleQueryStruct.Output.OI, 12000, "OI not matching of output struct")
	assert.Equal(t, singleQueryStruct.Output.LastPrice, float64(9000), "LastPrice not matching of output struct")
}
