package admin

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/category"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/seller"
	"github.com/shopeeProject/shopee/util"
)

const (
	routePrefix          = "/admin"
	addCategory          = "/add-category"
	deleteCategory       = "/delete-category"
	approveSeller        = "/approve-seller"
	getSellers           = "/get-sellers"
	getUnapprovedSellers = "/get-unapproved-sellers"
	createAdmin          = "/create-admin"
	adminLogin           = "/admin-login"
)

// todo
func AuthoriseAdmin(r *util.Repository) gin.HandlerFunc {
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
		if !tokenValidationResponse.Success || tokenValidationResponse.Data["Entity"] != "admin" {
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

func GetUnApprovedSellers(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		SellersList := seller.GetUnApprovedSellers(r)
		if !SellersList.Success {
			c.JSON(409, gin.H{"message": SellersList.Message})
			return
		}
		c.JSON(200, SellersList)
	}
}

func GetSellers(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		SellersList := seller.GetSellers(r)
		if !SellersList.Success {
			c.JSON(409, gin.H{"message": SellersList.Message})
			return
		}
		c.JSON(200, SellersList)
	}
}

func ApproveSeller(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		sellerId := c.Request.Form["sId"]
		i, err := strconv.Atoi(sellerId[0])
		if sellerId == nil || err != nil {
			c.JSON(409, gin.H{"message": "Please provide Seller Id to approve " + err.Error()})
			return
		}
		SellersList := seller.ApproveSeller(r, i)
		if !SellersList.Success {
			c.JSON(409, gin.H{"message": SellersList.Message})
			return
		}
		c.JSON(200, SellersList)
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {

	router.POST(adminLogin, AdminLogin(r))
	v1 := router.Group(routePrefix)
	v1.Use(AuthoriseAdmin(r))
	{
		router.POST(createAdmin, AdminSignUp(r))
		router.GET(getUnapprovedSellers, GetUnApprovedSellers(r))
		router.GET(getSellers, GetSellers(r))
		router.POST(addCategory, category.AddCategoryHandler(r))
		router.DELETE(deleteCategory, category.RemoveCategoryHandler(r))
	}
	return v1
}
