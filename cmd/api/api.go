package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/vickon16/go-jwt-mysql/cmd/services/cart"
	"github.com/vickon16/go-jwt-mysql/cmd/services/order"
	"github.com/vickon16/go-jwt-mysql/cmd/services/product"
	"github.com/vickon16/go-jwt-mysql/cmd/services/user"
	"github.com/vickon16/go-jwt-mysql/cmd/utils"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewAPIServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address,
		db,
	}
}

func (s *APIServer) Run() error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your React app
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Methods you want to allow
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Headers you want to allow
	})

	// using the gorilla mux router
	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "Server is running"})
	}).Methods("GET")

	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	user.NewHandler(userStore).RegisterRoutes(subRouter)

	productStore := product.NewStore(s.db)
	product.NewHandler(productStore, userStore).RegisterRoutes(subRouter)

	orderStore := order.NewStore(s.db)

	cart.NewHandler(orderStore, productStore, userStore).RegisterRoutes(subRouter)

	log.Println("Listening on", s.address)
	return http.ListenAndServe(s.address, c.Handler(router))
}
