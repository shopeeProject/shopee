package seller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/category"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

// product and cart tables and handlers to be created
const (
	routePrefix      = "/seller"
	addProduct       = "/add-product"
	updateProduct    = "/update-product"
	updateStatus     = "/update-seller-status"
	updateDetails    = "/update-details"
	createSeller     = "/create-seller"
	getSellerDetails = "/get-seller-details"
	sellerLogin      = "/seller-login"
)

type Seller struct {
	// SID          int   `json:"sid"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
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
		// paramPairs := c.Request.URL.Query()

		accessToken, ok := c.Request.Header["Authorization"]
		if !ok {
			returnString := map[string]interface{}{
				"message": "Access token not found in the request",
			}
			c.SecureJSON(409, returnString)
			c.Abort()
			return
		}
		fmt.Println(accessToken)
		tokenValidationResponse := jwthandler.JwtMiddleware(accessToken[0])
		fmt.Println(tokenValidationResponse.Message)
		if !tokenValidationResponse.Success || tokenValidationResponse.Data["Entity"] != "seller" {
			returnString := map[string]interface{}{
				"message": tokenValidationResponse.Message,
			}
			c.SecureJSON(409, returnString)
			c.Abort()
			return
		}
		c.Set("emailAddress", tokenValidationResponse.Data["Username"])

		c.Next()
	}
}

func getSellerDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch user details from db and return
		// usersList := []models.User{}
		// r.DB.Find(&usersList)
		var sellerdetails = Seller{}
		c.Bind(&sellerdetails)
		// fmt.Println(userdetails, c.Get("emailAddress"))
		currentSeller := []Seller{}

		email, ok := c.Get("emailAddress")
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find User",
			})
			return
		}
		str, ok := email.(string)
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find User",
			})
			return
		}
		condition := Seller{EmailAddress: str}
		fmt.Println(condition)
		r.DB.Limit(1).Find(&currentSeller, condition)

		if len(currentSeller) == 0 {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find Seller",
			})
			return
		}
		// delete(currentSeller[0],"Password")
		currentSeller[0].Password = ""
		m := map[string]interface{}{
			"message": "Details Fetched",
			"data":    currentSeller[0],
		}

		c.JSON(200, m)
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
		v1.GET(getSellerDetails, getSellerDetailsHandler(r))
	}
	return v1
}
