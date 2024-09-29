package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
		usersList := []models.User{}
		r.DB.Find(&usersList)
		m := map[string]interface{}{
			"message": "Details Fetched",
			"data":    usersList,
		}

		c.JSON(200, m)
	}

}

// todo
func AuthoriseUser(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramPairs := c.Request.URL.Query()
		for key, values := range paramPairs {
			fmt.Printf("key = %v, value(s) = %v\n", key, values)
		}
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
		// userGroup.GET(cart, getUserCartHandler(r))

	}
	return userGroup
}

/*

Uid 	Pid			COunt
karthik	mobile	2
karthik TV 		1


*/
