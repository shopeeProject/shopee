package rating

import (
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/product"
	"github.com/shopeeProject/shopee/util"
)

type returnMessage struct {
	successful bool
	message    string
}

type Rating struct {
	UID         int    `json:"uid"`
	PID         int    `json:"pid"`
	Rating      string `json:"rating"`
	RatingValue int    `json:"ratingValue"`
	Description string `json:"description"`
}

func checkForExistingReview(r *util.Repository, UID int, PID int) returnMessage {
	condition := Rating{UID: UID, PID: PID}
	ExistingReviews := []models.Rating{}
	err := r.DB.Find(&ExistingReviews, condition).Error
	if err != nil {
		return returnMessage{
			successful: false,
			message:    "Error while checking for existing review",
		}
	}
	if len(ExistingReviews) == 0 {
		return returnMessage{
			successful: true,
			message:    "No existing records present",
		}
	}
	return returnMessage{
		successful: false,
		message:    "Existing review",
	}
}
func computeRating(r *util.Repository, newRating Rating) returnMessage {
	condition := Rating{
		PID: newRating.PID,
	}
	var count int64
	err := r.DB.Model(&models.Rating{}).Where(condition).Count(&count).Error
	if err != nil {
		return returnMessage{
			successful: false,
			message:    "Error while checking for past ratings" + err.Error(),
		}
	}
	var resultSum int64
	err = r.DB.Model(&models.Rating{}).Select("sum(rating)").Where(condition).Group("p_i_d").First(&resultSum).Error
	if err != nil {
		return returnMessage{
			successful: false,
			message:    "Error while checking for past ratings" + err.Error(),
		}
	}
	newRatingValue := float32(resultSum) / float32(count)
	ratingUpdateResponseFromProduct := product.UpdateRating(r, newRating.PID, newRatingValue)
	return returnMessage{successful: ratingUpdateResponseFromProduct.Successful, message: ratingUpdateResponseFromProduct.Message}

}

func Addrating(r *util.Repository, newRating Rating) returnMessage {

	existingReview := checkForExistingReview(r, newRating.UID, newRating.PID)
	if existingReview.successful {
		err := r.DB.Create(newRating).Error
		if err != nil {
			return returnMessage{
				successful: false,
				message:    "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, newRating)
		return newRatingComputationResponse
	}
	return returnMessage{
		successful: false,
		message:    "Review already exists on the product by the user",
	}

}

func Deleterating(r *util.Repository, oldRating Rating) returnMessage {

	existingReview := checkForExistingReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.successful {
		err := r.DB.Where(condition).Delete(&models.Rating{}).Error
		if err != nil {
			return returnMessage{
				successful: false,
				message:    "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, oldRating)
		return newRatingComputationResponse

	}
	return returnMessage{
		successful: false,
		message:    "Review doesnot exists on the product by the user",
	}

}

func ModifyRating(r *util.Repository, oldRating Rating) returnMessage {

	existingReview := checkForExistingReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.successful {
		err := r.DB.Where(condition).Updates(oldRating).Error
		if err != nil {
			return returnMessage{
				successful: false,
				message:    "Error while adding new rating" + err.Error(),
			}
		}
		newRatingComputationResponse := computeRating(r, oldRating)
		return newRatingComputationResponse

	}
	return returnMessage{
		successful: false,
		message:    "Review doesnot exists on the product by the user",
	}

}
