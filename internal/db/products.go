package db

import (
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID         int
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}

type FilterParams struct {
	MinPrice      float64
	MaxPrice      float64
	MinCreateDate time.Time
	MaxCreateDate time.Time
}

func (r *Repository) GetProducts(params FilterParams) ([]Product, error) {
	prices := make([]Product, 0)
	conditions := make([]string, 0)
	statement := "SELECT id, name, category, price, create_date FROM prices"
	if params.MinPrice != 0 {
		conditions = append(conditions, fmt.Sprintf("price >= %f", params.MinPrice))
	}
	if params.MaxPrice != 0 {
		conditions = append(conditions, fmt.Sprintf("price <= %f", params.MaxPrice))
	}
	if !params.MinCreateDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("create_date >= '%s'", params.MinCreateDate.Format("2006-01-02")))
	}
	if !params.MaxCreateDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("create_date <= '%s'", params.MaxCreateDate.Format("2006-01-02")))
	}
	if len(conditions) != 0 {
		statement = fmt.Sprintf("%s WHERE %s", statement, strings.Join(conditions, " AND "))
	}
	rows, err := r.db.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.CreateDate)
		if err != nil {
			return nil, err
		}
		prices = append(prices, product)
	}
	return prices, nil
}

func (r *Repository) CreateProduct(product Product) error {
	statement := fmt.Sprintf("INSERT INTO prices (id, name, category, price, create_date) VALUES (%d, '%s', '%s', %f, '%s')", product.ID, product.Name, product.Category, product.Price, product.CreateDate.Format("2006-01-02"))
	_, err := r.db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTotalPriceAndUniqueCategories() (float64, int, error) {
	var totalPrice float64
	var totalCategories int
	err := r.db.QueryRow("SELECT SUM(price), COUNT(DISTINCT category) FROM prices").Scan(&totalPrice, &totalCategories)
	if err != nil {
		return 0, 0, err
	}
	return totalPrice, totalCategories, nil
}
