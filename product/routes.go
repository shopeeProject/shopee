package product

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/order"
	util "github.com/shopeeProject/shopee/util"
)

const (
	getProductDetails = "/:id"
)

func getProductDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch user details from db and return
		pId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Invalid Product ID",
			})
			return
		}
		product := []models.Product{}
		condition := models.Product{PID: pId}

		err = r.DB.Find(&product, condition).Error
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Error while fetching Product",
			})
			return
		}
		if len(product) == 0 {
			c.JSON(409, gin.H{
				"message": "Error while fetching Product",
			})
			return
		}
		m := map[string]interface{}{
			"message": "Details Fetched",
			"data":    product[0],
		}
		c.JSON(200, m)
	}

}

type Product struct {
	PID          int
	Name         string `json:"name"`
	Price        string `json:"price"`
	Availability bool   `json:"availability"`
	Rating       string `json:"rating"`
	CategoryID   int    `json:"category"`
	Description  string `json:"description"`
	SID          string `json:"sid"`
	Image        string `json:"image"`
}

func updateAvailability(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch user details from db and return
		pId, err := strconv.Atoi(c.Param("pid"))
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Invalid Product ID",
			})
			return
		}
		// product := []models.Product{}
		condition := models.Product{PID: pId}
		availability := c.Param("count")

		err = r.DB.Model(&models.Product{}).Where(condition).Update("availability", availability).Error
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Error while Updating Product",
			})
			return
		}

		m := map[string]interface{}{
			"message": "Data Updated successfully",
		}
		c.JSON(200, m)
	}

}

func updateCount(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch user details from db and return
		pId, err := strconv.Atoi(c.Param("pid"))
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Invalid Product ID",
			})
			return
		}
		// product := []models.Product{}
		condition := models.Product{PID: pId}
		count, err := strconv.Atoi(c.Param("count"))
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Invalid Count Value",
			})
			return
		}
		err = r.DB.Model(&models.Product{}).Where(condition).Update("count", count).Error
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Error while Updating Product",
			})
			return
		}

		m := map[string]interface{}{
			"message": "Data Updated successfully",
		}
		c.JSON(200, m)
	}

}

func GetProductDetails(r *util.Repository, PIDs []int) ([]models.Product, error) {
	var productDetails []models.Product
	err := r.DB.Model(&Product{}).Where("pid IN ?", PIDs).Find(productDetails).Error
	if err != nil {
		return nil, err
	}
	return productDetails, err
}

func BuyNow(r *util.Repository, UID int, PID int) util.ReturnMessage {
	PIDs := make([]int, 1)
	PIDs = append(PIDs, PID)
	productsList, err := GetProductDetails(r, PIDs)
	if err != nil {
		return util.ReturnMessage{Message: "error fetching product ids"}
	}
	fmt.Print(productsList)
	order.PlaceOrderHandler1(r, UID, PIDs, productsList)
	return util.ReturnMessage{}
}

func UpdateRating(r *util.Repository, pId int, rating float32) util.ReturnMessage {

	condition := models.Product{PID: pId}
	err := r.DB.Model(&models.Product{}).Where(condition).Update("rating", rating).Error
	if err != nil {
		return util.ReturnMessage{
			Successful: false,
			Message:    "Error while updating rating",
		}

	}

	return util.ReturnMessage{
		Successful: true,
		Message:    "Rating updated successfully",
	}

}

func GroupProductRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	// router.GET(createUser, UserSignUp(r))
	// router.POST(userLogin, UserLogin(r))
	userGroup := router.Group("/product")
	userGroup.Use()
	{
		userGroup.GET(getProductDetails, getProductDetailsHandler(r))

	}
	return userGroup
}
