package models

import (
	"encoding/json"
	"genetics_loader/driver/database"
	"testing"
)

func TestRecordHelperMethods(t *testing.T) {
	var entity Entity
	var user User
	var r Record

	// Fetch the Sleep entity.
	err := database.App.Where("name = ?", "Sleep").First(&entity).Error
	if err != nil {
		t.Fatal("Can't load 'Sleep' entity.")
	}
	// Fetch user (foo)
	err = database.App.First(&user).Error
	if err != nil {
		t.Fatal("Can't load a test user.")
	}

	// FromEntity
	err = r.FromEntity(&entity)
	if err != nil {
		t.Errorf("FromEntity returned error: %v", err)
	}
	t.Logf("After FromEntity, Meta: %v, MeasureData: %v", r.Meta, r.MeasureData)

	// SetMeasureValue, success
	err = r.SetMeasureValue("sleep", 8)
	if err != nil {
		t.Errorf("SetMeasureValue returned error: %v", err)
	}
	t.Logf("After SetMeasureValue, MeasureData: %v", r.MeasureData)

	// GetMeasureValue, success
	value, err := r.GetMeasureValue("sleep")
	if err != nil {
		t.Errorf("GetMeasureValue returned error: %v", err)
	}
	if value != 8 {
		t.Errorf("GetMeasureValue value does not match expected: %v", value)
	}

	// SetMeasureValue, error
	err = r.SetMeasureValue("foo", 9)
	if err == nil {
		t.Errorf("SetMeasureValue did not return expected error.")
	}
	t.Logf("After SetMeasureValue (no change expected), MeasureData: %v", r.MeasureData)

	// GetMeasureValue, error
	_, err = r.GetMeasureValue("foo")
	if err == nil {
		t.Errorf("GetMeasureValue did not return expected error.")
	}

	// SetMeta, one
	err = r.SetMeta(&Metadata{"one": 1})
	if err != nil {
		t.Errorf("SetMeta returned error: %v", err)
	}
	if _, ok := r.Meta["one"]; !ok {
		t.Errorf("Expected Meta key does not exist.")
	}
	t.Logf("After SetMeta, Meta: %v", r.Meta)

	// SetMeta, two
	err = r.SetMeta(&Metadata{"two": 2})
	if err != nil {
		t.Errorf("SetMeta returned error: %v", err)
	}
	if _, ok := r.Meta["one"]; ok {
		t.Errorf("Expected deleted Meta key does exist.")
	}
	if _, ok := r.Meta["two"]; !ok {
		t.Errorf("Expected Meta key does not exist.")
	}
	t.Logf("After SetMeta, Meta: %v", r.Meta)

	// GetMeta
	meta := r.GetMeta()
	if &(r.Meta) != meta {
		t.Errorf("GetMeta error.")
	}
	// TODO Iterate/Compare?

	// WriteMeta
	err = r.WriteMeta(&Metadata{"one": 1, "two": 1})
	if err != nil {
		t.Errorf("SetMeta returned error: %v", err)
	}
	if _, ok := r.Meta["one"]; !ok {
		t.Errorf("Expected Meta key, one, does not exist.")
	}
	if _, ok := r.Meta["two"]; !ok {
		t.Errorf("Expected Meta key, two, does not exist.")
	}
	if r.Meta["two"] != 1 {
		t.Errorf("Meta value for key, two, incorrect.")
	}
	t.Logf("After WriteMeta, Meta: %v", r.Meta)

	// AppendMeta, Success
	err = r.AppendMeta(&Metadata{"three": 3})
	if err != nil {
		t.Errorf("AppendMeta returned error: %v", err)
	}
	val, ok := r.Meta["three"]
	if !ok {
		t.Errorf("Expected Meta key, three, does not exist.")
	}
	if val != 3 {
		t.Errorf("Meta value for key, three, incorrect.")
	}
	t.Logf("After AppendMeta, Meta: %v", r.Meta)

	// AppendMeta, Failure
	err = r.AppendMeta(&Metadata{"two": 2, "four": 4})
	if err == nil {
		t.Errorf("AppendMeta returned error: %v", err)
	}
	if _, ok = r.Meta["four"]; ok {
		t.Errorf("Unexpected Meta key, four, does exist.")
	}
	t.Logf("After AppendMeta, Meta: %v", r.Meta)

	err = r.AppendMeta(&Metadata{"four": 4, "two": 2})
	if err == nil {
		t.Errorf("AppendMeta returned error: %v", err)
	}
	if _, ok = r.Meta["four"]; ok {
		t.Errorf("Unexpected Meta key, four, does exist.")
	}
	t.Logf("After AppendMeta, Meta: %v", r.Meta)

	// Save
	r.UserId = user.UserId
	err = r.Save()
	if err != nil {
		t.Errorf("Error calling record.Save(): %v", err)
	}

	json, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		t.Errorf("Error in json.Marshal record: %v", err)
	}
	t.Logf("Record after save:\n%v", string(json))
}
