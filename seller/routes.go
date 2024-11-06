package seller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/category"
	firebaseOp "github.com/shopeeProject/shopee/firebase"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/product"
	util "github.com/shopeeProject/shopee/util"
)

// product and cart tables and handlers to be created
const (
	routePrefix      = "/seller"
	addProduct       = "/add-product"
	editProduct      = "/edit-product"
	updateProduct    = "/update-product"
	updateStatus     = "/update-seller-status"
	updateDetails    = "/update-details"
	createSeller     = "/create-seller"
	getSellerDetails = "/get-seller-details"
	getProducts      = "/get-products"
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

type Product struct {
	// PID          int
	Name         string  `json:"name"`
	Price        int     `json:"price"`
	Availability bool    `json:"availability"`
	Rating       float32 `json:"rating"`
	CategoryID   int     `json:"category"`
	Count        int     `json:"count"`
	Description  string  `json:"description"`
	SID          string  `json:"sid"`
	Image        string  `json:"image"`
}

func addProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		// err := c.ShouldBindJSON(&req)
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(309, map[string]interface{}{
				"message": "Error while parsing image " + err.Error(),
			})
			return
		}
		name := c.PostForm("name")
		price, err := strconv.Atoi(c.PostForm("price"))
		count, err := strconv.Atoi(c.PostForm("count"))
		description := c.PostForm("description")
		email, ok := c.Get("emailAddress")
		categoryId := c.PostForm("category")
		categoryIdInt, err := strconv.Atoi(categoryId)
		if err != nil {
			c.JSON(409, gin.H{
				"message": "Error while parsing form " + err.Error(),
			})
			return
		}
		if !ok {
			c.JSON(409, gin.H{
				"message": "User not Found/ Authenticated",
			})
			return
		}
		str, ok := email.(string)
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find Seller",
			})
			return
		}
		fmt.Println(file.Filename, name, price, count, description)

		validationResponse := category.ValidateCategory(r.DB, categoryIdInt)
		if validationResponse.Success == false {
			c.SecureJSON(http.StatusBadRequest, gin.H{
				"message": validationResponse.Message,
			})
			return
		}
		imageStringResp := firebaseOp.UploadImageAndGetUrl(file, name, str)
		fmt.Println(imageStringResp.Data)
		req := Product{Name: name, Price: price, Count: count, Image: imageStringResp.Data["downloadURL"], Description: description, SID: str, CategoryID: categoryIdInt}
		products := []models.Product{}
		result := r.DB.Where(Product{SID: str, Name: name}).Find(&products).Error
		if result != nil || len(products) > 0 {
			c.JSON(400, gin.H{
				"Message": "Error while fetching for validating inputs " + err.Error(),
			})
			return
		}

		if len(products) > 0 {
			c.JSON(400, gin.H{
				"Message": "Error while fetching for validating inputs : Multiple products found ",
			})
			return
		}

		// .FirstOrCreate(&req)
		err = r.DB.Create(req).Error
		//if valid category add product to product table
		if err != nil {
			log.Fatal("Failed to add product ", err.Error())
		}

		c.SecureJSON(http.StatusOK, gin.H{
			"message": "product added ",
		})
	}
}

func editProductHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {

		// err := c.ShouldBindJSON(&req)
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(309, map[string]interface{}{
				"message": "Error while parsing image " + err.Error(),
			})
			return
		}
		name := c.PostForm("name")
		price, err := strconv.Atoi(c.PostForm("price"))
		count, err := strconv.Atoi(c.PostForm("count"))
		description := c.PostForm("description")
		email, ok := c.Get("emailAddress")
		categoryId := c.PostForm("category")
		pid := c.PostForm("pid")
		categoryIdInt, err := strconv.Atoi(categoryId)
		availability := c.PostForm("availability")
		availabilityBool := false
		if availability == "true" {
			availabilityBool = true
		}
		if pid == "" {
			c.JSON(409, gin.H{
				"message": "Product Id is empty",
			})
			return
		}
		pIdInt, err := strconv.Atoi(pid)

		if err != nil {
			c.JSON(409, gin.H{
				"message": "Error while parsing form " + err.Error(),
			})
			return
		}
		if !ok {
			c.JSON(409, gin.H{
				"message": "User not Found/ Authenticated",
			})
			return
		}
		str, ok := email.(string)
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find Seller",
			})
			return
		}
		fmt.Println(file.Filename, name, price, count, description)

		validationResponse := category.ValidateCategory(r.DB, categoryIdInt)
		if validationResponse.Success == false {
			c.SecureJSON(http.StatusBadRequest, gin.H{
				"message": validationResponse.Message,
			})
			return
		}
		imageStringResp := firebaseOp.UploadImageAndGetUrl(file, file.Filename+name, str)
		fmt.Println(imageStringResp.Data)
		req := Product{Name: name, Price: price, Count: count, Image: imageStringResp.Data["downloadURL"], Description: description, SID: str, CategoryID: categoryIdInt, Availability: availabilityBool}
		err = r.DB.Where(models.Product{PID: pIdInt}).Updates(req).Error
		if err != nil {
			c.JSON(400, gin.H{
				"Message": "Error while fetching for validating inputs " + err.Error(),
			})
			return
		}

		// // .FirstOrCreate(&req)
		// err = r.DB.Create(req).Error
		// //if valid category add product to product table
		// if err != nil {
		// 	log.Fatal("Failed to Update product ", err.Error())
		// }

		c.SecureJSON(http.StatusOK, gin.H{
			"message": "product Edit successful ",
		})
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
		if !tokenValidationResponse.Success {
			returnString := map[string]interface{}{
				"message": tokenValidationResponse.Message,
			}
			c.SecureJSON(401, returnString)
			c.Abort()
			return
		}

		if tokenValidationResponse.Data["Entity"] != "seller" {
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

func AuthoriseSellerForOperations(r *util.Repository) gin.HandlerFunc {
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
		if !tokenValidationResponse.Success {
			returnString := map[string]interface{}{
				"message": tokenValidationResponse.Message,
			}
			c.SecureJSON(401, returnString)
			c.Abort()
			return
		}

		if tokenValidationResponse.Data["Entity"] != "seller" {
			returnString := map[string]interface{}{
				"message": tokenValidationResponse.Message,
			}
			c.SecureJSON(409, returnString)
			c.Abort()
			return
		}
		currentseller := []models.Seller{}
		res := r.DB.Find(&currentseller, Seller{EmailAddress: tokenValidationResponse.Data["Username"], IsApproved: true})
		if res.Error != nil {
			returnString := map[string]interface{}{
				"message": "Error while validating seller permissions " + res.Error.Error(),
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

func getProductsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		email, ok := c.Get("emailAddress")
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find User",
			})
			return
		}
		emailstr, ok := email.(string)
		if !ok {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find User",
			})
			return
		}
		data, err := product.GetAllSellerProductsF(r, emailstr)
		if err != nil {
			c.JSON(309, map[string]interface{}{
				"message": "Error while fetching products " + err.Error(),
			})
			return
		}
		c.JSON(200, map[string]interface{}{
			"message": "Data fetched successfully",
			"data":    data,
		})
		return
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	router.POST(createSeller, SellerSignUp(r))
	router.POST(sellerLogin, SellerLogin(r))
	v1 := router.Group(routePrefix)
	v1.Use(AuthoriseSeller(r))
	{
		v1.PATCH(updateDetails, updateDetailsHandler(r))
		v1.GET(getSellerDetails, getSellerDetailsHandler(r))
		v1.GET(getProducts, getProductsHandler(r))
	}
	v2 := router.Group((routePrefix))
	v2.Use(AuthoriseSellerForOperations(r))
	{
		v2.POST(addProduct, addProductHandler(r))
		v2.POST(editProduct, editProductHandler(r))
		v2.PATCH(updateProduct, updateProductHandler(r))
		v2.PATCH(updateStatus, updateStatusHandler(r))
	}
	return v1
}
