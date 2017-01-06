package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Repo struct {
	db *gorm.DB
}

func (r *Repo) Open() error {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		return err
	}
	r.db = db

	// Migrate the schema
	r.db.AutoMigrate(&Donation{})
	return nil
}

func (r *Repo) Close() {
	r.db.Close()
}

// func example() {

// 	// Create
// 	db.Create(&Product{Code: "L1212", Price: 1000})

// 	// Read
// 	var product Product
// 	db.First(&product, 1)                   // find product with id 1
// 	db.First(&product, "code = ?", "L1212") // find product with code l1212

// 	// Update - update product's price to 2000
// 	db.Model(&product).Update("Price", 2000)

// 	// Delete - delete product
// 	db.Delete(&product)
// }

func (r *Repo) CreateDonation(refId int64, characterName string, characterID int64, date time.Time, amount float64) *Donation {

	d := Donation{
		RefID:         refId,
		CharacterName: characterName,
		CharacterID:   characterID,
		Date:          date,
		Amount:        amount,
	}

	r.db.Create(&d)

	return &d
}
