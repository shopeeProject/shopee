package payment

import (
	"time"

	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/util"
)

type Payment struct {
	UID           int    `json:"uid"`
	Amount        int    `json:"amount"`
	Timestamp     string `json:"timestamp"`
	PaymentStatus string `json:"paymentStatus"`
	Description   string `json:"description"`
}
type returnMessageWithData struct {
	Successful bool
	Message    string
	Data       models.Payment
}

func getPaymentDetails(r *util.Repository, pId int) returnMessageWithData {
	condition := models.Payment{PayID: pId}
	PaymentDetails := models.Payment{}
	err := r.DB.Find(&PaymentDetails, condition).Error
	r.DB.Where(&models.Order{UID: 12})
	if err != nil {
		return returnMessageWithData{
			Successful: false,
			Message:    "Error while checking for Payment Deatils",
		}
	}
	return returnMessageWithData{
		Successful: false,
		Message:    "Details Fetched successfully",
		Data:       PaymentDetails,
	}
}

func MakePayment(r *util.Repository, UID int, amount int) returnMessageWithData {
	mumbai, _ := time.LoadLocation("Asia/Kolkata")

	currentTime := time.Now()
	mumbaiTime := currentTime.In(mumbai)
	newPayment := Payment{
		UID:           UID,
		Amount:        amount,
		Timestamp:     string(mumbaiTime.Unix()),
		PaymentStatus: "Completed",
		Description:   "",
	}

	err := r.DB.Create(&newPayment).Error
	if err != nil {
		return returnMessageWithData{
			Successful: false,
			Message:    "Error while updating new Payment ID",
		}
	}

	GeneratedPayment := models.Payment{}

	err = r.DB.Find(GeneratedPayment, newPayment).Error
	if err != nil {
		return returnMessageWithData{
			Successful: false,
			Message:    "Error while fetching new Payment ",
		}
	}
	return returnMessageWithData{
		Successful: true,
		Message:    "New Payment details updated successfully",
		Data:       GeneratedPayment,
	}

}
