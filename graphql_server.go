package clickhousegraphqlgo

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GraphqlServer runs graphql query server for singleResolver or listResolver on localhost:8080
func (c *Client) GraphqlServer(reqSchemaType string) {
	schemaType := c.schemaSingle
	if reqSchemaType == "List" {
		schemaType = c.schemaList
	}
	h := handler.New(&handler.Config{
		Schema:   schemaType,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}

// GraphqlQuery returns graphql query output performed on the tick schema
func (c *Client) GraphqlQuery(reqQuery string) (Result, error) {

	params := graphql.Params{Schema: *c.schemaSingle, RequestString: reqQuery}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return Result{}, err
	}

	rJSON, _ := json.Marshal(r.Data)

	var result Result
	if err := json.Unmarshal(rJSON, &result); err != nil {
		return Result{}, err
	}

	return result, nil
}
