package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	cart "github.com/shopeeProject/shopee/cart"
	category "github.com/shopeeProject/shopee/category"
	order "github.com/shopeeProject/shopee/order"
	seller "github.com/shopeeProject/shopee/seller"
	user "github.com/shopeeProject/shopee/user"
	util "github.com/shopeeProject/shopee/util"
)

type APIServer struct {
	listenAddr string
}

func (s *APIServer) Run(r *util.ShopeeDatabase) {
	router := gin.Default()

	user.GroupUserRoutes(router, &util.Repository{DB: r.UserDB})
	seller.RegisterRoutes(router, &util.Repository{DB: r.SellerDB})
	cart.RegisterRoutes(router, &util.Repository{DB: r.CartDB})
	order.RegisterRoutes(router, &util.Repository{DB: r.OrderDB})
	category.RegisterRoutes(router, &util.Repository{DB: r.CategoryDB})

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router) // starts http server on on address specified and listens for incoming requests
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}

}
