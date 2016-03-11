package models

import (
	"genetics_loader/driver/database"
	"testing"
)

func TestCreateAssociation(t *testing.T) {
	var users []User

	db := *(database.App)
	db.Find(&users)
	if len(users) < 2 {
		t.Fatal("Less than 2 users, can't associate.")
	}

	p := Participant{users[0]}
	c := Coach{users[1]}

	p.SetCoach(&c)

	var assoc Association
	err := db.Where("type = ? AND users #>> '{participant}' = ? AND users #>> '{coach}' = ?", "coach:participant", p.UserId.String(), c.UserId.String()).First(&assoc).Error
	if err != nil {
		t.Fatal(err)
	}

	if assoc.ID == 0 {
		t.Fatal("Association not found")
	}
}

func TestSetState(t *testing.T) {
	var user User
	var err error

	err = database.App.First(&user).Error
	if err != nil {
		t.Fatal("Couldn't get a user")
	}
	err = user.SetState("test", "myValue")
	if err != nil {
		t.Fatal("Couldn't set state", err)
	}
}

func TestRoles(t *testing.T) {
	var user User
	var roles []Role
	var err error

	name := "abcdefg"

	err = database.App.First(&user).Error
	if err != nil {
		t.Fatal("Couldn't get a user: ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	// Verify role not present
	if rolePresent(name, roles) != false {
		t.Fatal("Test role already present for user.")
	}

	// Add role
	err = user.AddRole(name)
	if err != nil {
		t.Fatal("Error adding role: ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	if rolePresent(name, roles) != true {
		t.Fatal("Test role not present after add.")
	}

	// Remove role
	err = user.RemoveRole(name)
	if err != nil {
		t.Fatal("Error removing role: ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	if rolePresent(name, roles) != false {
		t.Fatal("Test role present after remove.")
	}
}

func rolePresent(name string, roles []Role) bool {
	for _, r := range roles {
		if r.Name == name {
			return true
		}
	}

	return false
}

func TestRoleDuplication(t *testing.T) {
	var user User
	var roles []Role
	var err error

	name := "duplicate_role"

	err = database.App.First(&user).Error
	if err != nil {
		t.Fatal("Couldn't get a user: ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	startCount := len(roles)

	// Verify role not present
	if rolePresent(name, roles) != false {
		t.Fatal("Test role already present for user.")
	}

	// Add role
	err = user.AddRole(name)
	if err != nil {
		t.Fatal("Error adding role: ", err)
	}

	// Add role again
	err = user.AddRole(name)
	if err != nil {
		t.Fatal("Error adding role (dupe): ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	afterCount := len(roles)

	if afterCount > startCount+1 {
		t.Fatal("Role was duplicated.")
	}

	for _, r := range roles {
		t.Logf("Role: %v", r.Name)
	}

	// Cleanup, Remove role
	err = user.RemoveRole(name)
	if err != nil {
		t.Fatal("Error removing role: ", err)
	}

	err = database.App.Model(&user).Association("Roles").Find(&roles).Error
	if err != nil {
		t.Fatal("Error checking roles: ", err)
	}

	if rolePresent(name, roles) != false {
		t.Fatal("Test role present after remove.")
	}
}
