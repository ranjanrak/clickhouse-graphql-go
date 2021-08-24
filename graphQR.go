package main

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"net/http"
	"time"
	"log"
)

type TickData struct {
	InstrumentToken       int
	Timestamp             time.Time
	LastPrice             float64
	VolumeTraded          int
	OI                    int
}

func main() {
	// Use DSN as your clickhouse DB setup.
	// visit https://github.com/ClickHouse/clickhouse-go#dsn to know more
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")
	if err != nil {
		log.Fatal(err)
	}
	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
	}
	tickType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Tick",
		Description: "Tick Data",
		Fields: graphql.Fields{
			"instrument_token": &graphql.Field{
								Type:        graphql.NewNonNull(graphql.Int),
								Description: "Instrument token",
								Resolve: func(p graphql.ResolveParams) (interface{}, error) {
								if tick, ok := p.Source.(*TickData); ok {
								return tick.InstrumentToken, nil
								}
								return nil, nil
								},
							},
			"timestamp": &graphql.Field{
						 Type:        graphql.NewNonNull(graphql.DateTime),
						 Description: "Time stamp of tick",
						 Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						 if tick, ok := p.Source.(*TickData); ok {
							return tick.Timestamp, nil
						}

							return nil, nil
						},
					},
			"lastprice": &graphql.Field{
						 Type:        graphql.NewNonNull(graphql.Float),
						 Description: "Last Price of tick",
						 Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						 if tick, ok := p.Source.(*TickData); ok {
						 return tick.LastPrice, nil
						}

						return nil, nil
					},
				},
			"volumetraded": &graphql.Field{
							Type:        graphql.NewNonNull(graphql.Int),
							Description: "Total volume",
							Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							if tick, ok := p.Source.(*TickData); ok {
							return tick.VolumeTraded, nil
						}

						return nil, nil
					},
				},
			"oi": &graphql.Field{
				  Type:        graphql.NewNonNull(graphql.Int),
				  Description: "Net open interest",
				  Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				  if tick, ok := p.Source.(*TickData); ok {
				  return tick.OI, nil
				}

				return nil, nil
				},
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
				"Tick": &graphql.Field{
					Type:        tickType,
					Description: "Get tick detail",
					Args: graphql.FieldConfigArgument{
						"instrument_token": &graphql.ArgumentConfig{
											Type: graphql.Int,
						},
					},
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						instrument_token, _ := params.Args["instrument_token"].(int)

						tickData := &TickData{}
						rows, err := connect.Query(`select instrument_token, 	
												  timestamp, 
												  last_price,
												  volume_traded,
												  oi
												  from tickstore where (instrument_token = ?)`, instrument_token)
						if err != nil {
						log.Fatal(err)
						}
						defer rows.Close()
						// fetch latest tick data
						for rows.Next() {
							if err := rows.Scan(&tickData.InstrumentToken, &tickData.Timestamp, &tickData.LastPrice, &tickData.VolumeTraded, &tickData.OI); err != nil {
								log.Fatal(err)
							}
						}

						return tickData, nil
					},
				},
			},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
	})

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
	