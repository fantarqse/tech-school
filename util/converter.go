package util

import (
	"database/sql"

	"golang.org/x/exp/constraints"
)

func SQLNullInt64[T constraints.Integer](intenger T) sql.NullInt64 {
	return sql.NullInt64{
		Int64: int64(intenger),
		Valid: true,
	}
}
