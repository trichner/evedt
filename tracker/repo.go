package tracker

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const revTypeIdDonation = 10

// Repo stores all
type Repo interface {
	Open() error
	Close() error
	StoreDonation(donation *Donation) error
	LastDonation() *Donation
	FindDonations(limit int, days int) Donations
	FindTopDonations(limit int, days int) Donations
}

type gormRepo struct {
	db *gorm.DB
}

func NewSqliteRepo() Repo {
	return &gormRepo{}
}

func (r *gormRepo) Open() error {
	db, err := gorm.Open("sqlite3", "evedt.db")
	if err != nil {
		return err
	}
	r.db = db

	// Migrate the schema
	r.db.AutoMigrate(&Donation{})
	return nil
}

func (r *gormRepo) Close() error {
	return r.db.Close()
}

func (r *gormRepo) StoreDonation(donation *Donation) error {
	return r.db.Create(donation).Error
}

func (r *gormRepo) LastDonation() *Donation {
	donations := make([]Donation, 1)

	r.db.Limit(1).Order("ref_id desc").Find(&donations)
	if len(donations) == 0 {
		return nil
	}
	return &donations[0]
}

func (r *gormRepo) FindDonations(limit int, days int) Donations {

	db := r.prepareFindDonations(limit, days)

	donations := make([]Donation, limit)
	db.Order("ref_id desc").Find(&donations)

	return donations
}

func (r *gormRepo) FindTopDonations(limit int, days int) Donations {

	db := r.prepareFindDonations(limit, days)

	donations := make([]Donation, limit)
	db.Order("amount desc").Find(&donations)

	return donations
}

func (r *gormRepo) prepareFindDonations(limit int, days int) *gorm.DB {

	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	db := r.db.Limit(limit)

	limitDays := 365
	if days > 0 && days < 365 {
		limitDays = days
	}

	firstDay := time.Now().Add(-(time.Hour * 24 * time.Duration(limitDays)))
	db = db.Where("date >= ? AND ref_type_id = ?", firstDay, revTypeIdDonation)

	return db
}
