package database

import (
	"ecommerce-platform/internal/models"
	"time"

	"gorm.io/gorm"
)

func CreateAuditLog(db *gorm.DB, action string, userID uint) (*models.AuditLog, error) {
	logEntry := models.AuditLog{
		Action:    action,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	if err := db.Create(&logEntry).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

func GetUserSessions(db *gorm.DB, userID uint) ([]models.Session, error) {
	var sessions []models.Session
	if err := db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
func GetAuditLogs(db *gorm.DB, userID uint) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	if err := db.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func CreateSession(db *gorm.DB, userID uint, expiresAt time.Time) (*models.Session, error) {
	session := models.Session{
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := db.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteSession(db *gorm.DB, sessionID uint) error {
	if err := db.Delete(&models.Session{}, sessionID).Error; err != nil {
		return err
	}
	return nil
}
