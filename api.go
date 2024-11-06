package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	cart "github.com/shopeeProject/shopee/cart"
	category "github.com/shopeeProject/shopee/category"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	order "github.com/shopeeProject/shopee/order"
	product "github.com/shopeeProject/shopee/product"
	seller "github.com/shopeeProject/shopee/seller"
	user "github.com/shopeeProject/shopee/user"
	util "github.com/shopeeProject/shopee/util"
	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr string
}

func (s *APIServer) Run(db *gorm.DB) {
	router := gin.Default()

	user.RegisterRoutes(router, &util.Repository{DB: db})
	seller.RegisterRoutes(router, &util.Repository{DB: db})
	cart.RegisterRoutes(router, &util.Repository{DB: db})
	order.RegisterRoutes(router, &util.Repository{DB: db})
	category.RegisterRoutes(router, &util.Repository{DB: db})
	product.RegisterRoutes(router, &util.Repository{DB: db})
	jwthandler.RegisterRoutes(router, &util.Repository{DB: db})

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router) // starts http server on on address specified and listens for incoming requests
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}

}
