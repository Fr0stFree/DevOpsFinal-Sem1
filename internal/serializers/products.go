package serializers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"project_sem/internal/config"
	"project_sem/internal/db"
)

func SerializeProduct(prices []db.Product) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)
	defer csvWriter.Flush()

	csvWriter.Write([]string{"id", "name", "category", "price", "create_date"})
	for _, price := range prices {
		record := []string{
			fmt.Sprintf("%d", price.ID),
			price.Name,
			price.Category,
			fmt.Sprintf("%.2f", price.Price),
			price.CreateDate.Format(config.DateFormat),
		}
		if err := csvWriter.Write(record); err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func DeserializeProducts(r io.Reader) ([]db.Product, int, error) {
	products := make([]db.Product, 0)
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, 0, err
	}
	totalItems := len(records) - 1
	for idx, record := range records {
		if idx == 0 {
			continue
		}
		product, err := buildProductFromRecord(record)
		if err != nil {
			log.Printf("product conversion failed %v\n", err)
			continue
		}
		if err = product.Validate(); err != nil {
			log.Printf("product validation failed %v\n", err)
			continue
		}
		products = append(products, product)
	}
	return products, totalItems, nil
}

func buildProductFromRecord(record []string) (db.Product, error) {
	product := db.Product{}

	id, err := strconv.Atoi(record[0])
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert id to int %v", err)
	}
	product.ID = id
	product.Name = record[1]
	product.Category = record[2]
	price, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert price to float %v", err)
	}
	product.Price = price
	create_date, err := time.Parse(config.DateFormat, record[4])
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert creation date %v", err)
	}
	product.CreateDate = create_date
	return product, nil
}
