package tracker

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRepoOpen(t *testing.T) {
	repo := NewSqliteRepo()

	Convey("Should open db", t, func() {

		err := repo.Open()
		So(err, ShouldBeNil)
		defer repo.Close()

		Convey("Should store donations", func() {

			now := time.Now()
			_ = now
			repo.StoreDonation(&Donation{
				CharacterID:   7734553,
				CharacterName: "Thomion",
				RefID:         1234,
				Date:          time.Now(),
				Amount:        882232.44,
			})
			So(err, ShouldBeNil)

			last := repo.LastDonation()
			So(last, ShouldNotBeNil)

			So(last.RefID, ShouldEqual, 1234)
		})

		Convey("Should find donations", func() {
			donations := repo.FindDonations(1000, 5)

			fmt.Printf("\nFound %d donations:\n", len(donations))
			for _, d := range donations {
				fmt.Println(d)
			}
		})

		Convey("Should find top donations", func() {
			donations := repo.FindTopDonations(5, 0)

			fmt.Printf("\nFound %d top donations:\n", len(donations))
			for _, d := range donations {
				fmt.Println(d)
			}
		})

	})

}
