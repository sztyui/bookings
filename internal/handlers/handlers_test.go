package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/sztyui/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},
	// new router
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test case where reservation ID is not exists
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)

	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

func reqBodyBuilder(params map[string]string) string {
	var result string
	if len(params) == 0 {
		return ""
	}
	for key, value := range params {
		result = fmt.Sprintf("%s&%s=%s", result, key, value)
	}
	return result[1:]
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		FirstName: "John",
		LastName:  "Smith",
		Email:     "john@smith.com",
		Phone:     "123456789",
		StartDate: time.Date(2050, 01, 01, 00, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2050, 01, 02, 00, 0, 0, 0, time.UTC),
		RoomID:    1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	session.Put(ctx, "reservation", reservation)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing request body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong response code for missing post body: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	param := []struct {
		body       map[string]string
		resultcode int
	}{
		{
			map[string]string{"start_date": "invalid", "end_date": "2050-01-02", "first_name": "John", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "1"},
			http.StatusSeeOther,
		},
		{
			map[string]string{"start_date": "2050-01-01", "end_date": "invalid", "first_name": "John", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "1"},
			http.StatusSeeOther,
		},
		{
			map[string]string{"start_date": "2050-01-01", "end_date": "2050-01-02", "first_name": "John", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "invalid"},
			http.StatusSeeOther,
		},
		{
			map[string]string{"start_date": "2050-01-01", "end_date": "2050-01-02", "first_name": "J", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "1"},
			http.StatusOK,
		},
		{
			map[string]string{"start_date": "2050-01-01", "end_date": "2050-01-02", "first_name": "John", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "2"},
			http.StatusSeeOther,
		},
		{
			map[string]string{"start_date": "2050-01-01", "end_date": "2050-01-02", "first_name": "John", "last_name": "Smith", "email": "john@smith.com", "phone": "123456789", "room_id": "3"},
			http.StatusSeeOther,
		},
	}

	for _, value := range param {
		req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBodyBuilder(value.body)))
		ctx = getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr = httptest.NewRecorder()
		handler = http.HandlerFunc(Repo.PostReservation)
		handler.ServeHTTP(rr, req)
		if rr.Code != value.resultcode {
			t.Errorf("Reservation handler returned wrong response code for invalid input. Got %d, wanted %d, input: %s", rr.Code, value.resultcode, value.body)
		}
	}
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	postData := []struct {
		params         map[string]string
		expectedResult bool
	}{
		{map[string]string{"start_date": "2050-09-01", "end_date": "2050-09-02", "room_id": "1"}, true},
		{map[string]string{"start_date": "2050-09-01", "end_date": "2050-09-02", "room_id": "2"}, false},
		{map[string]string{"start_date": "2050-09-01", "end_date": "2050-09-02", "room_id": "3"}, false},
		{map[string]string{"start_date": "2050-09-01", "end_date": "2050-09-02", "room_id": "a"}, false},
		{map[string]string{}, false},
	}

	for _, val := range postData {
		req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBodyBuilder(val.params)))

		// get context with session
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set http request header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AvailabilityJSON)
		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal([]byte(rr.Body.String()), &j)
		if err != nil {
			t.Error(err)
		}
		if j.OK != val.expectedResult {
			t.Errorf("error with the json response: %v", j)
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{"valid-credentials", "me@here.ca", http.StatusSeeOther, "", "/"},
	{"invalid-credentials", "jack@nimble.com", http.StatusSeeOther, "", "/user/login"},
	{"invalid-data", "j", http.StatusOK, `action="/user/login"`, ""},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed  %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLocation, _ := rr.Result().Location()
			if actualLocation.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s but got location %s", e.name, e.expectedLocation, actualLocation.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var bookTests = []struct {
	Name               string
	RoomID             int
	startDate          string
	endDate            string
	expectedStatusCode int
	expectedLocation   string
}{
	{"bookTest valid", 1, "2021-06-26", "2021-06-27", http.StatusOK, "/make-reservation"},
	//{"bookTest invalid", 1, "2011-13-22", "2021-06-40", http.StatusInternalServerError, ""},
}

func TestRepository_BookRoom(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range bookTests {
		urlParam := fmt.Sprintf("?id=%d&s=%s&e=%s", e.RoomID, e.startDate, e.endDate)
		resp, err := ts.Client().Get(ts.URL + "/book-room" + urlParam)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("error with %s: got %d but expected %d", e.Name, resp.StatusCode, e.expectedStatusCode)
		}
	}
}
