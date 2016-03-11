package models

import (
	"errors"
	"fmt"
	"genetics_loader/driver/database"
	"github.com/dabfleming/gorm"
	"newtopia/log"
	"strings"
)

type User struct {
	ID             int       `json:"-"`
	UserId         UUID      `sql:"type:uuid;default:uuid_generate_v4();unique" json:"-"`
	PrimaryEmail   UserEmail `json:"primary_email"`
	PrimaryEmailId int       `sql:"default:NULL" json:"-"`
	DisplayName    string    `sql:"size:100" json:"display_name"`
	Password       []byte    `json:"-"`
	Roles          []Role    `gorm:"many2many:user_x_role" json:"-"`
	Meta           Metadata  `sql:"type:jsonb" json:"meta"`
	Timestamps
	SoftDelete
}

type Participant struct {
	User
}

type Coach struct {
	User
}

// New type to map target user from admin user
// route ( /admin/user/:userid/...) to handler
// functions... martini context requires distinct
// types and the session user is already mapped
type TargetUser struct {
	User
}

// tells gorm to use 'users' table
func (t TargetUser) TableName() string {
	return "users"
}

func (u *User) IsCoach() bool {
	db := database.App
	var isCoach bool = false

	var roles []Role = make([]Role, 0)
	err := db.Model(u).Association("Roles").Find(&roles).Error
	if err != nil {
		log.Println("Error finding roles for user: ", err.Error())
		return false
	}

	for _, role := range roles {
		if role.Name == "coach" {
			isCoach = true
		}
	}

	return isCoach
}

func (u *User) IsCareSpecialist() bool {
	db := database.App
	var isCm bool = false

	var roles []Role = make([]Role, 0)
	err := db.Model(u).Association("Roles").Find(&roles).Error
	if err != nil {
		log.Println("Error finding roles for user: ", err.Error())
		return false
	}

	for _, role := range roles {
		if role.Name == "cm" {
			isCm = true
		}
	}

	return isCm
}

func (u *User) IsParticipant() bool {
	db := database.App
	var isParticipant bool = false

	var roles []Role = make([]Role, 0)
	err := db.Model(u).Association("Roles").Find(&roles).Error
	if err != nil {
		log.Println("Error finding roles for user: ", err.Error())
		return false
	}

	for _, role := range roles {
		if role.Name == "participant" {
			isParticipant = true
		}
	}

	return isParticipant
}

func (u *User) IsLead() bool {
	db := database.App
	var isLead bool = false

	var roles []Role = make([]Role, 0)
	err := db.Model(u).Association("Roles").Find(&roles).Error
	if err != nil {
		log.Println("Error finding roles for user: ", err.Error())
		return false
	}

	for _, role := range roles {
		if role.Name == "lead" {
			isLead = true
		}
	}

	return isLead
}

