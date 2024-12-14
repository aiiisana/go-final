package database

import (
	"ecommerce-platform/internal/models"
	"errors"

	"gorm.io/gorm"
)

func CreatePayment(db *gorm.DB, payment *models.Payment) error {
	if err := db.Create(payment).Error; err != nil {
		return err
	}
	return nil
}

func GetPayments(db *gorm.DB) ([]models.Payment, error) {
	var payments []models.Payment
	if err := db.Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func GetPaymentByID(db *gorm.DB, paymentID string) (*models.Payment, error) {
	var payment models.Payment
	if err := db.First(&payment, "id = ?", paymentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func DeletePayment(db *gorm.DB, paymentID string) error {
	if err := db.Delete(&models.Payment{}, "id = ?", paymentID).Error; err != nil {
		return err
	}
	return nil
}
