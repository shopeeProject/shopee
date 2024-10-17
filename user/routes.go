package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

var user User
var users []User

const (
	updateUserDetails = "/update-user-details"
	validateUser      = "/validate-user"
	deleteFromDb      = "/delete-user-from-db"
	getUser           = "/get-user-details"
	createUser        = "/create-user"
	userLogin         = "/user-login"
	cart              = "/cart"
	orderlist         = "/orders"
	logout            = "/logout"
)

// func getorders (){
// 	orders.getorders(uid)
// }

type User struct {
	// UId           int
	Name          string `form:"name"`
	PhoneNumber   string `form:"phoneNumber"`
	EmailAddress  string `form:"emailAddress"`
	AccountStatus string `form:"accountStatus"`
	Address       string `form:"address"`
	Password      string `form:"password"`
}

type ExportUser struct {
	// UId           int
	Name          string `form:"name"`
	PhoneNumber   string `form:"phoneNumber"`
	EmailAddress  string `form:"emailAddress"`
	AccountStatus string `form:"accountStatus"`
	Address       string `form:"address"`
}

func updateUserDetailsHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userdetails User
		c.Bind(&userdetails)

		// fetch details of user from db and update details
		currentUser := []User{}
		condition := User{EmailAddress: userdetails.EmailAddress}
		r.DB.Limit(1).Find(&currentUser, condition)
		if len(currentUser) == 0 {
			c.JSON(409, gin.H{
				"message": "No Users found with given Email",
			})
			return
		}

		err := r.DB.Model(&models.User{}).Updates(userdetails).Error
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Error while updating the User",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "Details Updated successfully",
		})
	}
}

func getUserHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		//fetch user details from db and return
		// usersList := []models.User{}
		// r.DB.Find(&usersList)
		var userdetails = User{}
		c.Bind(&userdetails)
		// fmt.Println(userdetails, c.Get("emailAddress"))
		currentUser := []User{}

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
		condition := User{EmailAddress: str}
		fmt.Println(condition)
		r.DB.Limit(1).Find(&currentUser, condition)

		if len(currentUser) == 0 {
			c.JSON(309, map[string]interface{}{
				"message": "Couldn't find User",
			})
			return
		}
		// delete(currentUser[0],"Password")
		currentUser[0].Password = ""
		m := map[string]interface{}{
			"message": "Details Fetched",
			"data":    currentUser[0],
		}

		c.JSON(200, m)
	}

}

func getLogoutHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = User{}
		c.Bind(&user)
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
		fmt.Println(str)
		refreshTokenRemovalResponse := jwthandler.RemoveRefreshTokenFromDB(str)
		fmt.Println(refreshTokenRemovalResponse)
		if !refreshTokenRemovalResponse.Success {
			m := map[string]interface{}{
				"message": "Unable to logout " + refreshTokenRemovalResponse.Message,
			}
			c.JSON(409, m)
			c.Abort()
		}
		m := map[string]interface{}{
			"message": refreshTokenRemovalResponse.Message,
		}

		c.JSON(200, m)
	}

}

// todo
func AuthoriseUser(r *util.Repository) gin.HandlerFunc {
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
			c.SecureJSON(409, returnString)
			c.Abort()
			return
		}
		c.Set("emailAddress", tokenValidationResponse.Data["Username"])

		c.Next()
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	router.POST(createUser, UserSignUp(r))
	router.POST(userLogin, UserLogin(r))
	userGroup := router.Group("/user")
	userGroup.Use(AuthoriseUser(r))
	{
		userGroup.POST(updateUserDetails, updateUserDetailsHandler(r))
		userGroup.POST(validateUser, updateUserDetailsHandler(r))
		userGroup.POST(deleteFromDb, updateUserDetailsHandler(r))
		userGroup.GET(getUser, getUserHandler(r))
		userGroup.GET(logout, getLogoutHandler(r))
		// userGroup.GET(cart, getUserCartHandler(r))

	}
	return userGroup
}

/*

Uid 	Pid			COunt
karthik	mobile	2
karthik TV 		1


*/
