package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"event-explorer/models"

	"event-explorer/requests"
	"event-explorer/secrets"
	_ "unsafe"

	beego "github.com/beego/beego/v2/server/web"
	beegoctx "github.com/beego/beego/v2/server/web/context"
)

// Access unexported variables from the secrets package only for tests.
//
//go:linkname testBaseURL event-explorer/secrets.baseUrl
var testBaseURL secrets.BaseUrls

//
//go:linkname testSecrets event-explorer/secrets.secrets
var testSecrets secrets.Secrets

// connects a controller to a recorder, no router needed
func initController(c *beego.Controller, target string, params map[string]string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, target, nil)

	ctx := beegoctx.NewContext()
	ctx.Reset(w, r)

	for key, value := range params {
		ctx.Input.SetParam(key, value)
	}

	c.Init(ctx, "TestController", "Test", nil)
	c.EnableRender = false // no .tpl files needed in unit tests

	return w
}

// starts a fake ticketmaster and google server that always answers with status and body
func startServer(t *testing.T, status int, body string) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))

	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		GoogleBaseURL:       server.URL,
		TicketmasterBaseURL: server.URL + "/",
	}

	testSecrets = secrets.Secrets{
		GoogleKey:       "test-google-key",
		TicketmasterKey: "test-ticketmaster-key",
	}

	t.Cleanup(func() {
		server.Close()
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	})
}

func clearCaches() {
	requests.ListCache.ClearAll(context.Background())
	requests.SingleCache.ClearAll(context.Background())
}

func TestMainControllerGet(t *testing.T) {
	c := &MainController{}
	initController(&c.Controller, "/", nil)

	c.Get()

	if c.TplName != "index.tpl" {
		t.Errorf("tpl-name = %s, wanted-tpl-name index.tpl", c.TplName)
	}
}

const (
	listBody  = `{"_embedded":{"events":[{"id":"e1","name":"One"}]}}`
	eventBody = `{"id":"e1","name":"One","url":"https://www.ticketmaster.com/event/e1"}`
)

func TestListEvents(t *testing.T) {
	cases := []struct {
		name         string
		target       string
		serverStatus int
		wantTpl      string
		wantStatus   int // 0 means do not check
	}{
		{
			name:         "missing city",
			target:       "/events?countryCode=CA",
			serverStatus: http.StatusOK,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "missing country code",
			target:       "/events?city=Toronto",
			serverStatus: http.StatusOK,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "successful response",
			target:       "/events?city=Toronto&countryCode=CA",
			serverStatus: http.StatusOK,
			wantTpl:      "lists.tpl",
		},
		{
			name:         "ticketmaster not found",
			target:       "/events?city=Toronto&countryCode=CA",
			serverStatus: http.StatusNotFound,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusNotFound,
		},
		{
			name:         "ticketmaster rate limited",
			target:       "/events?city=Toronto&countryCode=CA",
			serverStatus: http.StatusTooManyRequests,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusTooManyRequests,
		},
		{
			name:         "ticketmaster server error",
			target:       "/events?city=Toronto&countryCode=CA",
			serverStatus: http.StatusInternalServerError,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusBadGateway,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearCaches()
			startServer(t, tc.serverStatus, listBody)

			c := &EventsController{}
			initController(&c.Controller, tc.target, nil)

			c.ListEvents()

			if c.TplName != tc.wantTpl {
				t.Errorf("tpl-name = %s, wanted-tpl-name %s", c.TplName, tc.wantTpl)
			}

			if tc.wantStatus != 0 && c.Ctx.Output.Status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					c.Ctx.Output.Status,
					tc.wantStatus,
				)
			}
		})
	}
}

