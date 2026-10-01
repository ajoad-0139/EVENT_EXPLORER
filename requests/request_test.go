package requests

import (
	"context"
	"encoding/json"
	"event-explorer/secrets"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	_ "unsafe"
)

func TestListKey(t *testing.T) {
	cases := []struct {
		name        string
		city        string
		countryCode string
		limit       string
		want        string
	}{
		{
			name:        "normal values",
			city:        "Toronto",
			countryCode: "CA",
			limit:       "6",
			want:        "events:list:toronto:CA:6",
		},
		{
			name:        "lowercase country code",
			city:        "Toronto",
			countryCode: "ca",
			limit:       "6",
			want:        "events:list:toronto:CA:6",
		},
		{
			name:        "uppercase city",
			city:        "TORONTO",
			countryCode: "CA",
			limit:       "6",
			want:        "events:list:toronto:CA:6",
		},
		{
			name:        "city with whitespace",
			city:        " Toronto ",
			countryCode: " CA ",
			limit:       "6",
			want:        "events:list:toronto:CA:6",
		},
		{
			name:        "different limit",
			city:        "Dhaka",
			countryCode: "BD",
			limit:       "10",
			want:        "events:list:dhaka:BD:10",
		},
		{
			name:        "empty values",
			city:        "",
			countryCode: "",
			limit:       "",
			want:        "events:list:::",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := ListKey(tc.city, tc.countryCode, tc.limit)

			if res != tc.want {
				t.Errorf("list-cache-key = %s, wanted-list-cache-key %s", res, tc.want)
			}
		})
	}
}

func TestSingleKey(t *testing.T) {
	cases := []struct {
		name    string
		eventId string
		want    string
	}{
		{
			name:    "normal event id",
			eventId: "G5vZ9",
			want:    "events:single:G5vZ9",
		},
		{
			name:    "numeric event id",
			eventId: "12345",
			want:    "events:single:12345",
		},
		{
			name:    "empty event id",
			eventId: "",
			want:    "events:single:",
		},
		{
			name:    "event id with whitespace",
			eventId: " event123 ",
			want:    "events:single: event123 ",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := SingleKey(tc.eventId)

			if res != tc.want {
				t.Errorf("single-cache-key = %s, wanted-single-cache-key %s", res, tc.want)
			}
		})
	}
}

func TestCacheTTL(t *testing.T) {
	cases := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{
			name: "list cache ttl",
			got:  ListTTL,
			want: 5 * time.Minute,
		},
		{
			name: "single event cache ttl",
			got:  SingleTTL,
			want: 10 * time.Minute,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("cache-ttl = %v, wanted-cache-ttl %v", tc.got, tc.want)
			}
		})
	}
}

// Access unexported variables from the secrets package only for tests.
//
//go:linkname testBaseURL event-explorer/secrets.baseUrl
var testBaseURL secrets.BaseUrls

//
//go:linkname testSecrets event-explorer/secrets.secrets
var testSecrets secrets.Secrets

func TestGetConcurrentEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		city := r.URL.Query().Get("city")
		apiKey := r.URL.Query().Get("apikey")
		classification := r.URL.Query().Get("classificationName")
		countryCode := r.URL.Query().Get("countryCode")
		limit := r.URL.Query().Get("size")

		if apiKey != "test-api-key" {
			t.Errorf("api-key = %s, wanted-api-key test-api-key", apiKey)
		}

		if classification != "Music" {
			t.Errorf("classification-name = %s, wanted-classification-name Music", classification)
		}

		if countryCode != "CA" {
			t.Errorf("country-code = %s, wanted-country-code CA", countryCode)
		}

		if limit != "6" {
			t.Errorf("limit = %s, wanted-limit 6", limit)
		}

		switch city {

		case "Toronto":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"_embedded": map[string]interface{}{
					"events": []map[string]interface{}{
						{
							"id":   "event1",
							"name": "Music Festival",
						},
						{
							"id":   "event2",
							"name": "Concert",
						},
					},
				},
			})

		case "NotFound":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{}`))

		case "ServerError":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{}`))

		case "InvalidJSON":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`invalid json`))

		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"_embedded": {
					"events": []
				}
			}`))
		}
	}))

	defer server.Close()

	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		TicketmasterBaseURL: server.URL + "/",
		GoogleBaseURL:       "http://test-google/",
	}

	testSecrets = secrets.Secrets{
		TicketmasterKey: "test-api-key",
		GoogleKey:       "test-google-key",
	}

	defer func() {
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	}()

	cases := []struct {
		name               string
		city               string
		countryCode        string
		classificationName string
		limit              string
		wantCount          int
		wantStatus         int
		wantErr            string
	}{
		{
			name:               "successful response",
			city:               "Toronto",
			countryCode:        "CA",
			classificationName: "Music",
			limit:              "6",
			wantCount:          2,
			wantStatus:         http.StatusOK,
			wantErr:            "",
		},
		{
			name:               "not found response",
			city:               "NotFound",
			countryCode:        "CA",
			classificationName: "Music",
			limit:              "6",
			wantCount:          0,
			wantStatus:         http.StatusNotFound,
			wantErr:            "ticketmaster returned 404 Not Found",
		},
		{
			name:               "server error response",
			city:               "ServerError",
			countryCode:        "CA",
			classificationName: "Music",
			limit:              "6",
			wantCount:          0,
			wantStatus:         http.StatusInternalServerError,
			wantErr:            "ticketmaster returned 500 Internal Server Error",
		},
		{
			name:               "invalid json response",
			city:               "InvalidJSON",
			countryCode:        "CA",
			classificationName: "Music",
			limit:              "6",
			wantCount:          0,
			wantStatus:         http.StatusInternalServerError,
			wantErr:            "can not parse response",
		},
		{
			name:               "empty events",
			city:               "Unknown",
			countryCode:        "CA",
			classificationName: "Music",
			limit:              "6",
			wantCount:          0,
			wantStatus:         http.StatusOK,
			wantErr:            "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			events, status, errMsg := GetConcurrentEvents(
				context.Background(),
				tc.city,
				tc.countryCode,
				tc.classificationName,
				tc.limit,
			)

			if len(events) != tc.wantCount {
				t.Errorf(
					"event-count = %d, wanted-event-count %d",
					len(events),
					tc.wantCount,
				)
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if errMsg != tc.wantErr {
				t.Errorf(
					"error-message = %s, wanted-error-message %s",
					errMsg,
					tc.wantErr,
				)
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("classificationName")
		apiKey := r.URL.Query().Get("apikey")

		if apiKey != "test-api-key" {
			t.Errorf("api-key = %s, wanted-api-key test-api-key", apiKey)
		}

		w.Header().Set("Content-Type", "application/json")

		switch category {
		case "Music":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"_embedded": {
					"events": [
						{
							"id": "music1",
							"name": "Music Event"
						}
					]
				}
			}`))

		case "Sports":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"_embedded": {
					"events": [
						{
							"id": "sports1",
							"name": "Sports Event"
						}
					]
				}
			}`))

		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{}`))
		}
	}))

	defer server.Close()

	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		TicketmasterBaseURL: server.URL + "/",
		GoogleBaseURL:       "http://test-google/",
	}

	testSecrets = secrets.Secrets{
		TicketmasterKey: "test-api-key",
		GoogleKey:       "test-google-key",
	}

	defer func() {
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	}()

	cases := []struct {
		name        string
		city        string
		countryCode string
		limit       string
		wantStatus  int
		wantCount   int
		wantErr     string
	}{
		{
			name:        "music and sports events",
			city:        "Toronto",
			countryCode: "CA",
			limit:       "6",
			wantStatus:  http.StatusOK,
			wantCount:   2,
			wantErr:     "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, status, errMsg := GetEvents(
				t.Context(),
				tc.city,
				tc.countryCode,
				tc.limit,
			)

			if errMsg != tc.wantErr {
				t.Errorf(
					"error-message = %s, wanted-error-message %s",
					errMsg,
					tc.wantErr,
				)
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if len(res) != tc.wantCount {
				t.Errorf(
					"event-category-count = %d, wanted-event-category-count %d",
					len(res),
					tc.wantCount,
				)
			}

			if len(res) == 2 {
				if res[0].Name != "Music" {
					t.Errorf(
						"first-category = %s, wanted-first-category Music",
						res[0].Name,
					)
				}

				if res[1].Name != "Sports" {
					t.Errorf(
						"second-category = %s, wanted-second-category Sports",
						res[1].Name,
					)
				}

				if len(res[0].Events) != 1 {
					t.Errorf(
						"music-event-count = %d, wanted-music-event-count 1",
						len(res[0].Events),
					)
				}

				if len(res[1].Events) != 1 {
					t.Errorf(
						"sports-event-count = %d, wanted-sports-event-count 1",
						len(res[1].Events),
					)
				}
			}
		})
	}
}

