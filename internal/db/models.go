package db

import (
	"fmt"
	"time"
)


type Product struct {
	ID         int
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}

func (p Product) Validate() error {
	if p.ID < 0 {
		return fmt.Errorf("id cannot be negative")
	}
	if p.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if p.Category == "" {
		return fmt.Errorf("category cannot be empty")
	}
	if p.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	return nil
}
