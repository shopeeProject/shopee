package admin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwthandler "github.com/shopeeProject/shopee/jwt"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/util"
	"golang.org/x/crypto/bcrypt"
)

type validation struct {
	isValid bool
	message string
}
type Admin struct {
	Name         string `json:"name"`
	PhoneNumber  string `json:"phoneNumber"`
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
}

func ValidateEmail(r *util.Repository, email string) validation {
	if email == "" {
		return validation{false, "Email field cannot be empty"}
	}
	u := Admin{EmailAddress: email}
	AdminModel := []models.Admin{}
	err := r.DB.Where(u).Find(&AdminModel).Error
	fmt.Println(AdminModel, len(AdminModel), err)
	if err == nil {
		if len(AdminModel) == 0 {
			return validation{true, "Email is new"}
		}
		return validation{false, "Email already Exists in DB"}
	}
	// fmt.Println(entries)
	return validation{false, "Error while validating Email" + err.Error()}
}

func AdminSignUp(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var user1 User1
		var admindetails = Admin{}
		c.Bind(&admindetails)
		fmt.Println(admindetails)
		hashCost := 8
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(admindetails.Password), hashCost)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Error while generating hash" + err.Error(),
			})
			return
		}
		admindetails.Password = string(passwordHash)

		if ValidateEmail(r, admindetails.EmailAddress).isValid {
			// fmt.Println(r.DB.Create(models.User{EmailAddress: admindetails.EmailAddress}).Error)
			err := r.DB.Where(Admin{EmailAddress: admindetails.EmailAddress}).FirstOrCreate(&admindetails).Error
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

func validateAdminCredentials(r *util.Repository, admindetails Admin) validation {
	email := admindetails.EmailAddress
	password := admindetails.Password
	emailValidator := ValidateEmail(r, email)
	if !emailValidator.isValid {
		u := Admin{EmailAddress: email}
		AdminModel := []models.Admin{}
		err := r.DB.Where(u).Find(&AdminModel).Error
		if err == nil {
			fmt.Println(AdminModel, "AdminModel print User login")
			if len(AdminModel) == 1 {
				if bcrypt.CompareHashAndPassword([]byte(AdminModel[0].Password), []byte(password)) == nil {
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

func AdminLogin(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var admindetails = Admin{}
		c.Bind(&admindetails)
		credentialValidator := validateAdminCredentials(r, admindetails)
		if credentialValidator.isValid {
			accessToken, err := jwthandler.GenerateAccessToken(admindetails.EmailAddress, "admin")
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating accessToken" + err.Error(),
				})
			}
			refreshToken, err := jwthandler.GenerateRefreshToken(admindetails.EmailAddress, "admin")
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating refresh token" + err.Error(),
				})
			}
			InsertResponse := jwthandler.InsertRefreshTokenToDB(r, refreshToken, admindetails.EmailAddress, "admin")
			if !InsertResponse.Success {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": InsertResponse.Message,
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
