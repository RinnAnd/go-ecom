package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    string
	Name  string
	Price int64
	Brand string
}

type ProductService struct {
	DB *sql.DB
}

type ServerResponse struct {
	Message string
	Prod    *Product
}

type ProductRepo interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
}

func (ps *ProductService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	_, err = ps.DB.Exec("INSERT INTO products (name, price, brand) VALUES ($1, $2, $3)", product.Name, product.Price, product.Brand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(ServerResponse{
		Message: "Product created",
		Prod:    &product,
	})
}

func (ps *ProductService) GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := ps.DB.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer rows.Close()

	products := []Product{}
	for rows.Next() {
		item := Product{}
		err := rows.Scan(&item.ID, &item.Name, &item.Price, &item.Brand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		products = append(products, item)
	}
	json.NewEncoder(w).Encode(products)
}

func (ps *ProductService) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := ps.DB.Query("SELECT * FROM products WHERE id = $1", id)
	product := Product{}
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Brand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	if product.ID == "" {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func (ps *ProductService) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := ps.DB.Query("SELECT * FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	currProd := Product{}
	for rows.Next() {
		err := rows.Scan(&currProd.ID, &currProd.Name, &currProd.Price, &currProd.Brand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
	newProd := Product{}
	json.NewDecoder(r.Body).Decode(&newProd)

	if newProd.Name == "" {
		newProd.Name = currProd.Name
	}
	if newProd.Brand == "" {
		newProd.Brand = currProd.Name
	}
	if newProd.Price == 0 {
		newProd.Price = currProd.Price
	}
	_, err = ps.DB.Exec("UPDATE products SET name = $2, price = $3, brand = $4 WHERE id = $1", id, newProd.Name, newProd.Price, newProd.Brand)
	if err != nil {
		http.Error(w, "Could not update the product", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(ServerResponse{
		Message: "Product updated",
		Prod:    &newProd,
	})
}

func (ps *ProductService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := ps.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	response := fmt.Sprintf("Deleted product with ID: %v", id)
	json.NewEncoder(w).Encode(response)
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{
		DB: db,
	}
}
