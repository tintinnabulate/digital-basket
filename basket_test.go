package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
	"github.com/tintinnabulate/aecontext-handlers/handlers"
	"github.com/tintinnabulate/aecontext-handlers/testers"
)

func TestMain(m *testing.M) {
	testSetup()
	retCode := m.Run()
	os.Exit(retCode)
}

func testSetup() {
	configInit("config.example.json")
	stripeInit()
}

// TestGetHomePage does just that
func TestGetHomePage(t *testing.T) {
	ctx, inst := testers.GetTestingContext()
	defer inst.Close()

	c.Convey("When you visit the root URL", t, func() {
		r := createHTTPRouter(handlers.ToHTTPHandlerConverter(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("GET", "/", nil)

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"I am an alcoholic\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `I am an alcoholic`)
		})
	})
}

func TestPostHomePage(t *testing.T) {
	ctx, inst := testers.GetTestingContext()
	defer inst.Close()

	c.Convey("When you post at the root URL", t, func() {
		r := createHTTPRouter(handlers.ToHTTPHandlerConverter(ctx))
		record := httptest.NewRecorder()

		req, err := http.NewRequest("POST", "/", nil)

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"Charge successful\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Charge successful`)
		})
	})
}
