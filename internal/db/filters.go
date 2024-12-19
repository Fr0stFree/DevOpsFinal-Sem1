package db

import (
	"fmt"
	"strings"
	"time"

	"project_sem/internal/config"
)

type ProductsFilter struct {
	MinPrice      float64
	MaxPrice      float64
	MinCreateDate time.Time
	MaxCreateDate time.Time
}

func (p ProductsFilter) ToSQLStmt() string {
	conditions := make([]string, 0)
	if p.MinPrice != 0 {
		conditions = append(conditions, fmt.Sprintf("price >= %f", p.MinPrice))
	}
	if p.MaxPrice != 0 {
		conditions = append(conditions, fmt.Sprintf("price <= %f", p.MaxPrice))
	}
	if !p.MinCreateDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("create_date >= '%s'", p.MinCreateDate.Format(config.DateFormat)))
	}
	if !p.MaxCreateDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("create_date <= '%s'", p.MaxCreateDate.Format(config.DateFormat)))
	}
	if len(conditions) != 0 {
		return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
	}
	return ""
}
