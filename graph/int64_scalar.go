package graph

import (
	"fmt"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalInt64(i int64) graphql.Marshaler {
	return graphql.MarshalString(strconv.FormatInt(i, 10))
}

func UnmarshalInt64(v interface{}) (int64, error) {
	switch val := v.(type) {
	case string:
		return strconv.ParseInt(val, 10, 64)
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	case float64:
		return int64(val), nil
	default:
		return 0, fmt.Errorf("invalid type for Int64: %T", v)
	}
}