func TestListEventsFromCache(t *testing.T) {
	clearCaches()

	cached := []models.EventCategory{{Name: "Music"}}
	key := requests.ListKey("Toronto", "CA", "6")
	requests.ListCache.Put(context.Background(), key, cached, requests.ListTTL)

	c := &EventsController{}
	initController(&c.Controller, "/events?city=Toronto&countryCode=CA", nil)

	c.ListEvents()

	got, ok := c.Data["events"].([]models.EventCategory)
	if !ok || len(got) != 1 || got[0].Name != "Music" {
		t.Errorf("events = %v, wanted-events cached categories", c.Data["events"])
	}

	if c.TplName != "lists.tpl" {
		t.Errorf("tpl-name = %s, wanted-tpl-name lists.tpl", c.TplName)
	}
}

func TestSingleEvent(t *testing.T) {
	cases := []struct {
		name         string
		eventId      string
		serverStatus int
		wantTpl      string
		wantStatus   int // 0 means do not check
	}{
		{
			name:         "missing event id",
			eventId:      "",
			serverStatus: http.StatusOK,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "successful response",
			eventId:      "e1",
			serverStatus: http.StatusOK,
			wantTpl:      "single-event.tpl",
		},
		{
			name:         "event not found",
			eventId:      "e1",
			serverStatus: http.StatusNotFound,
			wantTpl:      "error.tpl",
			wantStatus:   http.StatusNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearCaches()
			startServer(t, tc.serverStatus, eventBody)

			c := &EventsController{}
			initController(&c.Controller, "/events/"+tc.eventId, map[string]string{":eventId": tc.eventId})

			c.SingleEvent()

			if c.TplName != tc.wantTpl {
				t.Errorf("tpl-name = %s, wanted-tpl-name %s", c.TplName, tc.wantTpl)
			}

			if tc.wantStatus != 0 && c.Ctx.Output.Status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					c.Ctx.Output.Status,
					tc.wantStatus,
				)
			}
		})
	}
}

func TestSingleEventFromCache(t *testing.T) {
	clearCaches()

	cached := &models.EventResponse{ID: "e1", Name: "Cached"}
	requests.SingleCache.Put(context.Background(), requests.SingleKey("e1"), cached, requests.SingleTTL)

	c := &EventsController{}
	initController(&c.Controller, "/events/e1", map[string]string{":eventId": "e1"})

	c.SingleEvent()

	got, ok := c.Data["single-event"].(*models.EventResponse)
	if !ok || got.Name != "Cached" {
		t.Errorf("single-event = %v, wanted-single-event cached event", c.Data["single-event"])
	}

	if c.TplName != "single-event.tpl" {
		t.Errorf("tpl-name = %s, wanted-tpl-name single-event.tpl", c.TplName)
	}
}

func TestRedirectToTickets(t *testing.T) {
	cases := []struct {
		name         string
		eventId      string
		serverStatus int
		serverBody   string
		wantStatus   int
		wantLocation string
	}{
		{
			name:         "missing event id",
			eventId:      "",
			serverStatus: http.StatusOK,
			serverBody:   eventBody,
			wantStatus:   http.StatusBadRequest,
			wantLocation: "",
		},
		{
			name:         "approved ticket url",
			eventId:      "e1",
			serverStatus: http.StatusOK,
			serverBody:   eventBody,
			wantStatus:   http.StatusFound,
			wantLocation: "https://www.ticketmaster.com/event/e1",
		},
		{
			name:         "event not found",
			eventId:      "e1",
			serverStatus: http.StatusNotFound,
			serverBody:   `{}`,
			wantStatus:   http.StatusNotFound,
			wantLocation: "",
		},
		{
			name:         "event has no ticket url",
			eventId:      "e1",
			serverStatus: http.StatusOK,
			serverBody:   `{"id":"e1","name":"One"}`,
			wantStatus:   http.StatusNotFound,
			wantLocation: "",
		},
		{
			name:         "unapproved ticket url",
			eventId:      "e1",
			serverStatus: http.StatusOK,
			serverBody:   `{"id":"e1","name":"One","url":"https://evil.example.com/e1"}`,
			wantStatus:   http.StatusBadGateway,
			wantLocation: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearCaches()
			startServer(t, tc.serverStatus, tc.serverBody)

			c := &EventsController{}
			w := initController(&c.Controller, "/redirect/"+tc.eventId, map[string]string{":eventId": tc.eventId})

			c.RedirectToTickets()

			status := c.Ctx.Output.Status
			if status == 0 {
				status = w.Code // redirects only set the recorder status
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if w.Header().Get("Location") != tc.wantLocation {
				t.Errorf(
					"location = %s, wanted-location %s",
					w.Header().Get("Location"),
					tc.wantLocation,
				)
			}
		})
	}
}

