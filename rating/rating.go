package rating

import (
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/product"
	util "github.com/shopeeProject/shopee/util"
)

type Rating struct {
	UID         int    `json:"uid"`
	PID         int    `json:"pid"`
	Rating      string `json:"rating"`
	RatingValue int    `json:"ratingValue"`
	Description string `json:"description"`
}

func checkForExistingReview(r *util.Repository, UID int, PID int) util.ReturnMessage {
	condition := Rating{UID: UID, PID: PID}
	ExistingReviews := []models.Rating{}
	err := r.DB.Find(&ExistingReviews, condition).Error
	if err != nil {
		return util.ReturnMessage{
			Message: "Error while checking for existing review",
		}
	}
	if len(ExistingReviews) == 0 {
		return util.ReturnMessage{
			Successful: true,
			Message:    "No existing records present",
		}
	}
	return util.ReturnMessage{
		Message: "Existing review",
	}
}
func computeRating(r *util.Repository, newRating Rating) util.ReturnMessage {
	condition := Rating{
		PID: newRating.PID,
	}
	var count int64
	err := r.DB.Model(&models.Rating{}).Where(condition).Count(&count).Error
	if err != nil {
		return util.ReturnMessage{
			Message: "Error while checking for past ratings" + err.Error(),
		}
	}
	var resultSum int64
	err = r.DB.Model(&models.Rating{}).Select("sum(rating)").Where(condition).Group("p_i_d").First(&resultSum).Error
	if err != nil {
		return util.ReturnMessage{
			Message: "Error while checking for past ratings" + err.Error(),
		}
	}
	newRatingValue := float32(resultSum) / float32(count)
	ratingUpdateResponseFromProduct := product.UpdateRating(r, newRating.PID, newRatingValue)
	return util.ReturnMessage{Successful: ratingUpdateResponseFromProduct.Successful, Message: ratingUpdateResponseFromProduct.Message}

}

func Addrating(r *util.Repository, newRating Rating) util.ReturnMessage {

	existingReview := checkForExistingReview(r, newRating.UID, newRating.PID)
	if existingReview.Successful {
		err := r.DB.Create(newRating).Error
		if err != nil {
			return util.ReturnMessage{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, newRating)
		return newRatingComputationResponse
	}
	return util.ReturnMessage{
		Message: "Review already exists on the product by the user",
	}

}

func Deleterating(r *util.Repository, oldRating Rating) util.ReturnMessage {

	existingReview := checkForExistingReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.Successful {
		err := r.DB.Where(condition).Delete(&models.Rating{}).Error
		if err != nil {
			return util.ReturnMessage{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, oldRating)
		return newRatingComputationResponse

	}
	return util.ReturnMessage{
		Message: "Review doesnot exists on the product by the user",
	}

}

func ModifyRating(r *util.Repository, oldRating Rating) util.ReturnMessage {

	existingReview := checkForExistingReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.Successful {
		err := r.DB.Where(condition).Updates(oldRating).Error
		if err != nil {
			return util.ReturnMessage{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, oldRating)
		return newRatingComputationResponse

	}
	return util.ReturnMessage{
		Message: "Review doesnot exists on the product by the user",
	}

}
