package seller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

// product and cart tables and handlers to be created
const (
	routePrefix   = "/seller"
	addProduct    = "/add-product"
	updateProduct = "/update-product"
	updateStatus  = "/update-seller-status"
	updateDetails = "/update-details"
	createSeller  = "/create-seller"
	sellerLogin   = "/seller-login"
)

type Seller struct {
	// SID          uint   `json:"sid"`
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	Rating       uint   `json:"rating"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	Status       string `json:"status"`
}

type productDetails struct {
	SellerId  string `json:"sid"`
	ProductId string `json:"name"`
}

func addProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo: replace product type
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		//we should use r.productdb not seller r.db
		if err := r.DB.Create(req); err != nil {
			log.Fatal("failed to update seller:", err)
		}

		c.SecureJSON(http.StatusOK, "product added")
	}
}

func updateProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := productDetails{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		//use product db here
		if err := r.DB.Model(&models.Product{}).Where("p_id = ?", req.ProductId).Updates(req).Error; err != nil {
			log.Fatal("failed to update seller:", err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func updateStatusHandler(r *util.Repository) gin.HandlerFunc {
	//what status to be u[dated] todo
	return func(c *gin.Context) {
		sellerNewDetails := Seller{}
		err := c.ShouldBindJSON(&sellerNewDetails)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		if err := r.DB.Model(&models.Seller{}).Where("email_address = ?", sellerNewDetails.EmailAddress).Updates(sellerNewDetails).Error; err != nil {
			log.Fatal("failed to update seller:", err)
		}
		c.SecureJSON(http.StatusOK, "updated seller status")
	}
}

func updateDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sellerNewDetails := Seller{}
		err := c.ShouldBindJSON(&sellerNewDetails)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		if err := r.DB.Model(&models.Seller{}).Where("email_address = ?", sellerNewDetails.EmailAddress).Updates(sellerNewDetails).Error; err != nil {
			log.Fatal("failed to update seller:", err)
		}
		c.SecureJSON(http.StatusAccepted, "updated seller details successfully")
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
	router.POST(createSeller, SellerSignUp(r))
	router.POST(sellerLogin, SellerLogin(r))
	v1 := router.Group(routePrefix)
	v1.Use(AuthoriseSeller(r))
	{
		v1.POST(addProduct, addProductHandler(r))
		v1.POST(updateProduct, updateProductHandler(r))
		v1.POST(updateStatus, updateStatusHandler(r))
		v1.POST(updateDetails, updateDetailsHandler(r))
	}
	return v1
}
