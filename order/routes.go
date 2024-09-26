package order

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/payment"
	util "github.com/shopeeProject/shopee/util"
	"gorm.io/gorm"
)

type returnMessage struct {
	Successful bool
	Message    string
}

const (
	routePrefix = "/order"
	placeOrder  = "/place-order"
	trackOrder  = "/track-order/:order_id"
	cancelOrder = "/cancel-order/:order_id"
	updateOrder = "/update-order"
	addStage    = "/add-stage"
)

type Order struct {
	UID           int      `json:"uid"`
	Products      []int    `json:"products"`
	PaymentID     int      `json:"payment_id"`
	PaymentStatus string   `json:"payment_status"`
	Address       string   `json:"address"`
	Stages        []string `json:"stages"`
	Price         int      `json:"price"`
	OrderStatus   string   `json:"order_status"`
}

// Place Order Handler
func PlaceOrderHandler(r *util.Repository) gin.HandlerFunc {
	//todo: buynow and checkout should go through this
	return func(c *gin.Context) {
		var order Order
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := r.DB.Create(&order).Error; err != nil {
			log.Fatal("failed to place order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to place order"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "order placed successfully"})
	}
}

// Track Order Handler
func trackOrderHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		var order models.Order
		if err := r.DB.First(&order, orderID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		c.JSON(http.StatusOK, order)
	}
}

// Cancel Order Handler
func cancelOrderHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		//todo update product status when order cancelled
		if err := r.DB.Delete(&models.Order{}, orderID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to cancel order"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order canceled successfully"})
	}
}

// Update Order Handler
func updateOrderHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		order := models.Order{}
		if err := c.ShouldBindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := r.DB.Model(&models.Order{}).Where("oid = ?", order.OID).Updates(order).Error; err != nil {
			log.Fatal("failed to update order:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "order updated successfully"})
	}
}

// Add Stage Handler
func addStageHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stage struct {
			OrderID uint   `json:"order_id"`
			Stage   string `json:"stage"`
		}
		if err := c.ShouldBindJSON(&stage); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Assuming stages are stored in a related table
		if err := r.DB.Model(&models.Order{}).Where("oid = ?", stage.OrderID).Update("stages", gorm.Expr("array_append(stages, ?)", stage.Stage)).Error; err != nil {
			log.Fatal("failed to add stage:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add stage"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "stage added successfully"})
	}
}

type returnMessageWithData struct {
	Successful bool
	Message    string
	Data       models.Payment
}

func PlaceOrderHandler1(r *util.Repository, UID int, PIDs []int, productsList map[string]interface{}) returnMessage {
	order := Order{
		UID:      UID,
		Products: PIDs,
	}

	if err := r.DB.Create(&order).Error; err != nil {
	}

	var Amount int
	var orderid int
	returnmessage := payment.MakePayment(r, UID, Amount)
	updateFields := models.Payment{}
	if returnmessage.Successful {
		err := r.DB.Where(Order{UID: orderid}).Updates(updateFields).Error
		fmt.Print(err.Error())
		return returnMessage{}

	}

	return returnMessage{}

}

func RegisterRoutes(router *gin.Engine, r *util.Repository) {
	v1 := router.Group(routePrefix)
	{
		v1.POST(placeOrder, PlaceOrderHandler(r))
		v1.GET(trackOrder, trackOrderHandler(r))
		v1.DELETE(cancelOrder, cancelOrderHandler(r))
		v1.PUT(updateOrder, updateOrderHandler(r))
		v1.POST(addStage, addStageHandler(r))
	}
}
