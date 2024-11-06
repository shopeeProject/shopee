package seller

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
	seller := models.Seller{EmailAddress: email}
	sellerModel := []models.Seller{}
	err := r.DB.Find(&sellerModel, seller).Error
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
		fmt.Println(sellerDetails)
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
			} else {
				if len(sellerModel) > 1 {
					return validation{false, "Multiple entries found with same Email"}
				}
				return validation{false, "Seller Not found"}
			}
		}
		return validation{false, "Error while validating user" + err.Error()}
	}
	return validation{false, "Email is not a valid one"}

}

type SellersListResponse struct {
	util.Response
	data []Seller
}

func NullifyPassowrd(sellers []Seller) []Seller {
	for i := 0; i < len(sellers); i++ {
		sellers[i].Password = ""
	}
	return sellers
}

// todo
func GetUnApprovedSellers(r *util.Repository) SellersListResponse {
	sellerdetails := Seller{IsApproved: false}
	sellerResponse := []Seller{}
	err := r.DB.Where(sellerdetails).Find(&sellerResponse).Error
	sellerResponse = NullifyPassowrd(sellerResponse)
	if err != nil {
		return SellersListResponse{
			Response: util.Response{Message: "Error while fetching details " + err.Error(), Success: false},
			data:     sellerResponse,
		}
	}
	return SellersListResponse{
		Response: util.Response{
			Message: "Data fetched Successfully",
			Success: true,
		},
		data: sellerResponse,
	}
}

func GetSellers(r *util.Repository) SellersListResponse {
	sellerResponse := []Seller{}
	err := r.DB.Find(&sellerResponse).Error
	sellerResponse = NullifyPassowrd(sellerResponse)
	if err != nil {
		return SellersListResponse{
			Response: util.Response{Message: "Error while fetching details " + err.Error(), Success: false},
			data:     sellerResponse,
		}
	}
	return SellersListResponse{
		Response: util.Response{
			Message: "Data fetched Successfully",
			Success: true,
		},
		data: sellerResponse,
	}
}

// todo
func ApproveSeller(r *util.Repository, sellerId int) util.Response {
	sellerdetails := models.Seller{SID: sellerId}
	err := r.DB.Where(sellerdetails).Update("isApproved", true).Error
	if err != nil {
		return util.Response{
			Success: false,
			Message: "Error while Approving Seller " + err.Error(),
		}
	}
	return util.Response{
		Success: true,
		Message: "Seller approved Successfully",
	}
}

func SellerLogin(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sellerdetails = models.Seller{}
		c.Bind(&sellerdetails)
		credentialValidator := validateSellerCredentials(r, sellerdetails)
		if credentialValidator.isValid {
			accessToken, err := jwthandler.GenerateAccessToken(sellerdetails.EmailAddress, "seller")
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating accessToken" + err.Error(),
				})
			}
			refreshToken, err := jwthandler.GenerateRefreshToken(sellerdetails.EmailAddress, "seller")
			if err != nil {
				c.SecureJSON(http.StatusInternalServerError, &map[string]string{
					"message": "Error while generating refresh token" + err.Error(),
				})
			}
			InsertResponse := jwthandler.InsertRefreshTokenToDB(r, refreshToken, sellerdetails.EmailAddress, "user")
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
