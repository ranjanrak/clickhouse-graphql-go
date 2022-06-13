package clickhousegraphqlgo

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// GraphqlServerList runs graphql query server for resolverList on localhost:8080
func (c *Client) GraphqlServerList() {
	h := handler.New(&handler.Config{
		Schema:   c.schemaList,
		Pretty:   true,
		GraphiQL: true,
	})

	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}

// GraphqlQueryList returns graphql query output performed on the tick list schema
func (c *Client) GraphqlQueryList(reqQuery string) (ResultList, error) {

	params := graphql.Params{Schema: *c.schemaList, RequestString: reqQuery}

	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		return ResultList{}, err
	}

	rJSON, _ := json.Marshal(r.Data)

	var result ResultList
	if err := json.Unmarshal(rJSON, &result); err != nil {
		return ResultList{}, err
	}

	return result, nil
}
