package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"genetics_loader/driver/database"
	"genetics_loader/models"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type GeneticRecord struct {
	ContainerID string    `json:"-"`    //col 1
	Result1     string    `json:"FTO"`  //col 8
	Result3     string    `json:"DRD2"` //col 22
	Result2     string    `json:"MC4R"` //col 15
	Date        time.Time `json:"-"`    //col 24
}

func main() {

	log.Print("Load Genetics Data from CSV")

	err := loadGeneticsData("genetics_data.csv")
	if err != nil {
		panic(err)
	}

	log.Print("End CSV Load")
	return
}

func loadGeneticsData(filename string) (err error) {

	//Check if CSV file
	ext := filepath.Ext(filename)
	if ext != ".csv" {
		err := errors.New("Error: Input file is not .csv")
		return err
	}

	//Open File
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := csv.NewReader(file)

	for i := 0; ; i++ {
		var geneticsRecord GeneticRecord

		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		geneticsRecord = GeneticRecord{
			ContainerID: row[1],
			Result1:     row[8],
			Result2:     row[15],
			Result3:     row[22],
		}

		// Date
		t, err := time.Parse("2006-01-02", row[24])
		if err != nil {
			return err
		}
		geneticsRecord.Date = t
		//log.Print(geneticsRecord)

		// Begin TXs
		app := database.App.Begin()
		if app.Error != nil {
			log.Fatalf("Error starting transaction(s).\n\tApp: %v\n", app.Error)
			err = app.Error
			return err
		}

		//Get USER ID from Container ID
		rowID := app.Raw("select distinct o.participant_id from orders o, order_items oi, sku_items si where oi.order_id = o.id and oi.sku_item_id = si.id and si.name like 'Test Tube' and oi.tracking_value = ?", geneticsRecord.ContainerID).Row()
		var user_id interface{}
		rowID.Scan(&user_id)

		if user_id == nil {
			app.Rollback()
			log.Print("Error: UserId from ContainerID not found in orders and order_items: ", geneticsRecord.ContainerID)
			continue
		}
		userId := string(user_id.([]uint8))

		//Check if UserID already has a Genetics Record
		var existGeneticsRecord models.Record
		err = app.Where("entity_id in (select id from entities where name = 'GeneticResults') and user_id = ?", userId).Find(&existGeneticsRecord).Error
		if err == nil {
			app.Rollback()
			log.Print("Genetics record already exists for ContainerID: ", geneticsRecord.ContainerID, " userId: ", userId)
			continue
		}

		//Create New Record if no Genetics Record exists
		var newRecord models.Record
		var entity models.Entity

		err = app.Where("name = ?", "GeneticResults").Find(&entity).Error
		if err != nil {
			app.Rollback()
			return err
		}

		err = newRecord.CreateFromEntity(&entity)
		if err != nil {
			app.Rollback()
			return err
		}

		newRecord.RecordAt = geneticsRecord.Date

		// Meta
		bytes, err := json.Marshal(geneticsRecord)
		if err != nil {
			app.Rollback()
			return err
		}

		var meta map[string]interface{}
		err = json.Unmarshal(bytes, &meta)
		if err != nil {
			app.Rollback()
			return err
		}

		var resultMeta = map[string]interface{}{
			"results": meta,
		}

		newRecord.Meta = resultMeta

		var id models.UUID
		id.Parse(userId)
		newRecord.UserId = id

		err = app.Create(&newRecord).Error
		if err != nil {
			return err
		}

		err = app.Commit()
		if err != nil {
			app.Rollback()
			return err
		}
		log.Print("Successful ingestion of GeneticResults Record for ContainerID: ", geneticsRecord.ContainerID, " UserId: ", userId)

	}
	return err

}

// 		err = app.Where("meta #>> '{serials,Medivo}' = ? and meta #>> '{serials,TestTube}' = ?", row[0], row[1]).Find(&serialRecord).Error
// 		if err != nil {
// 			app.Rollback()
// 			log.Print("No Medivo and TestTube data for this participant: ", row[0], ", ", row[1])
// 			return err
// 		}
