package seller

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

type validation struct {
	isValid bool
	message string
}

func ValidateEmail(r *util.Repository, email string) validation {
	if email == "" {
		return validation{false, "Email field cannot be empty"}
	}
	seller := models.Seller{EmailAddress: email}
	sellerModel := []models.Seller{}
	err := r.DB.Where(seller).Find(&sellerModel).Error
	fmt.Println(sellerModel, len(sellerModel), err)
	if err == nil {
		if len(sellerModel) == 0 {
			return validation{true, "Email is new"}
		}
		return validation{false, "Email already Exists in DB"}
	}
	return validation{false, "Error while validating Email" + err.Error()}
}

func SellerSignUp(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sellerDetails = models.Seller{}
		c.ShouldBindJSON(&sellerDetails)

		hashCost := 8
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(sellerDetails.Password), hashCost)
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Error while generating hash" + err.Error(),
			})
			return
		}

		sellerDetails.Password = string(passwordHash)

		if ValidateEmail(r, sellerDetails.EmailAddress).isValid {

			r.DB.Where(models.Seller{EmailAddress: sellerDetails.EmailAddress}).FirstOrCreate(&sellerDetails)
			c.SecureJSON(http.StatusOK, &map[string]string{
				"message": "Seller Created successfully",
			})
			return
		}
		c.SecureJSON(http.StatusConflict, &map[string]string{
			"message": "Seller Creation failed",
		})
	}
}

func validateSellerCredentials(r *util.Repository, sellerdetails models.Seller) validation {
	email := sellerdetails.EmailAddress
	password := sellerdetails.Password
	emailValidator := ValidateEmail(r, email)
	if !emailValidator.isValid {
		u := models.Seller{EmailAddress: email}
		sellerModel := []models.Seller{}
		err := r.DB.Where(u).Find(&sellerModel).Error
		if err == nil {
			if len(sellerModel) == 1 {
				if bcrypt.CompareHashAndPassword([]byte(sellerModel[0].Password), []byte(password)) == nil {
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

func SellerLogin(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sellerdetails = models.Seller{}
		c.Bind(&sellerdetails)
		credentialValidator := validateSellerCredentials(r, sellerdetails)
		if credentialValidator.isValid {
			c.SecureJSON(http.StatusOK, &map[string]string{
				"message": "Seller Validated successfully",
			})
			return
		}
		c.JSON(409, gin.H{
			"message": credentialValidator.message,
		})
	}
}