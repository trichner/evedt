package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestRepoOpen(t *testing.T) {
	repo := Repo{}

	Convey("Should open db", t, func() {

		err := repo.Open()
		So(err, ShouldBeNil)
		defer repo.Close()

		now := time.Now()
		_ = now
		repo.CreateDonation(7734553, "Thomion", 1234, time.Now(), 882232.44)
		So(err, ShouldBeNil)

	})
}
