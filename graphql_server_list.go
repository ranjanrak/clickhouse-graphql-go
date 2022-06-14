package clickhousegraphqlgo

import (
	"encoding/json"

	"github.com/graphql-go/graphql"
)

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
