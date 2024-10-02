package rating

import (
	"github.com/shopeeProject/shopee/models"
	util "github.com/shopeeProject/shopee/util"
)

type Rating struct {
	UID         int    `json:"uid"`
	PID         int    `json:"pid"`
	Rating      string `json:"rating"`
	RatingValue int    `json:"ratingValue"`
	Description string `json:"description"`
}

func isUserFirstReview(r *util.Repository, UID int, PID int) util.Response {
	condition := Rating{UID: UID, PID: PID}
	ExistingReviews := []models.Rating{}
	err := r.DB.Find(&ExistingReviews, condition).Error
	if err != nil {
		return util.Response{
			Message: "Error while checking for existing review",
		}
	}
	if len(ExistingReviews) == 0 {
		return util.Response{
			Success: true,
			Message: "No existing records present",
		}
	}
	return util.Response{
		Message: "Review Exists",
	}
}

// func UpdateRating(r *util.Repository, pID int, rating float32) util.Response {
// 	// Update the product rating in the database
// 	if err := r.DB.Model(&models.Product{}).Where("pid = ?", pID).Update("rating", rating).Error; err != nil {
// 		return util.Response{
// 			Message: err.Error(),
// 		}
// 	}
// 	return util.Response{Successful: true}
// }

func AddRating(r *util.Repository, newRating Rating) util.Response {

	existingReview := isUserFirstReview(r, newRating.UID, newRating.PID)
	if existingReview.Success {
		err := r.DB.Create(newRating).Error
		if err != nil {
			return util.Response{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		return util.Response{Success: true}
		// newRatingComputationResponse := computeRating(r, newRating)
		// return newRatingComputationResponse
	}
	return util.Response{
		Message: "Review already exists on the product by the user",
	}

}

func DeleteRating(r *util.Repository, oldRating Rating) util.Response {

	existingReview := isUserFirstReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.Success && existingReview.Message == "Review Exists" {
		err := r.DB.Where(condition).Delete(&models.Rating{}).Error
		if err != nil {
			return util.Response{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		return util.Response{Success: true}
		// newRatingComputationResponse := computeRating(r, oldRating)
		// return newRatingComputationResponse

	}
	return util.Response{
		Message: "Review does not exists on the product by the user",
	}

}

func ModifyRating(r *util.Repository, oldRating Rating) util.Response {

	existingReview := isUserFirstReview(r, oldRating.UID, oldRating.PID)
	condition := Rating{UID: oldRating.UID, PID: oldRating.PID}
	if !existingReview.Success && existingReview.Message == "Review Exists" {
		err := r.DB.Where(condition).Updates(oldRating).Error
		if err != nil {
			return util.Response{
				Message: "Error while adding new rating" + err.Error(),
			}
		}
		return util.Response{Success: true}
		// newRatingComputationResponse := computeRating(r, oldRating)
		// return newRatingComputationResponse

	}
	return util.Response{
		Message: "Review does not exists on the product by the user",
	}

}