func (u *User) AddLog(name, message string, meta Metadata) error {
	var tx *gorm.DB
	tx = database.App.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := u.AddLogWithTx(name, message, meta, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) AddLogWithTx(name, message string, meta Metadata, tx *gorm.DB) error {
	userLog := UserLog{
		UserId:  u.UserId,
		Name:    name,
		Message: message,
		Meta:    meta,
	}
	err := tx.Create(&userLog).Error
	if err != nil {
		log.Print("Error creating UserLog.", err)
		return err
	}
	return nil
}

func (u *User) RemoveRole(name string) error {
	var db *gorm.DB
	var role Role
	var err error

	// TODO version with/without TX
	db = database.App

	err = db.Where("name = ?", name).First(&role).Error
	if err == gorm.RecordNotFound {
		return nil
	} else if err != nil {
		return err
	}

	// Remove role
	err = db.Model(&u).Association("Roles").Delete(role).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) AddRoleWithTx(name string, db *gorm.DB) error {
	var role Role
	var err error

	// TODO version with/without TX

	err = db.FirstOrCreate(&role, Role{Name: name}).Error
	if err != nil {
		return err
	}

	// Associate role
	err = db.Model(&u).Association("Roles").Append(role).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) AddRole(name string) error {
	var db *gorm.DB
	var role Role
	var err error

	// TODO version with/without TX
	db = database.App

	err = db.FirstOrCreate(&role, Role{Name: name}).Error
	if err != nil {
		return err
	}

	// Associate role
	err = db.Model(&u).Association("Roles").Append(role).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) SetState(typeStr, state string) error {
	var tx *gorm.DB
	tx = database.App.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	err := u.SetStateWithTx(typeStr, state, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *User) SetStateWithTx(typeStr, state string, tx *gorm.DB) error {
	var userState UserState
	var err error

	// First soft-delete previous state of same type
	err = tx.Where("user_id = ? AND type = ?", u.UserId, typeStr).Delete(UserState{}).Error
	if err != nil {
		log.Print("db error clearing previous user state", u, typeStr)
		return err
	}

	// Set new state
	userState = UserState{
		UserId: u.UserId,
		Type:   typeStr,
		State:  state,
	}

	err = tx.Create(&userState).Error
	if err != nil {
		log.Print("db error creating new user state", u, typeStr)
		return err
	}

	return nil
}

func (u *User) SetAssociation(userKey, targetKey string, target *User) error {
	return u.associate(userKey, targetKey, target)
}

func (p *Participant) SetCoach(coach *Coach) error {
	var assoc Association
	err := database.App.Where("type = ? AND users #>> '{participant}' = ?", "coach:participant", p.UserId.String()).First(&assoc).Error
	if err != nil && err != gorm.RecordNotFound {
		return err
	}
	if assoc.ID != 0 {
		// Association exists, update
		assoc.Users["coach"] = coach.UserId.String()
		err = database.App.Save(&assoc).Error
		if err != nil {
			return err
		}
		return nil
	}
	return coach.associate("coach", "participant", &(p.User))
}

func (u *User) associate(userKey, targetKey string, target *User) error {
	var assoc Association

	if u.ID == target.ID {
		return errors.New("Can't associate user to self.")
	}

	userKeyLower := strings.ToLower(userKey)
	targetKeyLower := strings.ToLower(targetKey)

	var assocUuid UUID = UUID{}
	assocUuid.New()
	assocId := "assoc-" + assocUuid.String()

	assoc.Type = fmt.Sprintf("%s:%s", userKeyLower, targetKeyLower)
	assoc.Users = Metadata{
		userKeyLower:     u.UserId.String(),
		targetKeyLower:   target.UserId.String(),
		"association_id": assocId,
	}
	err := database.App.Create(&assoc).Error
	if err != nil {
		return err
	}

	return nil
}

func (user *User) HasRegistered() (bool, Record) {
	db := database.App
	var entityRegistration Entity
	var RegistrationRecord Record
	db.Where("name = ? and CAST (meta #>> '{active}' as BOOLEAN)", "Registration").First(&entityRegistration)
	RegistrationComplete := false

	if !db.Where("entity_id = ? AND user_id = ?", entityRegistration.ID, user.UserId).Preload("Type").Preload("Entity").Preload("Measure").First(&RegistrationRecord).RecordNotFound() {

		RegistrationComplete = RegistrationRecord.MeasureData["completed"].Value == 1

	}
	return RegistrationComplete, RegistrationRecord

}

func (user *User) HasAssessment() (bool, Record) {

	db := database.App

	var entityAssessment Entity
	var assessmentRecord Record
	db.Where("name = ? and CAST (meta #>> '{active}' as BOOLEAN)", "Assessment").First(&entityAssessment)

	assessmentComplete := false
	if !db.Where("entity_id = ? AND user_id = ?", entityAssessment.ID, user.UserId).First(&assessmentRecord).RecordNotFound() {

		assessmentComplete = assessmentRecord.MeasureData["done"].Value == 1

	}
	return assessmentComplete, assessmentRecord

}