func TestRedirectToTicketsFromCache(t *testing.T) {
	clearCaches()

	url := "https://www.ticketmaster.com/event/e1"
	cached := &models.EventResponse{ID: "e1", URL: url}
	requests.SingleCache.Put(context.Background(), requests.SingleKey("e1"), cached, requests.SingleTTL)

	c := &EventsController{}
	w := initController(&c.Controller, "/redirect/e1", map[string]string{":eventId": "e1"})

	c.RedirectToTickets()

	if w.Code != http.StatusFound {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			w.Code,
			http.StatusFound,
		)
	}

	if w.Header().Get("Location") != url {
		t.Errorf("location = %s, wanted-location %s", w.Header().Get("Location"), url)
	}
}

func TestInvalidateEventList(t *testing.T) {
	key := requests.ListKey("Toronto", "CA", "6")

	cases := []struct {
		name       string
		target     string
		cached     bool
		wantStatus int
	}{
		{
			name:       "missing city",
			target:     "/api/cache/events?countryCode=CA",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing country code",
			target:     "/api/cache/events?city=Toronto",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nothing cached",
			target:     "/api/cache/events?city=Toronto&countryCode=CA",
			cached:     false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "cached list is deleted",
			target:     "/api/cache/events?city=Toronto&countryCode=CA",
			cached:     true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearCaches()

			if tc.cached {
				requests.ListCache.Put(context.Background(), key, "data", requests.ListTTL)
			}

			c := &CacheController{}
			w := initController(&c.Controller, tc.target, nil)

			c.InvalidateEventList()

			if w.Code != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					w.Code,
					tc.wantStatus,
				)
			}

			exists, _ := requests.ListCache.IsExist(context.Background(), key)
			if exists {
				t.Error("list-cache-exists = true, wanted-list-cache-exists false")
			}
		})
	}
}

func TestInvalidateSingleEvent(t *testing.T) {
	key := requests.SingleKey("e1")

	cases := []struct {
		name       string
		eventId    string
		cached     bool
		wantStatus int
	}{
		{
			name:       "missing event id",
			eventId:    "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "nothing cached",
			eventId:    "e1",
			cached:     false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "cached event is deleted",
			eventId:    "e1",
			cached:     true,
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearCaches()

			if tc.cached {
				requests.SingleCache.Put(context.Background(), key, "data", requests.SingleTTL)
			}

			c := &CacheController{}
			w := initController(&c.Controller, "/api/cache/events/"+tc.eventId, map[string]string{":eventId": tc.eventId})

			c.InvalidateSingleEvent()

			if w.Code != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					w.Code,
					tc.wantStatus,
				)
			}

			exists, _ := requests.SingleCache.IsExist(context.Background(), key)
			if exists {
				t.Error("single-cache-exists = true, wanted-single-cache-exists false")
			}
		})
	}
}

