package models

import ()

type Role struct {
	ID              int
	Name            string
	InheritedRoleId int `sql:"default:NULL"`
	Timestamps
	SoftDelete
}
