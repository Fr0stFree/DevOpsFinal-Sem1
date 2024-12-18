package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"project_sem/internal/config"
)

type Repositorier interface {
	GetProducts(ProductsFilter) ([]Product, error)
	CreateProduct(Product) error
	GetTotalPriceAndUniqueCategories() (float64, int, error)
	Begin() (Transactioner, error)
}

type repositoriy struct {
	db *sql.DB
}

func NewRepository(cfg config.DBConfig) (Repositorier, error) {
	log.Println("connecting to database...")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &repositoriy{}, err
	}
	if err = db.Ping(); err != nil {
		return &repositoriy{}, err
	}
	log.Printf("successfully connected to database '%s'\n", cfg.Name)
	return &repositoriy{db}, nil
}

type Transactioner interface {
	Commit() error
	Rollback() error
}

func (r *repositoriy) Begin() (Transactioner, error) {
	return r.db.Begin()
}
