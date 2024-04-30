package server

import (
	"database/sql"
	product "ecom/pkg/services/product"
	user "ecom/pkg/services/user"
	"fmt"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	db         *sql.DB
	gateWay    *Gateway
}

type Gateway struct {
	ProductService *product.ProductService
	UserService    *user.UserService
}

func NewGateway(ps *product.ProductService, us *user.UserService) *Gateway {
	return &Gateway{
		ProductService: ps,
		UserService:    us,
	}
}

func (s *Server) CreateTables() error {
	_, err := s.db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id uuid DEFAULT uuid_generate_v4 (),
		name varchar(255),
		price int,
		brand varchar(255),
		PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}
	//! DONT FORGET TO HASH THE PASSWORDS
	_, err = s.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4 (),
		name varchar(255),
		email int,
		role varchar(255),
		password varchar(255),
		PRIMARY KEY (id)
		)`)
	if err != nil {
		return err
	}
	return nil
}

func NewServer(addr string, db *sql.DB) *Server {
	productService := product.NewProductService(db)
	userService := user.NewUserService(db)
	return &Server{
		httpServer: &http.Server{
			Addr: addr,
		},
		db:      db,
		gateWay: NewGateway(productService, userService),
	}
}

func (s *Server) Start() error {
	fmt.Println("[APP] is now up and listening on port", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) MakeProduct(w http.ResponseWriter, r *http.Request) {
	s.gateWay.ProductService.CreateProduct(w, r)
}

func (s *Server) GetProducts(w http.ResponseWriter, r *http.Request) {
	s.gateWay.ProductService.GetProducts(w, r)
}

func (s *Server) GetProduct(w http.ResponseWriter, r *http.Request) {
	s.gateWay.ProductService.GetProduct(w, r)
}

func (s *Server) EditProduct(w http.ResponseWriter, r *http.Request) {
	s.gateWay.ProductService.UpdateProduct(w, r)
}

func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	s.gateWay.ProductService.DeleteProduct(w, r)
}
