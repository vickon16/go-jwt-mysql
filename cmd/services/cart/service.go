package cart

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vickon16/go-jwt-mysql/cmd/types"
)

func getCartItemsIds(items []types.CartItem) ([]int, error) {
	ids := make([]int, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %v", item.ProductID)
		}

		ids = append(ids, item.ProductID)
	}

	return ids, nil
}

func (h *Handler) createOrder(products []types.Product, items []types.CartItem, userId uuid.UUID) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	// check if all products are available
	if err := checkCartsInStock(productMap, items); err != nil {
		return 0, 0, err
	}

	var totalPrice float64

	for _, item := range items {
		// calculate total price
		product := productMap[item.ProductID]
		totalPrice += (float64(item.Quantity) * product.Price)

		// reduce quantity of products in db
		product.Quantity -= item.Quantity
		if product.Quantity <= 0 {
			product.Quantity = 0
		}

		h.productStore.UpdateProduct(product.ID, types.UpdateProductPayload{
			Name:        product.Name,
			Description: product.Description,
			Image:       product.Image,
			Price:       product.Price,
			Quantity:    product.Quantity,
		})
	}

	// create the order
	orderId, err := h.orderStore.CreateOrder(types.CreateOrderPayload{
		UserID:  userId,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})

	if err != nil {
		return 0, 0, err
	}

	fmt.Println(orderId)

	// create order items
	for _, item := range items {
		err := h.orderStore.CreateOrderItem(types.CreateOrderItemPayload{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})

		if err != nil {
			return 0, 0, err
		}
	}

	return orderId, totalPrice, nil
}

func checkCartsInStock(productMap map[int]types.Product, items []types.CartItem) error {
	if len(items) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range items {
		product, ok := productMap[item.ProductID]
		if !ok {
			return fmt.Errorf("product %v not found", item.ProductID)
		}

		if item.Quantity > product.Quantity {
			return fmt.Errorf("product %v is out of stock", item.ProductID)
		}
	}

	return nil
}
