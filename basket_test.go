package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
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

		c.Convey("The next page body should contain \"Click here to Contribute\"", func() {
			r.ServeHTTP(record, req)
			c.So(record.Code, c.ShouldEqual, http.StatusOK)
			c.So(fmt.Sprint(record.Body), c.ShouldContainSubstring, `Click here to Contribute`)
		})
	})
}

func TestPostHomePage(t *testing.T) {
	ctx, inst := testers.GetTestingContext()
	defer inst.Close()

	c.Convey("When you register with a valid email address", t, func() {
		r := createHTTPRouter(handlers.ToHTTPHandlerConverter(ctx))
		record2 := httptest.NewRecorder()

		formData2 := url.Values{}
		formData2.Set("stripeEmail", "member@meeting.attendee.glom")
		formData2.Set("stripeToken", "tok_visa")

		req2, err := http.NewRequest("POST", "/charge", strings.NewReader(formData2.Encode()))
		req2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req2.Header.Add("Content-Length", strconv.Itoa(len(formData2.Encode())))

		c.So(err, c.ShouldBeNil)

		c.Convey("The next page body should contain \"Payment successful\"", func() {
			r.ServeHTTP(record2, req2)
			c.So(fmt.Sprint(record2.Body), c.ShouldContainSubstring, `Payment successful`)
			c.So(record2.Code, c.ShouldEqual, http.StatusOK)
		})
	})
}