func TestGetSingleEvent(t *testing.T) {
	cases := []struct {
		name             string
		serverStatus     int
		serverBody       string
		serverStatusText string
		wantStatus       int
		wantErrMsg       string
		wantEvent        bool
	}{
		{
			name:         "successful response",
			serverStatus: http.StatusOK,
			serverBody: `{
				"id": "event-123",
				"name": "Test Event",
				"url": "https://ticketmaster.com/event-123"
			}`,
			serverStatusText: "200 OK",
			wantStatus:       http.StatusOK,
			wantErrMsg:       "",
			wantEvent:        true,
		},
		{
			name:             "event not found",
			serverStatus:     http.StatusNotFound,
			serverBody:       `{"message":"not found"}`,
			serverStatusText: "404 Not Found",
			wantStatus:       http.StatusNotFound,
			wantErrMsg:       "event not found",
			wantEvent:        false,
		},
		{
			name:             "ticketmaster returns server error",
			serverStatus:     http.StatusInternalServerError,
			serverBody:       `{"message":"server error"}`,
			serverStatusText: "500 Internal Server Error",
			wantStatus:       http.StatusInternalServerError,
			wantErrMsg:       "ticketmaster returned 500 Internal Server Error",
			wantEvent:        false,
		},
		{
			name:             "invalid json response",
			serverStatus:     http.StatusOK,
			serverBody:       `{"id":`,
			serverStatusText: "200 OK",
			wantStatus:       http.StatusInternalServerError,
			wantErrMsg:       "can not parse response",
			wantEvent:        false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf(
						"request method = %s, wanted-method %s",
						r.Method,
						http.MethodGet,
					)
				}

				if !strings.HasPrefix(r.URL.Path, "/events/") {
					t.Errorf(
						"request path = %s, wanted path with /events/ prefix",
						r.URL.Path,
					)
				}

				if r.URL.Query().Get("apikey") != "test-api-key" {
					t.Errorf(
						"apikey = %s, wanted-api-key %s",
						r.URL.Query().Get("apikey"),
						"test-api-key",
					)
				}

				w.WriteHeader(tc.serverStatus)
				_, _ = w.Write([]byte(tc.serverBody))
			}))
			defer server.Close()

			oldBaseURL := testBaseURL
			oldSecrets := testSecrets

			testBaseURL = secrets.BaseUrls{
				TicketmasterBaseURL: server.URL + "/",
				GoogleBaseURL:       server.URL,
			}

			testSecrets = secrets.Secrets{
				TicketmasterKey: "test-api-key",
				GoogleKey:       "test-google-api-key",
			}

			defer func() {
				testBaseURL = oldBaseURL
				testSecrets = oldSecrets
			}()

			event, statusCode, errMsg := GetSingleEvent(
				t.Context(),
				"event-123",
			)

			if statusCode != tc.wantStatus {
				t.Errorf(
					"statusCode = %d, wanted-status %d",
					statusCode,
					tc.wantStatus,
				)
			}

			if errMsg != tc.wantErrMsg {
				t.Errorf(
					"errMsg = %q, wanted-errMsg %q",
					errMsg,
					tc.wantErrMsg,
				)
			}

			if tc.wantEvent {
				if event == nil {
					t.Fatal("event = nil, wanted-event")
				}

				if event.ID != "event-123" {
					t.Errorf(
						"event.ID = %q, wanted-ID %q",
						event.ID,
						"event-123",
					)
				}

				if event.Name != "Test Event" {
					t.Errorf(
						"event.Name = %q, wanted-name %q",
						event.Name,
						"Test Event",
					)
				}

				if event.URL != "https://ticketmaster.com/event-123" {
					t.Errorf(
						"event.URL = %q, wanted-URL %q",
						event.URL,
						"https://ticketmaster.com/event-123",
					)
				}

				return
			}

			if event != nil {
				t.Errorf(
					"event = %+v, wanted-event nil",
					event,
				)
			}
		})
	}
}

func TestGetSingleEventRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))

	serverURL := server.URL
	server.Close()

	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		TicketmasterBaseURL: serverURL + "/",
		GoogleBaseURL:       serverURL,
	}

	testSecrets = secrets.Secrets{
		TicketmasterKey: "test-api-key",
		GoogleKey:       "test-google-api-key",
	}

	defer func() {
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	}()

	event, statusCode, errMsg := GetSingleEvent(
		t.Context(),
		"event-123",
	)

	if event != nil {
		t.Errorf(
			"event = %+v, wanted-event nil",
			event,
		)
	}

	if statusCode != http.StatusInternalServerError {
		t.Errorf(
			"statusCode = %d, wanted-status %d",
			statusCode,
			http.StatusInternalServerError,
		)
	}

	if errMsg == "" {
		t.Error("errMsg = empty, wanted-request error")
	}
}

func TestGetAutoCompletePlaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Header.Get("X-Goog-Api-Key") {
		case "test-api-key":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"suggestions": [
					{
						"placePrediction": {
							"placeId": "place1",
							"text": {
								"text": "Toronto, Canada"
							}
						}
					},
					{
						"placePrediction": {
							"placeId": "place2",
							"text": {
								"text": "Toronto Metropolitan Area"
							}
						}
					}
				]
			}`))

		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{}`))
		}
	}))

	defer server.Close()

	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		GoogleBaseURL:       server.URL,
		TicketmasterBaseURL: "http://test-ticketmaster/",
	}

	testSecrets = secrets.Secrets{
		GoogleKey:       "test-api-key",
		TicketmasterKey: "test-ticketmaster-key",
	}

	defer func() {
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	}()

	cases := []struct {
		name       string
		input      string
		wantCount  int
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "successful autocomplete",
			input:      "Toronto",
			wantCount:  2,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "successful autocomplete with different input",
			input:      "Dhaka",
			wantCount:  2,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			suggestions, status, err := GetAutoCompletePlaces(
				t.Context(),
				tc.input,
				"test-session-token",
			)

			if len(suggestions) != tc.wantCount {
				t.Errorf(
					"suggestion-count = %d, wanted-suggestion-count %d",
					len(suggestions),
					tc.wantCount,
				)
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf(
					"has-error = %v, wanted-has-error %v",
					err != nil,
					tc.wantErr,
				)
			}

			if !tc.wantErr && len(suggestions) == 2 {
				if suggestions[0].PlaceID != "place1" {
					t.Errorf(
						"first-place-id = %s, wanted-first-place-id place1",
						suggestions[0].PlaceID,
					)
				}

				if suggestions[0].Text != "Toronto, Canada" {
					t.Errorf(
						"first-place-text = %s, wanted-first-place-text Toronto, Canada",
						suggestions[0].Text,
					)
				}

				if suggestions[1].PlaceID != "place2" {
					t.Errorf(
						"second-place-id = %s, wanted-second-place-id place2",
						suggestions[1].PlaceID,
					)
				}

				if suggestions[1].Text != "Toronto Metropolitan Area" {
					t.Errorf(
						"second-place-text = %s, wanted-second-place-text Toronto Metropolitan Area",
						suggestions[1].Text,
					)
				}
			}
		})
	}
}

