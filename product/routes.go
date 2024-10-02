package product

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	order "github.com/shopeeProject/shopee/order"
	rating "github.com/shopeeProject/shopee/rating"
	util "github.com/shopeeProject/shopee/util"
)

const (
	routePrefix       = "/product"
	updateCount       = "/update-count"
	updateRating      = "/update-rating"
	addRating         = "/add-rating"
	deleteRating      = "/delete-rating"
	buyNow            = "/buy-now"
	getProductDetails = "/get-product"
)

func ComputeRating(r *util.Repository, newRating rating.Rating) util.Response {
	condition := rating.Rating{
		PID: newRating.PID,
	}
	var count int64
	err := r.DB.Model(&models.Rating{}).Where(condition).Count(&count).Error
	if err != nil {
		return util.Response{
			Message: "Error while checking for past ratings" + err.Error(),
		}
	}
	var resultSum int64
	err = r.DB.Model(&models.Rating{}).Select("sum(rating)").Where(condition).Group("p_i_d").First(&resultSum).Error
	if err != nil {
		return util.Response{
			Message: "Error while checking for past ratings" + err.Error(),
		}
	}
	newRatingValue := float32(resultSum) / float32(count)

	ratingUpdateResponseFromProduct := UpdateRatingForProduct(r, newRating.PID, newRatingValue)
	return ratingUpdateResponseFromProduct

}

func UpdateRatingForProduct(r *util.Repository, pid int, ratingValue float32) util.Response {

	// Update the product count in the database
	if err := r.DB.Model(&models.Product{}).Where("pid = ?", pid).Update("rating", ratingValue).Error; err != nil {
		log.Fatal("failed to update product rating:", err)
		return util.Response{Message: "failed to update product rating"}
	}

	return util.Response{Success: true}
}

// Update product count
func UpdateCountHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			PID   int `json:"pid" binding:"required"`
			Count int `json:"count" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
			return
		}

		// Update the product count in the database
		if err := r.DB.Model(&models.Product{}).Where("pid = ?", req.PID).Update("count", req.Count).Error; err != nil {
			log.Fatal("failed to update product count:", err)
			c.SecureJSON(http.StatusInternalServerError, gin.H{"message": "error updating product count"})
			return
		}

		c.SecureJSON(http.StatusOK, "product count updated")
	}
}

// Update product rating
func UpdateRatingHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rating.Rating

		if err := c.ShouldBindJSON(&req); err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
			return
		}
		errStruct := rating.ModifyRating(r, req)
		// // Update the product rating in the database
		if !errStruct.Success {
			c.SecureJSON(http.StatusInternalServerError, gin.H{"message": "error updating product rating"})
			return
		}

		c.SecureJSON(http.StatusOK, "product rating updated")
	}
}

func AddRatingHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req rating.Rating

		if err := c.ShouldBindJSON(&req); err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
			return
		}
		errStruct := rating.AddRating(r, req)
		// // Update the product rating in the database
		if !errStruct.Success {
			c.SecureJSON(http.StatusInternalServerError, gin.H{"message": "error updating product rating"})
			return
		}

		c.SecureJSON(http.StatusOK, "product rating updated")
	}
}

// Buy product (simplified)
func BuyNowHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			PID    int `json:"pid" binding:"required"`
			UserID int `json:"uid" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
			return
		}

		productIdArray := []int{req.PID}
		productsList, err := GetProductDetails(r, productIdArray)
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, err)
		}
		returnMessage := order.PlaceOrderHandler1(r, req.UserID, productIdArray, productsList)
		if returnMessage.Success == false {

			c.SecureJSON(http.StatusInternalServerError, "buynow failed")
		}
		//exp code above
		c.SecureJSON(http.StatusOK, gin.H{"message": "purchase successful", "pid": req.PID})
	}
}

func GetProductDetails(r *util.Repository, PIDs []int) ([]models.Product, error) {
	var productDetails []models.Product
	err := r.DB.Model(&models.Product{}).Where("pid IN ?", PIDs).Find(productDetails).Error
	if err != nil {
		return nil, err
	}
	return productDetails, err
}

// Get product details using query parameters
func GetProductDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		pid, err := strconv.Atoi(c.Query("pid")) // Get PID from query parameter
		if err != nil {
			c.SecureJSON(http.StatusBadRequest, gin.H{"message": "pid is required"})
			return
		}

		productIDArray := []int{pid}
		productList, err := GetProductDetails(r, productIDArray)

		// Fetch product from the database
		if err != nil {
			log.Fatal("product not found:", err)
			c.SecureJSON(http.StatusNotFound, gin.H{"message": "product not found"})
			return
		}

		c.SecureJSON(http.StatusOK, productList[0])
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	// todo validate if request is coming from valid seller/user
	v1 := router.Group(routePrefix)
	{
		v1.PATCH(updateCount, UpdateCountHandler(r))
		v1.PATCH(updateRating, UpdateRatingHandler(r))
		v1.POST(addRating, UpdateRatingHandler(r))
		v1.DELETE(deleteRating, UpdateRatingHandler(r))
		v1.POST(buyNow, BuyNowHandler(r))
		v1.GET(getProductDetails, GetProductDetailsHandler(r))
	}
	return v1
}
