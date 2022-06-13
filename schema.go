package clickhousegraphqlgo

import (
	"database/sql"
	"log"

	"github.com/graphql-go/graphql"
)

// createSchema creates tick graphql schema
func createSchema(dbConnect *sql.DB, resolverType string) (graphql.Schema, error) {
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
				Type:        fieldType(resolverType, tickType),
				Description: "Get tick detail",
				Args: graphql.FieldConfigArgument{
					"instrument_token": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					instrument_token, _ := params.Args["instrument_token"].(int)

					tickDataRef := &tickData{}
					rows, err := dbConnect.Query(`select instrument_token, 	
												  timestamp, 
												  last_price,
												  volume_traded,
												  oi
												  from tickstore where (instrument_token = ?)`, instrument_token)
					if err != nil {
						log.Fatalf("Error quering tickstore DB : %v", err)
					}
					defer rows.Close()

					// fetch tick data as per resolverType i.e list or single object
					if resolverType == "List" {
						return resolverList(rows, tickDataRef), nil
					} else {
						return resolverSingle(rows, tickDataRef), nil
					}
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	if err != nil {
		return graphql.Schema{}, err
	}

	return schema, nil
}

// resolverSingle return just the last tick object
func resolverSingle(rows *sql.Rows, tickDataRef *tickData) *tickData {
	for rows.Next() {
		if err := rows.Scan(&tickDataRef.InstrumentToken, &tickDataRef.Timestamp, &tickDataRef.LastPrice, &tickDataRef.VolumeTraded, &tickDataRef.OI); err != nil {
			log.Fatalf("failed to parse DB rows for resolverSingle : %v", err)
		}
	}
	return tickDataRef
}

// resolverList returns all list of available tick object
func resolverList(rows *sql.Rows, tickDataRef *tickData) []*tickData {
	// fetch all available tick data as list
	tickDataSum := make([]*tickData, 0)
	for rows.Next() {
		if err := rows.Scan(&tickDataRef.InstrumentToken, &tickDataRef.Timestamp, &tickDataRef.LastPrice, &tickDataRef.VolumeTraded, &tickDataRef.OI); err != nil {
			log.Fatalf("failed to parse DB rows for resolverList: %v", err)
		}
		tickDataSum = append(tickDataSum, tickDataRef)
	}
	return tickDataSum
}

// Select output field type based on input resolverType for different schema
func fieldType(resolverType string, tickType *graphql.Object) graphql.Output {
	if resolverType == "List" {
		return graphql.NewList(tickType)
	} else {
		return tickType
	}
}
