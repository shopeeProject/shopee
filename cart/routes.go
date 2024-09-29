package cart

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/order"
	"github.com/shopeeProject/shopee/product"
	util "github.com/shopeeProject/shopee/util"
)

type Cart struct {
	UID   int ` json:"uid"`
	PID   int `json:"pid"`
	Count int `json:"count"`
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

func GetProductIDsOfUser(r *util.Repository, userID int) ([]int, error) {
	var productIDs []int
	err := r.DB.Model(&Cart{}).Where("uid = ?", userID).Pluck("pid", &productIDs).Error
	if err != nil {
		return nil, err
	}
	return productIDs, err
}

func checkoutHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//todo

		req := models.Cart{}
		err := c.ShouldBindJSON(&req)

		productIDs, err := GetProductIDsOfUser(r, req.UID)
		if err != nil {
			c.SecureJSON(http.StatusInternalServerError, "error fetching product ids")
		}
		productsList,err := product.GetProductDetails(r, productIDs)
		order.PlaceOrderHandler1(r, req.UID, productIDs, productsList)
		//exp code above
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		c.SecureJSON(http.StatusAccepted, req)
	}
}

func getCartDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		if uid, userIdPresent := c.GetQuery("userid"); userIdPresent {
			var cartDetails models.Cart
			if err := r.DB.Where("uid = ?", uid).Find(&cartDetails).Error; err != nil {
				c.SecureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.SecureJSON(http.StatusOK, cartDetails)
			return
		}
		c.SecureJSON(http.StatusBadRequest, "userid param is missing")

	}
}

func deleteFromCartHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		//todo test this for both count>0 & count=0
		uid, userIdPresent := c.GetQuery("userid")
		pid, productIdPresent := c.GetQuery("pid")

		if userIdPresent && productIdPresent {
			if err := r.DB.Where("uid = ? AND pid = ?", uid, pid).Delete(&models.Cart{}).Error; err != nil {
				c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
				return
			}
			c.SecureJSON(http.StatusAccepted, gin.H{"message": "product deleted successfully"})
			return
		}
		c.SecureJSON(http.StatusBadRequest, " expected params missing")

	}
}

func clearCartHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid, userIdPresent := c.GetQuery("userid")

		if userIdPresent {
			if err := r.DB.Where("uid = ?", uid).Delete(&models.Cart{}).Error; err != nil {
				c.SecureJSON(http.StatusInternalServerError, gin.H{"error": "failed to clear cart"})
				return
			}
			c.SecureJSON(http.StatusAccepted, gin.H{"message": "cart cleared successfully"})
			return
		}
		c.SecureJSON(http.StatusBadRequest, " user id param missing")
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
			existingCart.Count++
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
		c.SecureJSON(http.StatusAccepted, "added to cart successfully")
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	v1 := router.Group(routePrefix)
	{
		v1.GET(getCartDetails, getCartDetailsHandler(r))
		v1.DELETE(deleteFromCart, deleteFromCartHandler(r))
		v1.POST(checkout, checkoutHandler(r))
		v1.POST(clearCart, clearCartHandler(r))
		v1.POST(addToCart, addToCartHandler(r))
	}
	return v1
}
