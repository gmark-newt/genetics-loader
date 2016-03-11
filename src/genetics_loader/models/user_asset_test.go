package models

import (
	"genetics_loader/driver/database"
	"github.com/dabfleming/gorm"
	"testing"
)

var data = "This is a test of the user asset table."

func TestUserAssetTable(t *testing.T) {
	var user User
	var asset, lookup UserAsset
	var err error

	err = database.App.First(&user).Error
	if err != nil {
		t.Fatal("Could not find user.")
	}

	asset = UserAsset{
		UserId: user.UserId,
		Type:   "test",
		Data:   []byte(data),
	}

	err = database.App.Save(&asset).Error
	if err != nil {
		t.Fatalf("Error saving asset: %v", err)
	}

	err = database.App.First(&lookup, asset.ID).Error
	if err != nil {
		t.Fatalf("Error finding asset: ", err)
	}

	if lookup.UserId.String() != user.UserId.String() || lookup.Type != "test" || string(lookup.Data) != data {
		t.Fatal("Retrieved asset does not match saved asset.")
	}

	err = database.App.Unscoped().Delete(&lookup).Error
	if err != nil {
		t.Fatalf("Error (hard) deleting test asset: %v", err)
	}

	err = database.App.First(&lookup, asset.ID).Error
	if err != gorm.RecordNotFound {
		t.Fatalf("Unexpected error (expecting gorm.RecordNotFound) looking up (hard) deleted record: %v", err)
	}
}
