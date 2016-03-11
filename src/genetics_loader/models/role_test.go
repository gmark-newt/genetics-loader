package models

import (
	"genetics_loader/driver/database"
	"testing"
)

func TestInsertRole(t *testing.T) {
	var role Role
	err := database.App.FirstOrCreate(&role, Role{Name: "TestRole"}).Error
	if err != nil {
		t.Errorf("Database error: %v", err)
	}
}

func TestFindRole(t *testing.T) {
	var roles []Role
	err := database.App.Find(&roles, 1).Error
	if err != nil {
		t.Errorf("Database error: %v", err)
	}
	t.Logf("Roles %#v", roles)
}
