package product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/vickon16/go-jwt-mysql/cmd/services/auth"
	"github.com/vickon16/go-jwt-mysql/cmd/types"
	"github.com/vickon16/go-jwt-mysql/cmd/utils"
)

type Handler struct {
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{productStore, userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", auth.WithJWTAuth(h.handleCreateProducts, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{id}", auth.WithJWTAuth(h.handleUpdateProduct, h.userStore)).Methods(http.MethodPatch)
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	// get Json payload
	products, err := h.productStore.GetProducts()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleCreateProducts(w http.ResponseWriter, r *http.Request) {
	// get Json payload
	var payload types.CreateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate the payloads
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	err := h.productStore.CreateProduct(types.CreateProductPayload{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["id"]
	id, err := strconv.Atoi(params)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get Json payload
	var payload types.UpdateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// get product
	product, err := h.productStore.GetProductById(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if payload.Name == "" {
		payload.Name = product.Name
	}
	if payload.Description == "" {
		payload.Description = product.Description
	}
	if payload.Image == "" {
		payload.Image = product.Image
	}
	if payload.Price == 0 {
		payload.Price = product.Price
	}
	if payload.Quantity == 0 {
		payload.Quantity = product.Quantity
	}
	err = h.productStore.UpdateProduct(product.ID, payload)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, payload)
}
