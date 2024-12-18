package db

import (
	"github.com/lib/pq"
)

func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == pq.ErrorCode("23505") {
			return true
		}
	}
	return false
}
