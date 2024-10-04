package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
	"golang.org/x/crypto/bcrypt"
)

type validation struct {
	isValid bool
	message string
}

func ValidateEmail(r *util.Repository, email string) validation {
	if email == "" {
		return validation{false, "Email field cannot be empty"}
	}
	u := User{EmailAddress: email}
	UserModel := []models.User{}
	err := r.DB.Where(u).Find(&UserModel).Error
	fmt.Println(UserModel, len(UserModel), err)
	if err == nil {
		if len(UserModel) == 0 {
			return validation{true, "Email is new"}
		}
		return validation{false, "Email already Exists in DB"}
	}
	// fmt.Println(entries)
	return validation{false, "Error while validating Email" + err.Error()}
}

func UserSignUp(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var user1 User1
		var userdetails = User{}
		c.Bind(&userdetails)
		fmt.Println(userdetails)
		hashCost := 8
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userdetails.Password), hashCost)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Error while generating hash" + err.Error(),
			})
			return
		}
		userdetails.Password = string(passwordHash)

		if ValidateEmail(r, userdetails.EmailAddress).isValid {
			// fmt.Println(r.DB.Create(models.User{EmailAddress: userdetails.EmailAddress}).Error)
			err := r.DB.Where(User{EmailAddress: userdetails.EmailAddress}).FirstOrCreate(&userdetails).Error
			if err != nil {
				c.JSON(409, gin.H{
					"message": "Error While creating User",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "User Created succefully",
			})
			return
		}
		c.JSON(409, gin.H{
			"message": "User already exists in Database",
		})
	}
}

func validateUserCredentials(r *util.Repository, userdetails User) validation {
	email := userdetails.EmailAddress
	password := userdetails.Password
	emailValidator := ValidateEmail(r, email)
	if !emailValidator.isValid {
		u := User{EmailAddress: email}
		UserModel := []models.User{}
		err := r.DB.Where(u).Find(&UserModel).Error
		if err == nil {
			if len(UserModel) == 1 {
				if bcrypt.CompareHashAndPassword([]byte(UserModel[0].Password), []byte(password)) == nil {
					return validation{true, "Password verified successfully"}
				}
				return validation{false, "Invalid Password"}
			}
			return validation{false, "Multiple entries found with same Email"}
		}
		return validation{false, "Error while validating user" + err.Error()}
	}
	return validation{false, "Email is not a valid one"}

}

func UserLogin(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userdetails = User{}
		c.Bind(&userdetails)
		credentialValidator := validateUserCredentials(r, userdetails)
		if credentialValidator.isValid {
			accessToken, err := jwthandler.GenerateAccessToken(userdetails.EmailAddress)
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating accessToken" + err.Error(),
				})
			}
			refreshToken, err := jwthandler.GenerateRefreshToken(userdetails.EmailAddress)
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating refresh token" + err.Error(),
				})
			}
			c.SecureJSON(http.StatusOK, &map[string]string{
				"message":      "User Validated successfully",
				"accessToken":  accessToken,
				"refreshToken": refreshToken,
			})
			return

		}
		c.JSON(409, gin.H{
			"message": credentialValidator.message,
		})
	}
}
