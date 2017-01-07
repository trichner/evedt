package main

import (
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

func (r *Repo) StoreDonation(donation *Donation) {
	r.db.Create(donation)
}

func (r *Repo) LastDonation() *Donation {
	donations := make([]Donation, 1)

	r.db.Limit(1).Order("ref_id desc").Find(&donations)
	if len(donations) == 0 {
		return nil
	}
	return &donations[0]
}

func (r *Repo) FindDonations(limit int) Donations {

	donations := make([]Donation, limit)

	r.db.Limit(limit).Order("ref_id desc").Find(&donations)

	return donations
}
