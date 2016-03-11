package models

import (
	"bitbucket.org/newtopia2014/go-uuid/uuid"
	"database/sql/driver"
	"encoding/json"
	"genetics_loader/driver/database"
	"log"
	"time"
)

type UUID struct {
	uuid.UUID `json:"-"`
}

// Implement sql.Scanner interface for uuid.UUID
func (u *UUID) Scan(src interface{}) error {
	u.UUID = uuid.Parse(string(src.([]byte)))
	return nil
}

// Implement driver.Valuer interface for uuid.UUID
func (u UUID) Value() (driver.Value, error) {
	// if len(u.UUID) == 0 {
	// 	return []byte(""), nil
	// }
	return []byte(u.UUID.String()), nil
}

func (u UUID) String() string {
	return u.UUID.String()
}

func (u *UUID) Parse(uuidString string) {
	u.UUID = uuid.Parse(uuidString)
}

func (u *UUID) New() {
	uuidString := uuid.New()
	u.Parse(uuidString)
}

type UserSettings struct {
	ID          int      `json:"-"`
	UserId      UUID     `sql:"type:uuid" json:"-"`
	Settings    Metadata `sql:"type:jsonb" json:"settings"`
	Preferences Metadata `sql:"type:jsonb" json:"preferences"`
	Timestamps
	SoftDelete
}

type UserState struct {
	ID     int    `json:"id"`
	UserId UUID   `sql:"type:uuid" json:"-"`
	Type   string `sql:"size:100" json:"type"`
	State  string `sql:"size:100" json:"state"`
	Timestamps
	SoftDelete
}

type UserLog struct {
	ID      int      `json:"id"`
	UserId  UUID     `sql:"type:uuid" json:"-"`
	Name    string   `json:"name"`
	Meta    Metadata `sql:"type:jsonb" json:"meta"`
	Message string   `json="message"`
	Timestamps
	SoftDelete
}

type UserEmail struct {
	ID       int    `json:"id"`
	UserId   UUID   `sql:"type:uuid" json:"-"`
	Email    string `sql:"type:citext;size:100;" json:"email"`
	Verified bool   `json:"verified"`
	Timestamps
	SoftDelete
}

// TODO add Salutation?
type UserAddress struct {
	ID           int      `json:"id"`
	UserId       UUID     `sql:"type:uuid" json:"-"`
	Name         string   `json:"name"`
	AddressLine1 string   `json:"address_line_1"`
	AddressLine2 string   `json:"address_line_2"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	Country      string   `json:"country"`
	Zipcode      string   `sql:"size:6" json:"zipcode"`
	Plus4        string   `sql:"size:4" json:"plus_4"`
	IsBilling    bool     `json:"is_billing"`
	IsShipping   bool     `json:"is_shipping"`
	Meta         Metadata `json:"meta"`
	Timestamps
	SoftDelete
}

// TODO UserPhone or UserContact (Incl skype, etc)

type Session struct {
	ID         int  `json:"-"`
	Token      UUID `sql:"type:uuid;default:uuid_generate_v4();unique" json:"token"`
	UserId     UUID `sql:"type:uuid" json:"-"`
	Timestamps `json:"-"`
	SoftDelete `json:"-"`
}

type SessionDevice struct {
	SessionId  int
	OSType     string `gorm:"column:os_type"`
	DevicePush string
	Platform   string
	UserAgent  string
	Vendor     string
	Timestamps
	SoftDelete
}

type Permission struct {
	ID       int
	RoleId   int
	Type     string
	Action   string
	Resource string
	Timestamps
	SoftDelete
}

type Association struct {
	ID    int      `json:"-"`
	Type  string   `json:"type"`
	Users Metadata `sql:"type:jsonb" json:"users"`
	Timestamps
	SoftDelete
}

/**
 * Measure/Record Types:
 *   Clinical, Body, Well-being, Exercise
 */
type Type struct {
	ID   int    `json:"id"`
	Name string `sql:"size:100" json:"name"`
	Timestamps
	SoftDelete
}

type WellbeingRecordAggregate struct {
	Sleep    float64 `json:"sleep"`
	Mood     float64 `json:"mood"`
	Cravings float64 `json:"cravings"`
	Energy   float64 `json:"energy"`
	Anxiety  float64 `json:"anxiety"`
	Stress   float64 `json:"stress"`
}

type ExerciseRecordAggregate struct {
	Steps    float64 `json:"steps"`
	Time     float64 `json:"active_time"`
	Calories int64   `json:"calories_burned"`
}

type Measure struct {
	ID       int         `json:"id"`
	Name     string      `sql:"size:100" json:"name"`
	Type     Type        `json:"type"`
	TypeId   int         `json:"-"`
	Template MeasureInfo `sql:"type:jsonb" json:"template"`
	Timestamps
	SoftDelete
}

// Embedded in other structs to add CreatedAt, UpdatedAt
type Timestamps struct {
	CreatedAt time.Time `sql:"NOT NULL" json:"date_created"`
	UpdatedAt time.Time `sql:"NOT NULL" json:"date_updated"`
}

// Embedded in other structs to add DeletedAt
type SoftDelete struct {
	DeletedAt time.Time `sql:"default:NULL" json:"-"`
}

type Entity struct {
	ID          int       `json:"id"`
	Name        string    `sql:"size:100" json:"name"`
	Description string    `json:"description"`
	Type        Type      `json:"type"`
	TypeId      int       `json:"-"`
	Measures    []Measure `gorm:"many2many:entity_x_measure" json:"-"`
	Meta        Metadata  `sql:"type:jsonb" json:"meta"`
	Timestamps
	SoftDelete
}

func (e *Entity) CreateByName(name string) error {
	db := database.App
	var err error = nil
	var tempEntity Entity = Entity{}

	err = db.Where("name = ?", name).First(&tempEntity).Error
	if err != nil {
		log.Print("error getting entity from db : ", err.Error)
		return err
	}

	*e = tempEntity
	return nil
}

// MeasureInfo / MeasureData  -- For storing an array of measures
type MeasureInfo map[string]MeasureData

type MeasureData struct {
	Key            string    `json:"name"`
	Value          float64   `json:"value"`
	StringValue    string    `json:"string_value"`
	Unit           string    `json:"unit"`
	UnitShort      string    `json:"unit_short"`
	TimeStampValue time.Time `json:"timestamp_value"`
}

// Implement sql.Scanner interface for MeasureInfo
func (mi *MeasureInfo) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), mi)
}

// Implement driver.Valuer interface for MeasureInfo
func (mi MeasureInfo) Value() (driver.Value, error) {
	return json.Marshal(mi)
}

// Metadata -- Generic JSON
type Metadata map[string]interface{}

// Implement sql.Scanner interface for Metadata
func (mi *Metadata) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), mi)
}

// Implement driver.Valuer interface for Metadata
func (md Metadata) Value() (driver.Value, error) {
	return json.Marshal(md)
}