func TestInvalidateAll(t *testing.T) {
	clearCaches()

	requests.ListCache.Put(context.Background(), "list-key", "data", requests.ListTTL)
	requests.SingleCache.Put(context.Background(), "single-key", "data", requests.SingleTTL)

	c := &CacheController{}
	w := initController(&c.Controller, "/api/cache/all", nil)

	c.InvalidateAll()

	if w.Code != http.StatusOK {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			w.Code,
			http.StatusOK,
		)
	}

	listExists, _ := requests.ListCache.IsExist(context.Background(), "list-key")
	if listExists {
		t.Error("list-cache-exists = true, wanted-list-cache-exists false")
	}

	singleExists, _ := requests.SingleCache.IsExist(context.Background(), "single-key")
	if singleExists {
		t.Error("single-cache-exists = true, wanted-single-cache-exists false")
	}
}

const (
	autocompleteBody = `{"suggestions":[{"placePrediction":{"placeId":"p1","text":{"text":"Toronto, Canada"}}}]}`
	placeBody        = `{"addressComponents":[
		{"longText":"Toronto","shortText":"Toronto","types":["locality"]},
		{"longText":"Canada","shortText":"CA","types":["country"]}
	]}`
)

func TestAutoCompletedPlaces(t *testing.T) {
	cases := []struct {
		name         string
		target       string
		serverStatus int
		wantStatus   int
	}{
		{
			name:         "input too short",
			target:       "/api/locations/autocomplete?input=a&sessionToken=tok",
			serverStatus: http.StatusOK,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "input too long",
			target:       "/api/locations/autocomplete?input=" + strings.Repeat("a", 101) + "&sessionToken=tok",
			serverStatus: http.StatusOK,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "missing session token",
			target:       "/api/locations/autocomplete?input=Toronto",
			serverStatus: http.StatusOK,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "successful response",
			target:       "/api/locations/autocomplete?input=Toronto&sessionToken=tok",
			serverStatus: http.StatusOK,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "google returns error",
			target:       "/api/locations/autocomplete?input=Toronto&sessionToken=tok",
			serverStatus: http.StatusForbidden,
			wantStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			startServer(t, tc.serverStatus, autocompleteBody)

			c := &PlacesController{}
			w := initController(&c.Controller, tc.target, nil)

			c.AutoCompletedPlaces()

			if w.Code != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					w.Code,
					tc.wantStatus,
				)
			}
		})
	}
}

func TestGetAPlaceById(t *testing.T) {
	cases := []struct {
		name         string
		placeId      string
		token        string
		serverStatus int
		serverBody   string
		wantStatus   int
	}{
		{
			name:         "missing place id",
			placeId:      "",
			token:        "tok",
			serverStatus: http.StatusOK,
			serverBody:   placeBody,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "place id too long",
			placeId:      strings.Repeat("a", 201),
			token:        "tok",
			serverStatus: http.StatusOK,
			serverBody:   placeBody,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "missing session token",
			placeId:      "p1",
			token:        "",
			serverStatus: http.StatusOK,
			serverBody:   placeBody,
			wantStatus:   http.StatusBadRequest,
		},
		{
			name:         "successful response",
			placeId:      "p1",
			token:        "tok",
			serverStatus: http.StatusOK,
			serverBody:   placeBody,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "place not found",
			placeId:      "p1",
			token:        "tok",
			serverStatus: http.StatusNotFound,
			serverBody:   `{}`,
			wantStatus:   http.StatusNotFound,
		},
		{
			name:         "google returns error",
			placeId:      "p1",
			token:        "tok",
			serverStatus: http.StatusForbidden,
			serverBody:   `{}`,
			wantStatus:   http.StatusInternalServerError,
		},
		{
			name:         "no city in response",
			placeId:      "p1",
			token:        "tok",
			serverStatus: http.StatusOK,
			serverBody:   `{"addressComponents":[]}`,
			wantStatus:   http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			startServer(t, tc.serverStatus, tc.serverBody)

			c := &PlacesController{}
			w := initController(
				&c.Controller,
				"/api/locations/x?sessionToken="+tc.token,
				map[string]string{":placeId": tc.placeId},
			)

			c.GetAPlaceById()

			if w.Code != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					w.Code,
					tc.wantStatus,
				)
			}
		})
	}
}
