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

	"github.com/p3rfect05/go_proj/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name   string
	url    string
	method string

	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"make-reserves", "/make-reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-03"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-03"},
	// }, http.StatusOK},
	// {"post-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "alex"},
	// 	{key: "last_name", value: "v"},
	// 	{key: "email", value: "a@a.com"},
	// 	{key: "phone", value: "333-333-3333"},
	// }, http.StatusOK},
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
				t.Errorf("for %s expected %d, got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}

	}

}

func TestRepository_Reservation(t *testing.T) {
	res := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "General's quarters",
		},
	}
	req := httptest.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", res)

	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code, got %d, wanted %d", rr.Code, http.StatusOK)
	}
	// reservation is not in session
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code, got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test with non-existing room
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	res.RoomID = 0
	session.Put(ctx, "reservation", res)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code, got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2050-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-10")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=alex")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=xela")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=z@z.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=1-500")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req := httptest.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code, got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_JSONAvailability(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start", "2050-01-01")
	postedData.Add("end", "2050-01-01")
	postedData.Add("room_id", "2")


	//create request
	req := httptest.NewRequest("POST", "/search-availability-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	
	//set the request header
	req.Header.Set("Content-Type", "x-www-form-urlenencoded")
	// make handler a handler func
	handler := http.HandlerFunc(Repo.JSONAvailability)
	rr := httptest.NewRecorder()
	//make request to our handler
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("Failed to parse json response")
	}

}
