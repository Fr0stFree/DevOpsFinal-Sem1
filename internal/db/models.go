package db

import (
	"errors"
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
		return errors.New("id cannot be negative")
	}
	if p.Name == "" {
		return errors.New("name cannot be empty")
	}
	if p.Category == "" {
		return errors.New("category cannot be empty")
	}
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}
