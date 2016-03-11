package models

import ()

type UserAsset struct {
	ID     int
	UserId UUID `sql:"type:uuid" json:"-"`
	Type   string
	Data   []byte
	Timestamps
	SoftDelete
}
