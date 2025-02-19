package migrations

import "go.mongodb.org/mongo-driver/mongo"

func MigrateAll(db *mongo.Database) error {
	if err := MigrateUser(db); err != nil {
		return err
	}
	if err := MigrateProduct(db); err != nil {
		return err
	}
	if err := MigrateOrder(db); err != nil {
		return err
	}
	if err := MigratePersonalPet(db); err != nil {
		return err
	}
	if err := MigratePickupPoint(db); err != nil {
		return err
	}
	if err := MigrateSubscription(db); err != nil {
		return err
	}
	if err := MigrateSubscriptionPayment(db); err != nil {
		return err
	}
	if err := MigrateCart(db); err != nil {
		return err
	}

	return nil
}
