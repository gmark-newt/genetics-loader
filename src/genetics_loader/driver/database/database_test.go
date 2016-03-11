package database

import (
	"testing"
)

func TestConnection(t *testing.T) {
	err := DB.DB().Ping()
	if err != nil {
		t.Fatal("Error Ping()ing database.")
	}
}
