package clickhousegraphqlgo

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func (c *Client) GraphqlServerList() {
	tickType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Tick",
		Description: "Tick Data",
		Fields: graphql.Fields{
			"instrument_token": &graphql.Field{
				Type:        graphql.Int,
				Description: "Instrument token",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if tick, ok := p.Source.(*tickData); ok {
						return tick.InstrumentToken, nil
					}
					return nil, nil
				},
			},
			"timestamp": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "Time stamp of tick",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if tick, ok := p.Source.(*tickData); ok {
						return tick.Timestamp, nil
					}

					return nil, nil
				},
			},
			"lastprice": &graphql.Field{
				Type:        graphql.Float,
				Description: "Last Price of tick",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if tick, ok := p.Source.(*tickData); ok {
						return tick.LastPrice, nil
					}

					return nil, nil
				},
			},
			"volumetraded": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total volume",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if tick, ok := p.Source.(*tickData); ok {
						return tick.VolumeTraded, nil
					}

					return nil, nil
				},
			},
			"oi": &graphql.Field{
				Type:        graphql.Int,
				Description: "Net open interest",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if tick, ok := p.Source.(*tickData); ok {
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
				Type:        graphql.NewList(tickType),
				Description: "Get tick detail",
				Args: graphql.FieldConfigArgument{
					"instrument_token": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					instrument_token, _ := params.Args["instrument_token"].(int)

					tickDataRef := &tickData{}
					rows, err := c.dbClient.Query(`select instrument_token, 	
												  timestamp, 
												  last_price,
												  volume_traded,
												  oi
												  from tickstore where (instrument_token = ?)`, instrument_token)
					if err != nil {
						log.Fatal(err)
					}
					defer rows.Close()
					// fetch all available tick data as list
					tickDataSum := make([]*tickData, 0)
					for rows.Next() {
						if err := rows.Scan(&tickDataRef.InstrumentToken, &tickDataRef.Timestamp, &tickDataRef.LastPrice, &tickDataRef.VolumeTraded, &tickDataRef.OI); err != nil {
							log.Fatal(err)
						}
						tickDataSum = append(tickDataSum, tickDataRef)
					}

					return tickDataSum, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
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
