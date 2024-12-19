package db

import (
	"fmt"

	"project_sem/internal/config"
)

func (r *repositoriy) GetProducts(filter ProductsFilter) ([]Product, error) {
	products := make([]Product, 0)
	statement := fmt.Sprintf("SELECT id, name, category, price, create_date FROM prices %s", filter.ToSQLStmt())
	rows, err := r.db.Query(statement)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Category, &product.Price, &product.CreateDate)
		if err != nil {
			rows.Close()
			return nil, err
		}
		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repositoriy) CreateProduct(product Product) error {
	statement := fmt.Sprintf("INSERT INTO prices (id, name, category, price, create_date) VALUES (%d, '%s', '%s', %f, '%s')", product.ID, product.Name, product.Category, product.Price, product.CreateDate.Format(config.DateFormat))
	_, err := r.db.Exec(statement)
	if err != nil {
		return err
	}
	return nil
}

func (r *repositoriy) GetTotalPriceAndUniqueCategories() (float64, int, error) {
	var totalPrice float64
	var totalCategories int
	err := r.db.QueryRow("SELECT SUM(price), COUNT(DISTINCT category) FROM prices").Scan(&totalPrice, &totalCategories)
	if err != nil {
		return 0, 0, err
	}
	return totalPrice, totalCategories, nil
}
