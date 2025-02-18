package migrations

import (
	"gorm.io/gorm"
	"time"
)

type Subscription struct {
	ID           uint `gorm:"primaryKey"`
	UserID       uint
	StartDate    time.Time
	RenewalDate  time.Time
	IntervalDays int
	Type         string `gorm:"type:varchar(20)"`
	Status       string `gorm:"type:varchar(20)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func MigrateSubscription(db *gorm.DB) error {
	return db.AutoMigrate(&Subscription{})
}

func ExpireSubscriptionsNow(db *gorm.DB) error {
	now := time.Now()
	result := db.Model(&Subscription{}).
		Where("renewal_date < ? AND status = ?", now, "active").
		Update("status", "expired")

	if result.Error != nil {
		return result.Error
	}
	return nil
}