func TestGetPlaceById(t *testing.T) {
	cases := []struct {
		name        string
		response    string
		statusCode  int
		wantCity    string
		wantCountry string
		wantStatus  int
		wantErr     bool
	}{
		{
			name: "locality and country",
			response: `{
				"addressComponents": [
					{
						"longText": "Toronto",
						"shortText": "Toronto",
						"types": ["locality"]
					},
					{
						"longText": "Canada",
						"shortText": "CA",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "Toronto",
			wantCountry: "CA",
			wantStatus:  http.StatusOK,
			wantErr:     false,
		},
		{
			name: "postal town fallback",
			response: `{
				"addressComponents": [
					{
						"longText": "Cambridge",
						"shortText": "Cambridge",
						"types": ["postal_town"]
					},
					{
						"longText": "United Kingdom",
						"shortText": "GB",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "Cambridge",
			wantCountry: "GB",
			wantStatus:  http.StatusOK,
			wantErr:     false,
		},
		{
			name: "administrative area fallback",
			response: `{
				"addressComponents": [
					{
						"longText": "Greater Toronto",
						"shortText": "Greater Toronto",
						"types": ["administrative_area_level_2"]
					},
					{
						"longText": "Canada",
						"shortText": "CA",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "Greater Toronto",
			wantCountry: "CA",
			wantStatus:  http.StatusOK,
			wantErr:     false,
		},
		{
			name: "locality wins over fallback",
			response: `{
				"addressComponents": [
					{
						"longText": "Greater Toronto",
						"shortText": "Greater Toronto",
						"types": ["administrative_area_level_2"]
					},
					{
						"longText": "Toronto",
						"shortText": "Toronto",
						"types": ["locality"]
					},
					{
						"longText": "Canada",
						"shortText": "CA",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "Toronto",
			wantCountry: "CA",
			wantStatus:  http.StatusOK,
			wantErr:     false,
		},
		{
			name: "missing city",
			response: `{
				"addressComponents": [
					{
						"longText": "Canada",
						"shortText": "CA",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "",
			wantCountry: "CA",
			wantStatus:  http.StatusUnprocessableEntity,
			wantErr:     true,
		},
		{
			name: "invalid country code",
			response: `{
				"addressComponents": [
					{
						"longText": "Toronto",
						"shortText": "Toronto",
						"types": ["locality"]
					},
					{
						"longText": "Canada",
						"shortText": "C",
						"types": ["country"]
					}
				]
			}`,
			statusCode:  http.StatusOK,
			wantCity:    "",
			wantCountry: "",
			wantStatus:  http.StatusUnprocessableEntity,
			wantErr:     true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.response))
			}))

			defer server.Close()

			oldBaseURL := testBaseURL
			oldSecrets := testSecrets

			testBaseURL = secrets.BaseUrls{
				GoogleBaseURL:       server.URL,
				TicketmasterBaseURL: "http://test-ticketmaster/",
			}

			testSecrets = secrets.Secrets{
				GoogleKey:       "test-google-api-key",
				TicketmasterKey: "test-ticketmaster-key",
			}

			defer func() {
				testBaseURL = oldBaseURL
				testSecrets = oldSecrets
			}()

			city, status, err := GetPlaceById(
				t.Context(),
				"test-place-id",
				"test-session-token",
			)

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf(
					"has-error = %v, wanted-has-error %v",
					err != nil,
					tc.wantErr,
				)
			}

			if tc.wantErr {
				return
			}

			if city == nil {
				t.Fatal("city = nil, wanted-city response")
			}

			if city.City != tc.wantCity {
				t.Errorf(
					"city = %s, wanted-city %s",
					city.City,
					tc.wantCity,
				)
			}

			if city.CountryCode != tc.wantCountry {
				t.Errorf(
					"country-code = %s, wanted-country-code %s",
					city.CountryCode,
					tc.wantCountry,
				)
			}
		})
	}
}

// sets the test urls and keys, returns a function that restores the old ones
func setTestConfig(googleURL string, ticketmasterURL string) func() {
	oldBaseURL := testBaseURL
	oldSecrets := testSecrets

	testBaseURL = secrets.BaseUrls{
		GoogleBaseURL:       googleURL,
		TicketmasterBaseURL: ticketmasterURL,
	}

	testSecrets = secrets.Secrets{
		GoogleKey:       "test-google-key",
		TicketmasterKey: "test-ticketmaster-key",
	}

	return func() {
		testBaseURL = oldBaseURL
		testSecrets = oldSecrets
	}
}

func TestInvalidBaseURL(t *testing.T) {
	defer setTestConfig("http://[bad", "http://[bad")()

	cases := []struct {
		name       string
		wantStatus int
		call       func() (int, bool)
	}{
		{
			name:       "concurrent events",
			wantStatus: 400,
			call: func() (int, bool) {
				_, status, errMsg := GetConcurrentEvents(t.Context(), "Toronto", "CA", "Music", "6")
				return status, errMsg != ""
			},
		},
		{
			name:       "single event",
			wantStatus: 400,
			call: func() (int, bool) {
				_, status, errMsg := GetSingleEvent(t.Context(), "event-123")
				return status, errMsg != ""
			},
		},
		{
			name:       "autocomplete places",
			wantStatus: 500,
			call: func() (int, bool) {
				_, status, err := GetAutoCompletePlaces(t.Context(), "Toronto", "token")
				return status, err != nil
			},
		},
		{
			name:       "place by id",
			wantStatus: 500,
			call: func() (int, bool) {
				_, status, err := GetPlaceById(t.Context(), "place-1", "token")
				return status, err != nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, hasErr := tc.call()

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if !hasErr {
				t.Error("has-error = false, wanted-has-error true")
			}
		})
	}
}

func TestConnectionRefused(t *testing.T) {
	// start a server and close it, so nothing is listening on its url
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	closedURL := server.URL
	server.Close()

	defer setTestConfig(closedURL, closedURL+"/")()

	cases := []struct {
		name string
		call func() (int, bool)
	}{
		{
			name: "concurrent events",
			call: func() (int, bool) {
				_, status, errMsg := GetConcurrentEvents(t.Context(), "Toronto", "CA", "Music", "6")
				return status, errMsg != ""
			},
		},
		{
			name: "all events",
			call: func() (int, bool) {
				_, status, errMsg := GetEvents(t.Context(), "Toronto", "CA", "6")
				return status, errMsg != ""
			},
		},
		{
			name: "autocomplete places",
			call: func() (int, bool) {
				_, status, err := GetAutoCompletePlaces(t.Context(), "Toronto", "token")
				return status, err != nil
			},
		},
		{
			name: "place by id",
			call: func() (int, bool) {
				_, status, err := GetPlaceById(t.Context(), "place-1", "token")
				return status, err != nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, hasErr := tc.call()

			if status != http.StatusInternalServerError {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					http.StatusInternalServerError,
				)
			}

			if !hasErr {
				t.Error("has-error = false, wanted-has-error true")
			}
		})
	}
}

func TestGetEventsUpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	defer setTestConfig("http://test-google/", server.URL+"/")()

	res, status, errMsg := GetEvents(t.Context(), "Toronto", "CA", "6")

	if res != nil {
		t.Errorf("events = %v, wanted-events nil", res)
	}

	if status != http.StatusInternalServerError {
		t.Errorf(
			"response-status = %d, wanted-response-status %d",
			status,
			http.StatusInternalServerError,
		)
	}

	if errMsg != "ticketmaster returned 500 Internal Server Error" {
		t.Errorf(
			"error-message = %s, wanted-error-message %s",
			errMsg,
			"ticketmaster returned 500 Internal Server Error",
		)
	}
}

func TestGetAutoCompletePlacesErrors(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		response   string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "google returns error status",
			statusCode: http.StatusForbidden,
			response:   `{}`,
			wantStatus: 500,
			wantErr:    "google returned status 403",
		},
		{
			name:       "invalid json",
			statusCode: http.StatusOK,
			response:   `not json`,
			wantStatus: 500,
			wantErr:    "invalid response from google",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			defer setTestConfig(server.URL, "http://test-ticketmaster/")()

			suggestions, status, err := GetAutoCompletePlaces(t.Context(), "Toronto", "token")

			if suggestions != nil {
				t.Errorf("suggestions = %v, wanted-suggestions nil", suggestions)
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if err == nil || err.Error() != tc.wantErr {
				t.Errorf("error = %v, wanted-error %s", err, tc.wantErr)
			}
		})
	}
}

func TestGetPlaceByIdErrors(t *testing.T) {
	cases := []struct {
		name       string
		statusCode int
		response   string
		wantStatus int
		wantErr    string
	}{
		{
			name:       "place not found",
			statusCode: http.StatusNotFound,
			response:   `{}`,
			wantStatus: 404,
			wantErr:    "place not found",
		},
		{
			name:       "google returns error status",
			statusCode: http.StatusForbidden,
			response:   `{}`,
			wantStatus: 500,
			wantErr:    "google returned status 403",
		},
		{
			name:       "invalid json",
			statusCode: http.StatusOK,
			response:   `not json`,
			wantStatus: 500,
			wantErr:    "invalid response from google",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			defer setTestConfig(server.URL, "http://test-ticketmaster/")()

			city, status, err := GetPlaceById(t.Context(), "place-1", "token")

			if city != nil {
				t.Errorf("city = %v, wanted-city nil", city)
			}

			if status != tc.wantStatus {
				t.Errorf(
					"response-status = %d, wanted-response-status %d",
					status,
					tc.wantStatus,
				)
			}

			if err == nil || err.Error() != tc.wantErr {
				t.Errorf("error = %v, wanted-error %s", err, tc.wantErr)
			}
		})
	}
}
