package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

func AssignAdminRole(db *gorm.DB) error {
	var user User
	email := "azabraza061005@gmail.com"

	result := db.First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			adminUser := User{
				Email:     email,
				Password:  "aza061005",
				FirstName: "Admin",
				LastName:  "User",
				Role:      "admin",
			}
			if err := adminUser.HashPassword(); err != nil {
				return fmt.Errorf("error hashing password: %v", err)
			}

			if err := db.Create(&adminUser).Error; err != nil {
				return fmt.Errorf("error creating admin user: %v", err)
			}
			fmt.Println("Admin user created successfully.")
		} else {
			return fmt.Errorf("error checking user: %v", result.Error)
		}
	} else {
		user.Role = "admin"
		if err := db.Save(&user).Error; err != nil {
			return fmt.Errorf("error updating user role: %v", err)
		}
		fmt.Println("User role updated to admin.")
	}
	return nil
}
