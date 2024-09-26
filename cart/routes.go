package cart

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

type Cart struct {
	UID   uint ` json:"uid"`
	PID   uint `json:"pid"`
	Count uint `json:"count"`
}

// product and cart tables and handlers to be created
const (
	routePrefix    = "/cart"
	getCartDetails = "/get-cart-details"
	deleteFromCart = "/delete-from-cart"
	checkout       = "/checkout"
	clearCart      = "/clear-cart"
	createCart     = "/create-cart"
	addToCart      = "/add-to-cart"
)

func getCartDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req uint
		if err := c.ShouldBindJSON(&req); err != nil {
			c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := r.DB.Where("uid = ?", req).Delete(&models.Cart{}).Error; err != nil {
			log.Println("failed to delete records:", err)
			c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete records"})
			return
		}
		c.SecureJSON(http.StatusAccepted, gin.H{"message": "records deleted successfully"})
	}
}

func checkoutHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo
		req := models.Cart{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func deleteFromCartHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.Cart
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//todo:reduce count
		// Delete the specific product from the cart
		if err := r.DB.Where("uid = ? AND pid = ?", req.UID, req.PID).Delete(&models.Cart{}).Error; err != nil {
			log.Println("failed to delete product from cart:", err)
			c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
			return
		}

		c.SecureJSON(http.StatusAccepted, gin.H{"message": "product deleted successfully"})
	}
}

func clearCartHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := Cart{}
		err := c.ShouldBindJSON(&req)

		if err := r.DB.Where("uid = ?", req.UID).Delete(&models.Cart{}).Error; err != nil {
			log.Println("failed to clear cart:", err)
			c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
			return
		}
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, "updated details successfully")
	}
}

func addToCartHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := Cart{}
		err := c.ShouldBindJSON(&req)

		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}

		// Check if the product already exists in the cart
		var existingCart models.Cart
		if err := r.DB.Where("uid = ? AND pid = ?", req.UID, req.PID).First(&existingCart).Error; err == nil {
			// Product exists, update the count to 1
			existingCart.Count = 1
			if err := r.DB.Save(&existingCart).Error; err != nil {
				log.Println("failed to update cart count:", err)
				c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to update cart"})
				return
			}
			c.SecureJSON(http.StatusAccepted, gin.H{"message": "cart updated successfully"})
			return
		}

		// Product does not exist, create a new record with count = 1
		req.Count = 1 // Set count to 1 for the new record
		if err := r.DB.Create(&req).Error; err != nil {
			log.Println("failed to add product to cart:", err)
			c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to add product"})
			return
		}
		c.SecureJSON(http.StatusAccepted, "cart created successfully")
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	// router.POST(createCart, SellerSignUp(r))
	v1 := router.Group(routePrefix)
	{
		v1.POST(getCartDetails, getCartDetailsHandler(r))
		v1.POST(deleteFromCart, deleteFromCartHandler(r))
		v1.POST(checkout, checkoutHandler(r))
		v1.POST(clearCart, clearCartHandler(r))
		v1.POST(addToCart, addToCartHandler(r))
	}
	return v1
}
