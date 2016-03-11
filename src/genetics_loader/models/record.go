package models

import (
	"fmt"
	"genetics_loader/driver/database"
	"time"
)

type Record struct {
	ID          int         `json:"id"`
	UserId      UUID        `sql:"type:uuid" json:"-"`
	Type        Type        `json:"type"`
	TypeId      int         `json:"-"`
	Entity      Entity      `json:"entity"`
	EntityId    int         `json:"-"`
	Measure     Measure     `json:"measure"`
	MeasureId   int         `sql:"default:NULL" json:"-"`
	MeasureData MeasureInfo `sql:"type:jsonb" json:"measure_data"`
	RecordAt    time.Time   `json:"record_at"`
	Meta        Metadata    `sql:"type:jsonb" json:"meta"`
	Timestamps
	SoftDelete
}

// Load record fields based on given Entity
func (r *Record) FromEntity(e *Entity) error {
	var measures []Measure

	// Set Entity, EntityId
	r.Entity = *e
	r.EntityId = e.ID

	// Set Type, TypeId
	if r.Entity.Type.ID == 0 {
		var myType Type
		err := database.App.Where("id = ?", r.Entity.TypeId).First(&myType).Error
		if err != nil {
			return fmt.Errorf("Database error looking up type with id %v: %v", r.Entity.TypeId, err)
		}
		r.Entity.Type = myType
	}
	r.Type = r.Entity.Type
	r.TypeId = r.Type.ID

	// Set Measure, MeasureId
	database.App.Model(e).Association("Measures").Find(&measures)
	if len(measures) != 1 {
		return fmt.Errorf("Entity has multiple measures, cannot automatically populate. Entity: %v", e)
	}
	r.Measure = measures[0]
	r.MeasureId = measures[0].ID

	// Copy Measure template
	r.MeasureData = r.Measure.Template

	return nil
}

// Similar to FromEntity method, but adds all measure templates to
// the record's MeasureData map
func (r *Record) CreateFromEntity(e *Entity) error {
	var measures []Measure

	// Set Entity, EntityId
	r.Entity = *e
	r.EntityId = e.ID

	// Set Type, TypeId
	if r.Entity.Type.ID == 0 {
		var myType Type
		err := database.App.Where("id = ?", r.Entity.TypeId).First(&myType).Error
		if err != nil {
			return fmt.Errorf("Database error looking up type with id %v: %v", r.Entity.TypeId, err)
		}
		r.Entity.Type = myType
	}
	r.Type = r.Entity.Type
	r.TypeId = r.Type.ID

	// Set Measure, MeasureId
	database.App.Model(e).Association("Measures").Find(&measures)
	if len(measures) < 1 {
		return fmt.Errorf("Entity has no measures. Entity: %v", e)
	}
	r.Measure = measures[0]
	r.MeasureId = measures[0].ID

	// create record MeasureData (MeasureInfo which is map[string]MeasureData)
	r.MeasureData = make(MeasureInfo, 0)

	// go through all the measures and copy their templates
	// into the record's MeasureData map by key
	for _, measure := range measures {
		var key string = ""
		for k := range measure.Template {
			key = k
			r.MeasureData[key] = measure.Template[key]
		}
		//		r.MeasureData[key] = measure.Template[key]
	}

	r.Measure.Template = r.MeasureData

	return nil
}

// Set measure value for given key
func (r *Record) SetMeasureValue(key string, value float64) error {
	data, ok := r.MeasureData[key]
	if !ok {
		return fmt.Errorf("Specified measure key, '%v', does not exist in record's MeasureData.", key)
	}
	data.Value = value
	r.MeasureData[key] = data
	return nil
}

// Set measure string value for given key
func (r *Record) SetMeasureStringValue(key string, value string) error {
	data, ok := r.MeasureData[key]
	if !ok {
		return fmt.Errorf("Specified measure key, '%v', does not exist in record's MeasureData.", key)
	}
	data.StringValue = value
	r.MeasureData[key] = data
	return nil
}

// Get measure value for given key
func (r *Record) GetMeasureValue(key string) (float64, error) {
	data, ok := r.MeasureData[key]
	if !ok {
		return 0, fmt.Errorf("Specified measure key, '%v', does not exist in record's MeasureData.", key)
	}
	return data.Value, nil
}

// Get measure string value for given key
func (r *Record) GetMeasureStringValue(key string) (string, error) {
	data, ok := r.MeasureData[key]
	if !ok {
		return "", fmt.Errorf("Specified measure key, '%v', does not exist in record's MeasureData.", key)
	}
	return data.StringValue, nil
}

// Set Record Metadata, discards old Metadata
func (r *Record) SetMeta(newMeta *Metadata) error {
	r.Meta = *newMeta
	return nil
}

// Returns Record Metadata
func (r *Record) GetMeta() *Metadata {
	return &r.Meta
}

// Writes key/values from newMeta to Record Metadata.
// Overwrites existing values that are present in newMeta.
// Leaves other existing values intact.
func (r *Record) WriteMeta(newMeta *Metadata) error {
	for key, value := range *newMeta {
		r.Meta[key] = value
	}
	return nil
}

// Adds key/values from newMeta to Record Metadata.
// Throws error if key in newMeta already exists (and value is different?) in Record
func (r *Record) AppendMeta(newMeta *Metadata) error {
	// TODO Looping twice is inefficient but simple. Improve by operating on a copy? (Deep copy required?)
	for key, value := range *newMeta {
		current, ok := r.Meta[key]
		if ok && current != value {
			return fmt.Errorf("Overwrite error.")
		}
	}
	for key, value := range *newMeta {
		r.Meta[key] = value
	}
	return nil
}

// Save the record to the database.
// Assumes the relatedType, Measure, Entity do not/have not changed.
func (r *Record) Save() error {
	var tempType Type
	var tempMeasure Measure
	var tempEntity Entity

	if r.TypeId != r.Type.ID || r.EntityId != r.Entity.ID || r.MeasureId != r.Measure.ID {
		return fmt.Errorf("Error: record's related type/entity/measure out of sync with ID.")
	}
	if r.UserId.String() == "" {
		return fmt.Errorf("Record.UserId must be set.")
	}

	// Prevent gorm trying to save/update related
	r.TypeId = r.Type.ID
	tempType = r.Type
	r.Type = Type{}

	r.MeasureId = r.Measure.ID
	tempMeasure = r.Measure
	r.Measure = Measure{}

	r.EntityId = r.Entity.ID
	tempEntity = r.Entity
	r.Entity = Entity{}

	var err error
	if r.ID == 0 {
		err = database.App.Create(r).Error
	} else {
		err = database.App.Save(r).Error
	}
	if err != nil {
		return err
	}

	r.Type = tempType
	r.Measure = tempMeasure
	r.Entity = tempEntity

	return nil
}
