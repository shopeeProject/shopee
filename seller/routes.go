package seller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/category"
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
	// SID          int   `json:"sid"`
	Name         string `json:"name"`
	EmailAddress string `json:"email"`
	Password     string `json:"password"`
	Rating       int    `json:"rating"`
	Description  string `json:"description"`
	CategoryID   int    `json:"category"`
	Image        string `json:"image"`
	Status       string `json:"status"`
	IsApproved   bool   `json:"isApproved"`
}

type productDetails struct {
	// PID          int
	Name         string `json:"name"`
	Price        string `json:"price"`
	Availability bool   `json:"availability"`
	Rating       string `json:"rating"`
	CategoryID   int    `json:"category"`
	Description  string `json:"description"`
	SID          string `json:"sid"`
	Image        string `json:"image"`
}

func addProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		req := productDetails{}
		err := c.ShouldBindJSON(&req)

		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}

		validationResponse := category.ValidateCategory(r.DB, req.CategoryID)
		if validationResponse.Success == false {
			c.SecureJSON(http.StatusBadRequest, gin.H{
				"message": validationResponse.Message,
			})
		}

		//if valid category add product to product table
		if err := r.DB.Create(&req); err != nil {
			log.Fatal("failed to update seller:", err)
		}

		c.SecureJSON(http.StatusOK, "product added")
	}
}

func updateProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		req := models.Product{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}

		if err := r.DB.Model(&models.Product{}).Where("p_id = ?", req.PID).Updates(req).Error; err != nil {
			log.Fatal("failed to update seller:", err)
		}
		c.SecureJSON(http.StatusAccepted, "updated product details")
	}
}

func updateStatusHandler(r *util.Repository) gin.HandlerFunc {
	//todo what status to be updated
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

// todo
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
		v1.PATCH(updateProduct, updateProductHandler(r))
		v1.PATCH(updateStatus, updateStatusHandler(r))
		v1.PATCH(updateDetails, updateDetailsHandler(r))
	}
	return v1
}
