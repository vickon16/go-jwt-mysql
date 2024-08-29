package product

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vickon16/go-jwt-mysql/cmd/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetProductById(id int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	product := new(types.Product)
	for rows.Next() {
		product, err = scanRowIntoProduct(rows, product)
		if err != nil {
			return nil, err
		}
	}

	if product.ID == 0 {
		return nil, fmt.Errorf("product not found")
	}

	return product, nil
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]types.Product, 0)
	product := new(types.Product)

	for rows.Next() {
		product, err = scanRowIntoProduct(rows, product)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (s *Store) GetProductsByIds(productIds []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIds)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	// convert productIds to []interface{}
	args := make([]interface{}, len(productIds))
	for i, id := range productIds {
		args[i] = id
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]types.Product, 0)
	product := new(types.Product)

	for rows.Next() {
		product, err = scanRowIntoProduct(rows, product)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.CreateProductPayload) error {
	_, err := s.db.Exec(`INSERT INTO products (name, price, description, quantity, image) VALUES (?, ?, ?, ?, ?)`, product.Name, product.Price, product.Description, product.Quantity, product.Image)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateProduct(productId int, product types.UpdateProductPayload) error {
	_, err := s.db.Exec(`UPDATE products SET name = ?, price = ?, description = ?, quantity = ?, image = ? WHERE id = ?`, product.Name, product.Price, product.Description, product.Quantity, product.Image, productId)
	return err
}

func scanRowIntoProduct(rows *sql.Rows, product *types.Product) (*types.Product, error) {
	err := rows.Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Name,
		&product.Image,
		&product.Description,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
