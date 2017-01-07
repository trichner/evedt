package tracker

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

func (r *Repo) FindDonations(limit int, days int) Donations {

	db := r.prepareFindDonations(limit, days)

	donations := make([]Donation, limit)
	db.Order("ref_id desc").Find(&donations)

	return donations
}

func (r *Repo) FindTopDonations(limit int, days int) Donations {

	db := r.prepareFindDonations(limit, days)

	donations := make([]Donation, limit)
	db.Order("amount desc").Find(&donations)

	return donations
}

func (r *Repo) prepareFindDonations(limit int, days int) *gorm.DB {

	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	db := r.db.Limit(limit)

	limitDays := 365
	if days > 0 && days < 365 {
		limitDays = days
	}

	firstDay := time.Now().Add(-(time.Hour * 24 * time.Duration(limitDays)))
	db = db.Where("date >= ?", firstDay)

	return db
}
