package order

import (
	"database/sql"

	"github.com/vickon16/go-jwt-mysql/cmd/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) CreateOrder(order types.CreateOrderPayload) (int, error) {
	res, err := s.db.Exec(`INSERT INTO orders (userId, total, status, address) VALUES (?, ?, ?, ?)`, order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderItem types.CreateOrderItemPayload) error {
	_, err := s.db.Exec(`INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?, ?, ?, ?)`, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err
}
