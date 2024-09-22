package seller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	util "github.com/shopeeProject/shopee/util"
)

const (
	routePrefix   = "/seller"
	addProduct    = "/add-product"
	updateProduct = "/update-product"
	updateStatus  = "/update-seller-status"
	updateDetails = "/update-details"
)

type Seller struct {
	SellerId    string `json:"sid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Rating      uint   `json:"rating"`
	Description string
	Image       string
	Status      string
}

type productDetails struct {
	SellerId  string `json:"sid"`
	ProductId string `json:"name"`
}

func addProductHandler(v string) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func updateProductHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func updateStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func updateDetailsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func AuthoriseSeller(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramPairs := c.Request.URL.Query()
		for key, values := range paramPairs {
			fmt.Printf("key = %v, value(s) = %v\n", key, values)
		}
		c.Next()
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	v1 := router.Group(routePrefix)
	v1.Use(AuthoriseSeller(r))
	{
		v1.POST(addProduct, addProductHandler("hi"))
		v1.POST(updateProduct, updateProductHandler())
		v1.POST(updateStatus, updateStatusHandler())
		v1.POST(updateDetails, updateDetailsHandler())
	}
	return v1
}
