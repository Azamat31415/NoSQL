package migrations

import (
	"gorm.io/gorm"
)

// PickupPoint model for pickup points
type PickupPoint struct {
	ID           uint    `gorm:"primaryKey"`
	Name         string  `gorm:"type:varchar(255);not null"`
	Address      string  `gorm:"type:text;not null"`
	City         string  `gorm:"type:varchar(100);not null"`
	Latitude     float64 `gorm:"type:double precision;not null"`
	Longitude    float64 `gorm:"type:double precision;not null"`
	Phone        string  `gorm:"type:varchar(20)"`
	WorkingHours string  `gorm:"type:varchar(100)"`
}

// MigratePickupPoint for creation of pickup point table
func MigratePickupPoint(db *gorm.DB) error {
	if err := db.AutoMigrate(&PickupPoint{}, &Order{}); err != nil {
		return err
	}
	return nil
}
