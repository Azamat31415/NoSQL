package migrations

import (
	"gorm.io/gorm"
	"time"
)

type SubscriptionPayment struct {
	ID             uint `gorm:"primaryKey"`
	SubscriptionID uint
	Amount         float64
	PaymentDate    time.Time
	Status         string `gorm:"type:varchar(20)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func MigrateSubscriptionPayment(db *gorm.DB) error {
	return db.AutoMigrate(&SubscriptionPayment{})
}
