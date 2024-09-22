package user

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
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
)

type User struct {
	Uid           int
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

		c.JSON(200, gin.H{
			"message": "Details Updated successfully",
		})
	}
}

func getUserHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userdetails := User{
			Uid: 1,
		}
		userdetails.AccountStatus = "Down"
		userdetails.Address = "hyd"
		//fetch user details from db and return
		userdetailsJS, err := json.Marshal(userdetails)
		fmt.Println(string(userdetailsJS), err)
		if err == nil {
			c.JSON(200, gin.H{
				"message": "details fetched",
				"data":    userdetails,
			})
		}
	}

}

func AuthoriseUser(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		paramPairs := c.Request.URL.Query()
		for key, values := range paramPairs {
			fmt.Printf("key = %v, value(s) = %v\n", key, values)
		}
		c.Next()
	}
}
func GroupUserRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	router.POST(createUser, UserSignUp(r))
	router.POST(userLogin, UserLogin(r))
	userGroup := router.Group("/user")
	userGroup.Use(AuthoriseUser(r))
	{
		userGroup.POST(updateUserDetails, updateUserDetailsHandler(r))
		userGroup.POST(validateUser, updateUserDetailsHandler(r))
		userGroup.POST(deleteFromDb, updateUserDetailsHandler(r))
		userGroup.GET(getUser, getUserHandler(r))

	}
	return userGroup
}
